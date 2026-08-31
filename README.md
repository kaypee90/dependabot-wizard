# Dependabot Wizard

**Point it at a repository and it writes a `dependabot.yml` that covers every package manager you
actually have — validated before it touches disk.**

[**Try the wizard in your browser →**](https://kaypee90.github.io/dependabot-wizard/) — no install,
generates the config and hands it to you to copy.

![Depbot in action](https://github.com/kaypee90/dependabot-wizard/blob/main/assets/depbot.png)

## Why use it

Writing a `dependabot.yml` for a single-language repo is a copy-paste from the GitHub docs. Depbot is
for the cases where it stops being that:

- **It finds the ecosystems for you.** `depbot` walks the project and reports every manifest it
  finds — `go.mod` at the root, `package.json` under `ui/web`, `requirements.txt` under
  `services/api`, workflows under `.github/workflows` — and builds one update block per ecosystem, at
  the right directory. In a monorepo this is the part people get wrong, because it means knowing every
  manifest in a tree nobody has read end to end. Dependency directories such as `node_modules` and
  `vendor` are skipped, so you get real manifests rather than your dependencies' manifests.
- **It refuses to write an invalid config.** Every config is parsed back and checked — version,
  ecosystem names, directory format, schedule interval, no empty reviewer or label entries — before
  anything reaches disk. An unparseable file is rejected with the reason, not written.
- **It shows you whether Dependabot is working.** The web UI charts your dependency pull requests, and
  the [GitHub Action](#github-action) posts the same summary on a schedule. That is the question the
  config leads to, and GitHub makes you click through a PR list to answer it.

### Before and after

A monorepo with a Go service, a React app and CI workflows. Written by hand, this means finding every
manifest in the tree and remembering that `github-actions` is declared against the root regardless of
where the workflows live. `depbot` scans and produces all of it:

| By hand | With depbot |
| --- | --- |
| Read the docs, find every manifest, write ~30 lines, hope you did not miss a directory | `depbot`, answer two questions |
| Mistakes surface when Dependabot silently does not run | Rejected before the file is written, with the reason |
| One ecosystem is easy to miss; nothing tells you | Every manifest reported with the file that proved it |

## Install

Pick whichever suits you — the hosted wizard needs nothing at all.

```bash
# Go toolchain
go install github.com/kaypee90/dependabot-wizard/cmd/depbot@latest

# Prebuilt binary (see the releases page for your platform)
curl -sSL https://github.com/kaypee90/dependabot-wizard/releases/latest/download/depbot_linux_amd64.tar.gz | tar xz

# Docker (published on tagged releases)
docker run --rm -e GITHUB_TOKEN ghcr.io/kaypee90/dependabot-wizard --report --owner OWNER --repo REPO
```

The binary embeds the web UI, so there is nothing else to install and no Node toolchain at runtime.
If you used `go install`, make sure Go's bin directory is on your `PATH`:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

## Usage

Run it from the root of the project you want to configure:

```bash
cd <your-project>
depbot
```

It scans first and offers what it found:

```
Found 5 package ecosystems in this project:
  docker          /            (from Dockerfile)
  github-actions  /            (from .github/workflows/ci.yml)
  gomod           /            (from go.mod)
  npm             /ui/web      (from ui/web/package.json)
  pip             /services/api (from services/api/requirements.txt)

Use this detected configuration? [yes/no]
```

Answer `yes` and you are asked only for the schedule and pull request limit. Answer `no` and you get
the original question-by-question flow, one ecosystem at a time.

The result is written to `.github/dependabot.yml`, creating the directory if needed. Running inside
`.github` writes alongside you instead.

### Web UI

```bash
depbot --web
```

Opens a browser at `http://127.0.0.1:3001/`, with a **Scan project** button, a form, a live YAML
preview, and a **Pull Requests** dashboard.

| Flag | Meaning |
| --- | --- |
| `--help`, `-h` | Show help |
| `--version`, `-v` | Print the version |
| `--web`, `-w` | Start the web application |
| `--port` | Port for the web application (default `3001`) |
| `--report` | Print a dependabot health report |
| `--owner`, `--repo`, `--label` | Target for `--report` |

`DEPBOT_NO_BROWSER=1` starts the server without opening a browser, for remote or headless use.
`--report` reads `GITHUB_TOKEN`.

The dashboard needs a GitHub token with `repo` scope. It is held in memory for the request only:
never written to disk, and never logged.

## What it writes

```yaml
version: 2
updates:
- package-ecosystem: gomod
  directory: /
  schedule:
    interval: weekly
  open-pull-requests-limit: 5
- package-ecosystem: npm
  directory: /ui/web
  schedule:
    interval: weekly
  open-pull-requests-limit: 5
- package-ecosystem: github-actions
  directory: /
  schedule:
    interval: weekly
  open-pull-requests-limit: 5
```

## GitHub Action

Config is the beginning; the question that follows is whether anyone is merging what Dependabot
opens. This action answers it on a schedule and writes to the workflow summary.

```yaml
name: Dependabot health

on:
  schedule:
    - cron: '0 9 * * 1'   # Monday morning
  workflow_dispatch:

jobs:
  report:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      pull-requests: read
    steps:
      - uses: kaypee90/dependabot-wizard@v1
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
```

It produces:

> ## Dependabot health: kaypee90/zebrok
>
> | Metric | Value |
> | --- | ---: |
> | Total pull requests | 8 |
> | Open | 3 |
> | Merged | 4 |
> | Closed without merging | 1 |
> | Merge rate | 80% |
> | Oldest open | [#61](https://github.com/kaypee90/zebrok/pull/61), 34 days |
>
> 3 open, 80% of resolved pull requests merged. This repository is keeping up.

The closing line changes with the numbers: a queue nobody is reviewing, a merge rate low enough to
suggest `ignore` rules, or a backlog large enough to suggest `groups`.

### Inputs and outputs

| Input | Default | Meaning |
| --- | --- | --- |
| `token` | `github.token` | Token with repo scope |
| `owner` | current repository | Repository owner |
| `repo` | current repository | Repository name |
| `label` | `dependencies` | Pull request label to count |

Outputs `total`, `open`, `merged`, `closed` and `merge-rate`, so you can gate on them:

```yaml
      - uses: kaypee90/dependabot-wizard@v1
        id: health
        with:
          token: ${{ secrets.GITHUB_TOKEN }}
      - if: steps.health.outputs.open > 20
        run: echo "::warning::${{ steps.health.outputs.open }} dependabot PRs are open"
```

## Development

Requires Go 1.22+ and, only if you are changing the UI, Node 20+.

```bash
make deps        # Go and npm dependencies
make build       # build ./depbot with the version stamped in
make run         # run the CLI wizard from source
make web         # run the web application from source
```

### Testing

```bash
make test        # everything: Go tests, UI tests, embedded asset check
make test-go     # Go only
make test-ui     # React only
make typecheck   # TypeScript, no emit
```

To exercise the running server end to end:

```bash
make build
cd /some/project
DEPBOT_NO_BROWSER=1 /path/to/depbot --web --port 3999 &

curl http://127.0.0.1:3999/healthz
curl http://127.0.0.1:3999/api/detect

# A config that does not parse is rejected, and nothing is written
curl -i -X POST http://127.0.0.1:3999/api/configurations \
  -H 'Content-Type: application/json' \
  -d '{"configuration":"version: 2\n  updates:\n  - package-ecosystem: npm\n"}'

# A valid one returns 201 and the path it wrote
curl -i -X POST http://127.0.0.1:3999/api/configurations \
  -H 'Content-Type: application/json' \
  -d '{"configuration":"version: 2\nupdates:\n- package-ecosystem: npm\n  directory: /\n  schedule:\n    interval: weekly\n"}'
```

To check the action image the way CI does:

```bash
docker build -t depbot-action .
docker run --rm depbot-action --version
docker run --rm depbot-action --report --owner o --repo r      # exits 1: no token
GITHUB_TOKEN=<token> docker run --rm -e GITHUB_TOKEN depbot-action --report --owner OWNER --repo REPO
```

To preview the static wizard exactly as GitHub Pages serves it:

```bash
cd ui/depbot-app
REACT_APP_STATIC=true PUBLIC_URL=/dependabot-wizard npm run build:static
mkdir -p /tmp/pages/dependabot-wizard && cp -R build/. /tmp/pages/dependabot-wizard/
cd /tmp/pages && python3 -m http.server 4111
# open http://127.0.0.1:4111/dependabot-wizard/
```

### Changing the UI

The React app lives in `ui/depbot-app` and its production build is committed to
`cmd/depbot/static`, which the binary embeds — that is what lets `go install` work without Node.
After any UI change:

```bash
make ui          # rebuild and re-embed
make verify-assets
```

`make verify-assets` checks that `index.html` only references files that exist in the embedded tree
and that exactly one bundle is present. CI runs it on every push.

The same source builds two ways. `npm run build` produces the copy the binary serves, which can scan
the filesystem and write the file. `npm run build:static` sets `REACT_APP_STATIC`, which drops those
server-backed features and offers the config for copying instead — that is what is published to
Pages.

### Releasing

Tag and push; `.github/workflows/release.yml` runs GoReleaser and publishes binaries for Linux, macOS
and Windows on amd64 and arm64.

```bash
git tag v1.2.0 && git push origin v1.2.0
```

Run `goreleaser check` after editing `.goreleaser.yaml`, and `goreleaser release --snapshot --clean`
to build artifacts locally without publishing. A Homebrew tap is configured but commented out; the
file explains what to create before enabling it.

## How it fits together

```
action.yml           Marketplace action definition
Dockerfile           image the action runs from
.goreleaser.yaml     release build matrix
cmd/depbot/
  main.go      flags, CLI wizard, web launcher
  detect.go    manifest scanning -> ecosystems and directories
  models.go    config model, validation, YAML round-trip
  prompts.go   promptui questions
  report.go    dependabot health report
  server.go    HTTP API and embedded UI
  static/      committed React build (embedded)
github/
  graphqlclient.go   GitHub GraphQL client
  models.go          response types
ui/depbot-app/       React source
scripts/             embedded asset verification
```

| Endpoint | Method | Purpose |
| --- | --- | --- |
| `/` | GET | The web UI |
| `/healthz` | GET | Health check |
| `/api/detect` | GET | Ecosystems found in the working directory |
| `/api/configurations` | POST | Validate and write `dependabot.yml` |
| `/api/pull-requests` | POST | Dependency pull requests for a repository |

The server binds to `127.0.0.1` only and sends no CORS headers, so pages in your browser cannot reach
it while it is running.

## License

MIT — see [LICENSE](LICENSE).
