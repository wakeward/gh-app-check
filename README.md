# gh-app-check

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/wakeward/gh-app-check/badge)](https://securityscorecards.dev/viewer/?uri=github.com/wakeward/gh-app-check)
[![Go Report Card](https://goreportcard.com/badge/github.com/wakeward/gh-app-check)](https://goreportcard.com/report/github.com/wakeward/gh-app-check)
[![Release](https://img.shields.io/github/v/release/wakeward/gh-app-check)](https://github.com/wakeward/gh-app-check/releases)
[![License](https://img.shields.io/github/license/wakeward/gh-app-check)](LICENSE)

gh-app-check is a `gh` CLI extension and security auditing tool designed to evaluate GitHub App installations across your organization. It ensures apps adhere to the principle of least privilege, identifies toxic permission combinations, and traces the actual execution of app credentials within your codebase.

Unlike static configuration checkers, gh-app-check bridges the gap between the **Control Plane** (what permissions an app has) and the **Execution Plane** (how the app's token is actually generated and used in Actions workflows).

> **Status:** Phase 1 (`gh app-check org`) is implemented: org installation fetch, least-privilege rules, and toxic-combination checks via embedded data from [`gh-app-graph`](https://github.com/wakeward/gh-app-graph). Phase 1 assesses **installed App permissions only** - not user roles or install authorization paths (see [gh-app-graph installation gates](https://github.com/wakeward/gh-app-graph/blob/main/docs/installation-gates.md) for likelihood context in blog/threat-model material). Phase 2 (execution trace), Phase 3 (SARIF), and Phase 4 (drift guard) are not yet implemented - see [Implementation Phases](#implementation-phases) below.

## Documentation

- [**Installation, Deployment, & Troubleshooting Guide**](docs/INSTALLATION.md): learn how to install the CLI locally, deploy it as an automated enterprise CI/CD job, and resolve common API or rate-limit errors.

## Quick Start

### 1. Control Plane Audit (Default)

Audit all applications installed in your organization to identify blast radius issues (e.g. "All Repositories" access) and toxic write permissions.

```bash
gh app-check org my-organization
```

### 2. Execution Plane Trace

Trace exactly how a specific internal app is being used across your codebase. This checks your `.github/workflows` to ensure tokens are being scoped correctly.

```bash
gh app-check trace my-internal-deployer-app --org my-organization
```

### Output Formats

By default, the tool outputs a human-readable terminal table. You can modify this for automation pipelines:

- `--format table` (default)
- `--format json`: structured JSON output for piping to `jq` or external SIEMs.
- `--format markdown`: Markdown tables for automated PR comments or Issues.

## Implementation Phases

- **Phase 0:** secure foundation - LICENSE, SECURITY.md, branch protection, Dependabot, gosec/govulncheck as blocking CI checks, OpenSSF Scorecard workflow, command skeleton.
- **Phase 1 (implemented):** Control Plane auditor - `GET /orgs/{org}/installations` fetching, pagination, rules engine (blast radius, toxic permissions from `gh-app-graph`).
- **Phase 2:** Execution Plane tracer - Code Search API integration, YAML AST parsing of Actions workflows, `.pem`/hardcoded-key detection.
- **Phase 3:** CI/CD integration - `--strict` exit codes, SARIF output for GitHub Advanced Security.
- **Phase 4 (backlog, not yet scoped):** permission drift guard - track how installed apps' permissions change over time and alert on high-risk escalations, inspired by [google/capslock](https://github.com/google/capslock)'s capability-diffing approach for Go packages. See [docs/BACKLOG.md](docs/BACKLOG.md) for open questions.

## Development

Clone **`gh-app-graph`** and **`gh-app-check`** as sibling directories. Local builds
use a `replace` directive in `go.mod`:

```go
replace github.com/wakeward/gh-app-graph => ../gh-app-graph
```

```bash
go build ./...
go vet ./...
go test ./...
```

Product improvement backlog: [`docs/improvement-backlog.md`](docs/improvement-backlog.md).
Local org audit notes belong in `docs/FEEDBACK-*.md` (gitignored).

## Security & Supply Chain

This project takes supply chain security seriously.

- Dependency updates go through an explicit 48-hour cooldown window before being proposed (`.github/dependabot.yml`), to avoid pulling in just-published (and potentially compromised) releases.
- `gosec` and `govulncheck` run as blocking checks on every pull request.
- OpenSSF Scorecard runs on every push, pull request, and branch protection change.
- Future releases will be signed keylessly using **Sigstore/Cosign**, include an SPDX Software Bill of Materials (SBOM), and generate **SLSA Level 3** provenance attestations.

Please refer to [SECURITY.md](SECURITY.md) for vulnerability reporting guidelines.
