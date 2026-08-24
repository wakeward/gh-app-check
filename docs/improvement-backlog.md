# Control-plane auditor improvement backlog

Sanitized product backlog derived from Phase 1 field validation. **Do not** commit
org-specific feedback files (`docs/FEEDBACK-*.md` is gitignored).

Use synthetic fixtures in tests only.

## Done (2026-08-14 through 2026-08-24)

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
- [x] **Validation runbook** [`ORG-VALIDATION-RUNBOOK.md`](ORG-VALIDATION-RUNBOOK.md)
- [x] **2026-08-24 validation** First company org read-only run (16 installs, cloud)
- [x] **JSON** Emit `installation_id`, `app_id`, `permissions`, `html_url` per install
- [x] **Output** `--no-near-misses`; dedupe near-miss display; toxic technique `[id]` when duplicated
- [x] **Docs** Local `go build -o gh-app-check .`; 404 vs org role vs `admin:org` false lead
- [x] **UX** Warn when `--platform ghes` used against github.com
- [x] **Explain** `--explain` / `--explain-all`; JSON includes exploit_path and matched_grants

## Validation summary (2026-08-24, sanitized)

### First run (pre-calibration)

| Area | Result |
| --- | --- |
| Core scan (admin org) | PASS - mapping REST → scan correct on spot checks |
| Parent org (member only) | Expected 404 - role, not scope |
| Platform auto/cloud | PASS - 1 GHES toxic rule excluded |
| Name enrichment | PASS - private Apps 404 slug lookup, slug-only label |
| Signal: admin write, all-repos split, god-mode | PASS - matches intuition |
| Signal: `contents:write` alone → Critical toxic | **Too hot** - catalog tuning needed |
| Signal: `organization_administration:write` alone → org takeover | **Coarse** - GitHub bundles rulesets |
| Near misses on `contents:read` | Noisy for default table |
| Table UX | Unusable wide with near misses; JSON preferred |

### Retest (post-calibration + --explain)

| Area | Result |
| --- | --- |
| PR bots (`contents:write` only) | **PASS** - was false CRITICAL, now PASS |
| Org-wide CI App | **PASS** - CRITICAL with multi-grant toxics + structural |
| IT/runner App | **PASS** - CRITICAL via admin write + runners toxic |
| `--explain` / JSON fields | **PASS** - matched grants, exploit_path, structural rationales |
| Performance (16 installs) | **PASS** - ~4s |
| Risk mix | CRITICAL 8→6, PASS 5→7 |

**Verdict:** contents:write calibration validated. Pin v0.1.0 reasonable if org-admin single-grant Critical accepted. Drop `replace` only from a real git checkout with tagged graph.

## Remaining

### P0 - Catalog calibration (gh-app-graph)

| Item | Status |
| --- | --- |
| `contents: write` standalone → High | Done (2026-08-24) |
| Remove single-grant supply-chain Critical toxics | Done |
| `organization-takeover` single-grant | Open - see calibration-notes |

### P1 - Output / UX

| Item | Notes |
| --- | --- |
| Notable read row when write granted | Done (2026-08-24 retest) |
| Near-miss wrap repeats `Would enable:` | Done (2026-08-24 retest) |
| Notable header on CRITICAL with `--explain-all` | Done - header no longer claims "no toxic combo matched" |
| `notable_grants` in default JSON | Open - only populated with `--explain`/`--explain-all`; document or change |
| Table default without near misses | Consider default `--no-near-misses` for table format only |
| Enrichment visibility | Log or count private App name lookup 404s (stderr summary) |
| Selected-repo repo list | Needs installation token or higher privilege API (product gap) |

### Phase 2 trace priorities

| Item | Source pattern | Notes |
| --- | --- | --- |
| Detect `pull_request_target` + `create-github-app-token` + fork checkout | `credential-access-prt-fork-iatt-exfiltration` | GHSA-9g93-rxr5-xhqw |
| Detect permissive `[bot]` actor checks in agent workflows | `execution-ai-agent-external-bot-trust` | Allow-list by App slug / installation ID |
| Document Zizmor as complementary workflow linter | marketplace-trust-limitations.md | Not a gh-app-check dependency |

### Release

| Item | Notes |
| --- | --- |
| Tag `gh-app-graph` v0.1.0 | After catalog calibration + org retest PASS |
| Drop `go.mod` replace | Pin tagged graph; needs real git checkout with remotes |
| Optional org squash to v0.1.0 | After calibration + retest |

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
