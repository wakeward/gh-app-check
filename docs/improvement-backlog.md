# Control-plane auditor improvement backlog

Sanitized product backlog derived from Phase 1 field validation. **Do not** commit
org-specific feedback files (`docs/FEEDBACK-*.md` is gitignored).

Use synthetic fixtures in tests only.

## Done (2026-08-14)

- [x] **P0** Surface `near_misses` from `gh-app-graph` in JSON/table/markdown
- [x] **P0** All-repos + read-only grants → WARN (writes still → HIGH)
- [x] **P1** God-mode write count → HIGH (not CRITICAL); `write_scope_count` in output
- [x] **P1** Sort by risk; table/markdown columns for WRITES, TOXIC, NEAR_MISSES
- [x] **P2** `--timeout` on `org` command
- [x] **P2** `PermissionsMap` returns error instead of silent empty map
- [x] **Docs** OAuth scope vs org role; trace sections marked Phase 2 planned
- [x] **Hygiene** `go.mod` replace path → `../gh-app-graph`

## Remaining

### P2 - Friendly app name

Installation list API exposes `app_slug` but not display name; table duplicates slug.

**Shape:** Optional `GET /apps/{slug}` enrichment.

### CI / supply chain (non-blockers)

| Item | Action |
| --- | --- |
| `gofmt` / import order | Keep CI clean on touched files |
| Action comments | Pin helper actions with `# vX.Y.Z` matching SHA |
| Harden-runner | Consider `egress-policy: block` when ready |

## Non-goals

- Org-specific permission maps or customer identifiers in-repo
- Treating PASS as "safe" (Phase 1 is control-plane only)
