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

## Phase B - Release prep (private)

- [x] Fix CI graph checkout (composite action; requires `GH_GRAPH_READ_TOKEN` while graph is private)
- [x] Bump Go to 1.26.6 (stdlib govulncheck fixes on graph)
- [x] Pin graph in `gh-app-check` `go.mod` at v0.1.0; drop `replace`
- [ ] Drop `GH_GRAPH_READ_TOKEN` secret after graph is public
- [x] GoReleaser + signed release workflow (Cosign, SBOM)
- [ ] Signed GPG tags on all releases (see setup runbook)
- [ ] Apply branch protection + secret scanning (see setup runbook; after public)
- [x] Extend CI path filters to `data/**` (graph); add `go test` to `refresh.yml`
- [x] CodeQL analysis workflow

## Phase C - History reset

- [ ] Scrub working tree (this checklist + FEEDBACK gitignore verified)
- [ ] Orphan squash → single commit per repo (`Initial public release v0.1.0`)
- [ ] Re-create **GPG-signed** `v0.1.0` tags on squashed commits; force-push tags
- [ ] New public remote or force-push (explicit approval required)

## Phase D - Publish

- [ ] Flip repo visibility public
- [ ] Cut `gh-app-check` v0.1.0 release (binary for `gh extension install`)
- [ ] Verify OpenSSF Scorecard badge
- [ ] Announce with DRAFT catalog disclaimer (graph)

## Sensitive-data audit (2026-08-24)

Git history search found **no** customer org names, app slugs, or tokens.
Generic validation metadata and Cursor co-author trailers exist in history;
orphan squash removes both.

## OpenSSF Scorecard target (solo maintainer)

Realistic overall: **7.5-8.5** after publish. See setup runbook for branch
protection tiers. `Code-Review` will score low on a solo project; do not fake
with a second account.
