import { configStore } from "../store/configuration-store";

export type UpdateBlock = {
  packageEcosystem: string;
  directory: string;
  interval: string;
  reviewers: string[];
  openPullRequestsLimit?: number;
  labels: string[];
};

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
  }];
};

export const generateYaml = () => {
  configStore.saveYamlData(renderConfig(buildUpdates()));
};
