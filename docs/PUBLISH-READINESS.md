# Publish readiness checklist

Living checklist for taking `gh-app-check` and sibling `gh-app-graph` public.
See also [`GITHUB_SETUP_RUNBOOK.md`](GITHUB_SETUP_RUNBOOK.md) for GitHub settings.

## Phase A - Policy, scrub, OSS hygiene (done)

- [x] `organization-takeover` single-grant policy decided (see gh-app-graph `docs/calibration-notes.md`)
- [x] Doc accuracy: trace, SARIF, `notable_grants`, graph README typo, `evaluate-app` stub
- [x] Sanitize operator-specific paths in setup runbook
- [x] `CONTRIBUTING.md`, issue templates, graph `CODEOWNERS`
- [x] Generalize internal validation notes in `improvement-backlog.md`
- [x] Mark `ORG-VALIDATION-RUNBOOK.md` as maintainer-only template

## Phase B - Release prep (done)

- [x] Fix CI graph checkout (composite action while graph was private)
- [x] Bump Go to 1.26.6 (stdlib govulncheck fixes on graph)
- [x] Pin graph in `gh-app-check` `go.mod` at v0.1.1; drop `replace`
- [x] Drop `GH_GRAPH_READ_TOKEN` from CI (graph public; use `go mod download`)
- [x] GoReleaser + signed release workflow (Cosign, SBOM)
- [x] Signed GPG tags on graph releases (`v0.1.0`, `v0.1.1`)
- [x] Branch protection + secret scanning + CodeQL default setup (see setup runbook)
- [x] Extend CI path filters to `data/**` (graph); add `go test` to `refresh.yml`
- [x] CodeQL default setup (custom workflow removed)

## Phase C - History reset (done)

- [x] Orphan squash to single commit per repo
- [x] GPG-signed tags on squashed commits (`gh-app-graph` v0.1.0 / v0.1.1)
- [x] Mirror backups under operator home directory

## Phase D - Publish (in progress)

- [x] Flip repo visibility public
- [ ] Merge finish-publish PR (drop private-module CI, token-permissions on release.yml)
- [ ] Cut `gh-app-check` v0.1.0 release (signed tag push → GoReleaser binaries)
- [ ] Smoke test: `gh extension install wakeward/gh-app-check`
- [ ] Verify OpenSSF Scorecard badge (expect low solo scores; see runbook)
- [ ] Announce with DRAFT catalog disclaimer (graph)

## Sensitive-data audit (2026-08-24)

Git history search found **no** customer org names, app slugs, or tokens.
Orphan squash removed generic Cursor co-author trailers from public history.

## OpenSSF Scorecard target (solo maintainer)

Realistic overall: **7.5-8.5** after publish. See setup runbook for branch
protection tiers. `Code-Review` will score low on a solo project; do not fake
with a second account.
