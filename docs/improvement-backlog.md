# Control-plane auditor improvement backlog

Sanitized product backlog derived from Phase 1 field validation. **Do not** commit
org-specific feedback files (`docs/FEEDBACK-*.md` is gitignored).

Use synthetic fixtures in tests only.

## Done (2026-08-14 through 2026-08-21)

- [x] **P0** Surface `near_misses` from `gh-app-graph` in JSON/table/markdown
- [x] **P0** All-repos + read-only grants → WARN (writes still → HIGH)
- [x] **P1** God-mode write count → HIGH (not CRITICAL); `write_scope_count` in output
- [x] **P1** Sort by risk; table/markdown columns for WRITES, TOXIC, NEAR_MISSES
- [x] **P2** `--timeout` on `org` command
- [x] **P2** `PermissionsMap` returns error instead of silent empty map
- [x] **Docs** OAuth scope vs org role; trace sections marked Phase 2 planned
- [x] **Hygiene** `go.mod` replace path → `../gh-app-graph`

- [x] **P2** Friendly app name via `GET /apps/{slug}` (`--no-enrich-names` to skip)
- [x] **Platform** GHES-only rule filtering and scope highlighting

## Remaining

### Phase 2 trace priorities (from ecosystem research)

| Item | Source pattern | Notes |
| --- | --- | --- |
| Detect `pull_request_target` + `create-github-app-token` + fork checkout | `credential-access-prt-fork-iatt-exfiltration` | GHSA-9g93-rxr5-xhqw |
| Detect permissive `[bot]` actor checks in agent workflows | `execution-ai-agent-external-bot-trust` | Allow-list by App slug / installation ID |
| Document Zizmor as complementary workflow linter | marketplace-trust-limitations.md | Not a gh-app-check dependency |

### CI / supply chain (non-blockers)

| Item | Action |
| --- | --- |
| `gofmt` / import order | Keep CI clean on touched files |
| Action comments | Pin helper actions with `# vX.Y.Z` matching SHA |
| Harden-runner | Consider `egress-policy: block` when ready |

## Non-goals

- Org-specific permission maps or customer identifiers in-repo
- Treating PASS as "safe" (Phase 1 is control-plane only)
- User install authorization inference or likelihood scoring (see gh-app-graph `installation-gates.md`)
- Vendor OAuth callback or platform CVE assessment in org scan output
