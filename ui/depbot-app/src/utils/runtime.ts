/**
 * The wizard ships in two shapes from one codebase:
 *
 *  - served by the depbot binary, where it can scan the project on disk and
 *    write .github/dependabot.yml through the local API;
 *  - published as a static page, where there is no server, so it generates
 *    the config for copying instead.
 *
 * REACT_APP_STATIC is set at build time by the Pages workflow. It is read
 * through a function rather than a constant so it can be exercised in tests.
 */
export const isStaticBuild = (): boolean => process.env.REACT_APP_STATIC === 'true';

/** Where to send someone who needs the features only the CLI can provide. */
export const CLI_INSTALL_COMMAND = 'go install github.com/kaypee90/dependabot-wizard/cmd/depbot@latest';
