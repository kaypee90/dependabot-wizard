# Dependabot Wizard

**Point it at a repository and it writes a `dependabot.yml` that covers every package manager you
actually have — validated before it touches disk.**

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
- **It shows you whether Dependabot is working.** The web UI reads your repository's dependency pull
  requests and charts them: how many are open, closed and merged, and how they arrived over time.
  That is the question the config leads to, and GitHub makes you click through a PR list to answer it.

Two interfaces over the same engine: a terminal wizard, and a local web UI with a live preview of the
file it will write.

## Install

```bash
go install github.com/kaypee90/dependabot-wizard/cmd/depbot@latest
```

The binary embeds the web UI, so there is nothing else to install and no Node toolchain required at
runtime. Make sure Go's bin directory is on your `PATH`:

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

`DEPBOT_NO_BROWSER=1` starts the server without opening a browser, for remote or headless use.

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

## How it fits together

```
cmd/depbot/
  main.go      flags, CLI wizard, web launcher
  detect.go    manifest scanning -> ecosystems and directories
  models.go    config model, validation, YAML round-trip
  prompts.go   promptui questions
  server.go    HTTP API and embedded UI
  static/      committed React build (embedded)
github/
  graphqlclient.go   GitHub GraphQL client
  models.go          response types
ui/depbot-app/       React source
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
