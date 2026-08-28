# Contributing to gh-app-check

Thank you for helping improve GitHub App installation auditing for the
community.

## Scope

This repo implements the **CLI and org scan** (`gh app-check org`). Permission
catalog, toxic combinations, and methodology live in
[`gh-app-graph`](https://github.com/wakeward/gh-app-graph). Open catalog or
methodology issues in that repo unless the bug is in how this tool fetches or
displays results.

## Getting started

1. Clone **`gh-app-graph`** and **`gh-app-check`** as sibling directories.
2. `go test ./...`, `go vet ./...`, `go build -o gh-app-check .`
3. Local development uses a `replace` directive in `go.mod` for the sibling
   graph checkout.

## Pull requests

- Open a PR against `main` from a branch (direct pushes to `main` are
  blocked once branch protection is enabled).
- Keep changes focused; match existing style and test patterns.
- Run `go test ./...` locally before pushing.
- Do not commit org-specific audit output, tokens, or customer identifiers.
  Use synthetic fixtures in tests only.

## Reporting security issues

See [SECURITY.md](SECURITY.md). Do **not** open public issues for
vulnerabilities.

## Maintainer note

This project is primarily maintained by a solo maintainer. PRs and issues are
welcome; response times may vary.
