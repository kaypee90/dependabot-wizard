import { renderConfig, UpdateBlock, groupsFor, ignoreMajorUpdates } from './config-generator';

const baseUpdate: UpdateBlock = {
  packageEcosystem: 'npm',
  directory: '/',
  interval: 'weekly',
  reviewers: [],
  openPullRequestsLimit: 5,
  labels: [],
};

describe('renderConfig', () => {
  // The bug that shipped: the whole document was indented two spaces and then
  // trimmed, leaving "version: 2" at column zero and "  updates:" indented
  // under it, which does not parse.
  it('puts top-level keys at column zero', () => {
    const lines = renderConfig([baseUpdate]).split('\n');

    expect(lines[0]).toBe('version: 2');
    expect(lines[1]).toBe('updates:');
  });

  it('renders a complete single-ecosystem config', () => {
    expect(renderConfig([{ ...baseUpdate, reviewers: ['kaypee90'], labels: ['dependencies'] }])).toBe(
      [
        'version: 2',
        'updates:',
        '- package-ecosystem: npm',
        '  directory: /',
        '  schedule:',
        '    interval: weekly',
        '  reviewers:',
        '  - kaypee90',
        '  open-pull-requests-limit: 5',
        '  labels:',
        '  - dependencies',
        '',
      ].join('\n'),
    );
  });

  it('omits reviewers and labels when they are blank', () => {
    const rendered = renderConfig([{ ...baseUpdate, reviewers: [''], labels: ['   '] }]);

    expect(rendered).not.toContain('reviewers');
    expect(rendered).not.toContain('labels');
    expect(rendered).not.toContain('- \n');
  });

  it('omits the pull request limit when it is not set', () => {
    const rendered = renderConfig([{ ...baseUpdate, openPullRequestsLimit: undefined }]);

    expect(rendered).not.toContain('open-pull-requests-limit');
  });

  it('renders every ecosystem when several are detected', () => {
    const rendered = renderConfig([
      baseUpdate,
      { ...baseUpdate, packageEcosystem: 'gomod' },
      { ...baseUpdate, packageEcosystem: 'github-actions' },
    ]);

    expect(rendered.match(/^- package-ecosystem:/gm)).toHaveLength(3);
    expect(rendered).toContain('- package-ecosystem: gomod');
  });

  it('quotes values that would otherwise change meaning', () => {
    const rendered = renderConfig([{ ...baseUpdate, directory: '/apps: web' }]);

    expect(rendered).toContain('directory: "/apps: web"');
  });

  it('returns nothing when there is no update to render', () => {
    expect(renderConfig([])).toBe('');
  });
});

describe('grouping and ignore rules', () => {
  it('renders the all-minor-and-patch group', () => {
    const rendered = renderConfig([{ ...baseUpdate, groups: groupsFor('all') }]);

    expect(rendered).toContain([
      '  groups:',
      '    all-minor-and-patch:',
      '      patterns:',
      '      - "*"',
      '      update-types:',
      '      - minor',
      '      - patch',
    ].join('\n'));
  });

  it('renders production and development groups in a stable order', () => {
    const rendered = renderConfig([{ ...baseUpdate, groups: groupsFor('by-dependency-type') }]);

    // Sorted by name, so the output does not churn between renders.
    expect(rendered.indexOf('development-dependencies')).toBeLessThan(
      rendered.indexOf('production-dependencies'),
    );
    expect(rendered).toContain('      dependency-type: production');
    expect(rendered).toContain('      dependency-type: development');
  });

  it('renders nothing for the none strategy', () => {
    const rendered = renderConfig([{ ...baseUpdate, groups: groupsFor('none') }]);

    expect(rendered).not.toContain('groups:');
  });

  it('renders the ignore rule for major versions', () => {
    const rendered = renderConfig([{ ...baseUpdate, ignore: ignoreMajorUpdates() }]);

    expect(rendered).toContain([
      '  ignore:',
      '  - dependency-name: "*"',
      '    update-types:',
      '    - version-update:semver-major',
    ].join('\n'));
  });

  it('places groups and ignore after the simple keys', () => {
    const rendered = renderConfig([{
      ...baseUpdate,
      labels: ['deps'],
      groups: groupsFor('all'),
      ignore: ignoreMajorUpdates(),
    }]);

    expect(rendered.indexOf('labels:')).toBeLessThan(rendered.indexOf('groups:'));
    expect(rendered.indexOf('groups:')).toBeLessThan(rendered.indexOf('ignore:'));
  });

  it('applies the same rules to every detected ecosystem', () => {
    const rendered = renderConfig([
      { ...baseUpdate, groups: groupsFor('all') },
      { ...baseUpdate, packageEcosystem: 'gomod', groups: groupsFor('all') },
    ]);

    expect(rendered.match(/all-minor-and-patch:/g)).toHaveLength(2);
  });
});
