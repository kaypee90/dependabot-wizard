import { renderConfig, UpdateBlock } from './config-generator';

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
