# gh-app-check

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/wakeward/gh-app-check/badge)](https://securityscorecards.dev/viewer/?uri=github.com/wakeward/gh-app-check)
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
gh app-check org my-organization --explain --no-near-misses
gh app-check org my-organization --explain-all --format json
```

### 2. Execution Plane Trace (Phase 2 - not yet implemented)

> **Not available yet.** `gh app-check trace` returns an error today. Phase 2 will
> add Code Search and workflow YAML analysis. Use Phase 1 (`org`) for control-plane
> auditing now.

```bash
# Planned - do not rely on this until Phase 2 ships
gh app-check trace my-internal-deployer-app --org my-organization
```

### Output Formats

By default, the tool outputs a human-readable terminal table. You can modify this for automation pipelines:

- `--format table` (default)
- `--format json`: structured JSON output for piping to `jq` or external SIEMs. Includes `exploit_path`, `matched_grants`, and `control_plane_findings` on each installation. Field `notable_grants` is populated only when `--explain` or `--explain-all` is used (or with `--explain-all` in JSON-only runs).
- `--format markdown`: Markdown tables for automated PR comments or Issues.
- `--explain`: human-readable narrative for **CRITICAL/HIGH** findings (why each toxic combo and structural rule fired).
- `--explain-all`: with `--explain`, also include PASS/WARN installations and standalone permission notes from the catalog.

## Implementation Phases

- **Phase 0:** secure foundation - LICENSE, SECURITY.md, branch protection, Dependabot, gosec/govulncheck as blocking CI checks, OpenSSF Scorecard workflow, command skeleton.
- **Phase 1 (implemented):** Control Plane auditor - `GET /orgs/{org}/installations` fetching, pagination, rules engine (blast radius, toxic permissions from `gh-app-graph`).
- **Phase 2 (planned):** Execution Plane tracer - Code Search API integration, YAML AST parsing of Actions workflows, `.pem`/hardcoded-key detection.
- **Phase 3 (planned):** CI/CD integration - `--strict` exit codes, SARIF output for GitHub Advanced Security.
- **Phase 4 (backlog, not yet scoped):** permission drift guard - track how installed apps' permissions change over time and alert on high-risk escalations, inspired by [google/capslock](https://github.com/google/capslock)'s capability-diffing approach for Go packages. See [docs/BACKLOG.md](docs/BACKLOG.md) for open questions.

## Development

Clone **`gh-app-graph`** and **`gh-app-check`** as sibling directories for local
catalog work. CI consumes a tagged module (`github.com/wakeward/gh-app-graph`).

**Local dev with sibling checkout** (optional):

```go
replace github.com/wakeward/gh-app-graph => ../gh-app-graph
```

```bash
go build -o gh-app-check .
go vet ./...
go test ./...
```

Without a `replace` line, `go mod download` fetches `github.com/wakeward/gh-app-graph`
from the public module proxy.

Product improvement backlog: [`docs/improvement-backlog.md`](docs/improvement-backlog.md).
Publish checklist: [`docs/PUBLISH-READINESS.md`](docs/PUBLISH-READINESS.md).
Maintainer org validation template (read-only): [`docs/ORG-VALIDATION-RUNBOOK.md`](docs/ORG-VALIDATION-RUNBOOK.md).
Local org audit notes belong in `docs/FEEDBACK-*.md` (gitignored).

## Security & Supply Chain

This project takes supply chain security seriously.

- Dependency updates go through an explicit 48-hour cooldown window before being proposed (`.github/dependabot.yml`), to avoid pulling in just-published (and potentially compromised) releases.
- `gosec` and `govulncheck` run as blocking checks on every pull request.
- OpenSSF Scorecard runs on every push, pull request, and branch protection change.
- Releases are signed with GPG tags; release artifacts use **Sigstore/Cosign** with SPDX SBOM (see [releases](https://github.com/wakeward/gh-app-check/releases)).

Please refer to [SECURITY.md](SECURITY.md) for vulnerability reporting guidelines.
Contributing: [CONTRIBUTING.md](CONTRIBUTING.md).
