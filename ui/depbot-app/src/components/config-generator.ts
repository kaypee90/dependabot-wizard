import { configStore } from "../store/configuration-store";

export type Group = {
  dependencyType?: string;
  patterns?: string[];
  updateTypes?: string[];
};

export type IgnoreRule = {
  dependencyName: string;
  updateTypes?: string[];
};

/** Named, ready-made grouping configurations. Mirrors the Go definitions. */
export type GroupingStrategy = 'none' | 'all' | 'by-dependency-type';

export type UpdateBlock = {
  packageEcosystem: string;
  directory: string;
  interval: string;
  reviewers: string[];
  openPullRequestsLimit?: number;
  labels: string[];
  groups?: Record<string, Group>;
  ignore?: IgnoreRule[];
};

export const GROUPING_LABELS: Record<GroupingStrategy, string> = {
  'none': 'No grouping — one pull request per dependency',
  'all': 'All minor and patch updates together, majors separate',
  'by-dependency-type': 'Production and development dependencies separately',
};

/**
 * Grouping is the single biggest reducer of pull request volume, which is why
 * the health report recommends it and why it is offered here: the diagnosis
 * and the fix belong in the same tool.
 */
export const groupsFor = (strategy: GroupingStrategy): Record<string, Group> | undefined => {
  switch (strategy) {
    case 'all':
      return {
        'all-minor-and-patch': { patterns: ['*'], updateTypes: ['minor', 'patch'] },
      };
    case 'by-dependency-type':
      return {
        'production-dependencies': {
          dependencyType: 'production', patterns: ['*'], updateTypes: ['minor', 'patch'],
        },
        'development-dependencies': {
          dependencyType: 'development', patterns: ['*'], updateTypes: ['minor', 'patch'],
        },
      };
    default:
      return undefined;
  }
};

/** Stop opening major bumps, and take them deliberately instead. */
export const ignoreMajorUpdates = (): IgnoreRule[] => ([
  { dependencyName: '*', updateTypes: ['version-update:semver-major'] },
]);

const INDENT = "  ";

// Characters that would change how YAML reads a bare scalar. Anything holding
// one gets quoted rather than emitted raw.
const NEEDS_QUOTING = /^$|^[-?:,[\]{}#&*!|>'"%@`]|[:#]\s|\s$|^\s|^(true|false|null|yes|no|on|off|~)$/i;

const scalar = (value: string): string => {
  const text = String(value);
  return NEEDS_QUOTING.test(text) ? JSON.stringify(text) : text;
};

// A sequence of plain strings, rendered at the given indent level. yaml.v2 on
// the Go side does not indent sequence items under their key, so neither do we
// and the two paths produce identical files.
const sequence = (key: string, values: string[], indent: string): string[] => {
  const present = values.map((value) => value.trim()).filter((value) => value.length > 0);
  if (present.length === 0) {
    return [];
  }

  return [
    indent + key + ":",
    ...present.map((value) => indent + "- " + scalar(value)),
  ];
};

const renderGroups = (groups: Record<string, Group>, indent: string): string[] => {
  const names = Object.keys(groups).sort();
  if (names.length === 0) {
    return [];
  }

  const lines = [indent + 'groups:'];
  names.forEach((name) => {
    const group = groups[name];
    const inner = indent + INDENT;
    lines.push(inner + scalar(name) + ':');

    if (group.dependencyType) {
      lines.push(inner + INDENT + 'dependency-type: ' + scalar(group.dependencyType));
    }
    lines.push(...sequence('patterns', group.patterns ?? [], inner + INDENT));
    lines.push(...sequence('update-types', group.updateTypes ?? [], inner + INDENT));
  });

  return lines;
};

const renderIgnore = (rules: IgnoreRule[], indent: string): string[] => {
  if (rules.length === 0) {
    return [];
  }

  const lines = [indent + 'ignore:'];
  rules.forEach((rule) => {
    lines.push(indent + '- dependency-name: ' + scalar(rule.dependencyName));
    lines.push(...sequence('update-types', rule.updateTypes ?? [], indent + INDENT));
  });

  return lines;
};

export const renderUpdate = (update: UpdateBlock): string[] => {
  const indent = INDENT;
  const lines: string[] = [
    "- package-ecosystem: " + scalar(update.packageEcosystem),
    indent + "directory: " + scalar(update.directory),
    indent + "schedule:",
    indent + INDENT + "interval: " + scalar(update.interval),
  ];

  lines.push(...sequence("reviewers", update.reviewers, indent));

  if (typeof update.openPullRequestsLimit === "number" && !Number.isNaN(update.openPullRequestsLimit)) {
    lines.push(indent + "open-pull-requests-limit: " + update.openPullRequestsLimit);
  }

  lines.push(...sequence("labels", update.labels, indent));

  if (update.groups) {
    lines.push(...renderGroups(update.groups, indent));
  }

  if (update.ignore && update.ignore.length > 0) {
    lines.push(...renderIgnore(update.ignore, indent));
  }

  return lines;
};

/**
 * Renders a complete dependabot.yml. Top-level keys sit at column zero: the
 * previous version indented the whole document and then trimmed the first
 * line, which produced a file that would not parse.
 */
export const renderConfig = (updates: UpdateBlock[]): string => {
  if (updates.length === 0) {
    return "";
  }

  const lines = ["version: 2", "updates:"];
  updates.forEach((update) => lines.push(...renderUpdate(update)));

  return lines.join("\n") + "\n";
};

const parseLimit = (value: string): number | undefined => {
  const trimmed = value.trim();
  if (trimmed === "") {
    return undefined;
  }

  const parsed = Number(trimmed);
  return Number.isInteger(parsed) && parsed >= 0 ? parsed : undefined;
};

const noiseControls = () => ({
  groups: groupsFor(configStore.grouping),
  ignore: configStore.ignoreMajor ? ignoreMajorUpdates() : undefined,
});

/** Builds the update blocks for whichever mode the user is in. */
export const buildUpdates = (): UpdateBlock[] => {
  if (configStore.useDetected && configStore.detected.length > 0) {
    const interval = configStore.schedule || "weekly";
    const limit = parseLimit(configStore.openPrLimit);

    return configStore.detected.map((item) => ({
      packageEcosystem: item.ecosystem,
      directory: item.directory,
      interval,
      reviewers: [configStore.reviewer],
      openPullRequestsLimit: limit,
      labels: [configStore.label],
      ...noiseControls(),
    }));
  }

  if (!configStore.ecosystem || !configStore.directory || !configStore.schedule) {
    return [];
  }

  return [{
    packageEcosystem: configStore.ecosystem,
    directory: configStore.directory,
    interval: configStore.schedule,
    reviewers: [configStore.reviewer],
    openPullRequestsLimit: parseLimit(configStore.openPrLimit),
    labels: [configStore.label],
    ...noiseControls(),
  }];
};

export const generateYaml = () => {
  configStore.saveYamlData(renderConfig(buildUpdates()));
};
