# gh-app-check

[![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/wakeward/gh-app-check/badge)](https://securityscorecards.dev/viewer/?uri=github.com/wakeward/gh-app-check)
[![Go Report Card](https://goreportcard.com/badge/github.com/wakeward/gh-app-check)](https://goreportcard.com/report/github.com/wakeward/gh-app-check)
[![Release](https://img.shields.io/github/v/release/wakeward/gh-app-check)](https://github.com/wakeward/gh-app-check/releases)
[![License](https://img.shields.io/github/license/wakeward/gh-app-check)](LICENSE)

gh-app-check is a `gh` CLI extension and security auditing tool designed to evaluate GitHub App installations across your organization. It ensures apps adhere to the principle of least privilege, identifies toxic permission combinations, and traces the actual execution of app credentials within your codebase.

Unlike static configuration checkers, gh-app-check bridges the gap between the **Control Plane** (what permissions an app has) and the **Execution Plane** (how the app's token is actually generated and used in Actions workflows).

> **Status:** this repository currently contains the secure-foundation scaffolding and a command skeleton. The Control Plane auditor (Phase 1), Execution Plane tracer (Phase 2), and CI/CD SARIF integration (Phase 3) are in active development - see [Implementation Phases](#implementation-phases) below.

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

- **Phase 0 (this repo, today):** secure foundation - LICENSE, SECURITY.md, branch protection, Dependabot, gosec/govulncheck as blocking CI checks, OpenSSF Scorecard workflow, command skeleton.
- **Phase 1:** Control Plane auditor - `GET /orgs/{org}/installations` fetching, pagination, rules engine (blast radius, toxic permissions).
- **Phase 2:** Execution Plane tracer - Code Search API integration, YAML AST parsing of Actions workflows, `.pem`/hardcoded-key detection.
- **Phase 3:** CI/CD integration - `--strict` exit codes, SARIF output for GitHub Advanced Security.

## Security & Supply Chain

This project takes supply chain security seriously.

- Dependency updates go through an explicit 48-hour cooldown window before being proposed (`.github/dependabot.yml`), to avoid pulling in just-published (and potentially compromised) releases.
- `gosec` and `govulncheck` run as blocking checks on every pull request.
- OpenSSF Scorecard runs on every push, pull request, and branch protection change.
- Future releases will be signed keylessly using **Sigstore/Cosign**, include an SPDX Software Bill of Materials (SBOM), and generate **SLSA Level 3** provenance attestations.

Please refer to [SECURITY.md](SECURITY.md) for vulnerability reporting guidelines.
