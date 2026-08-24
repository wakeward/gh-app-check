# Company org validation runbook (Phase 1)

**Purpose:** Safe, read-only validation of `gh-app-check` against a real organization.
Copy this file (or the Cursor prompt section below) when running from a company
clone next week.

**Hard constraints for this session:**

- **Read-only by default** - no installs, uninstalls, permission upgrades, workflow
  edits, branch protection changes, or secret rotation unless an optional owned-App
  check is explicitly approved in chat first.
- **No audit log work** - you do not have org audit log API access; skip all
  `observable_artifacts` / detection validation that depends on it.
- **No Phase 2 `trace`** - command is not implemented; do not attempt Code Search
  sweeps unless implementing Phase 2 in the same session.
- **No org names in git** - capture notes in `docs/FEEDBACK-*.md` (gitignored) or
  a local file outside the repo; redact slugs if sharing externally.

---

## Cursor agent prompt (copy from here)

```
You are helping validate gh-app-check Phase 1 against my company GitHub organization.

## Context

- Tool: gh-app-check (gh CLI extension), sibling gh-app-graph for embedded rules.
- Phase 1 scope: GET /orgs/{org}/installations + permission/toxic evaluation only.
- I have org admin (or sufficient) access with gh auth; I own several GitHub Apps
  installed in this org that I may use for spot checks.
- I do NOT have audit log access. Do not plan tests that require integration_installation.* events.
- I do NOT want destructive or state-changing tests unless I explicitly approve a
  single bounded action in chat.

## Your job

1. Verify build and gh extension wiring locally.
2. Run the read-only test matrix below in order; stop on auth/API errors and diagnose.
3. For each thread, record: pass/fail, unexpected output, false positive/negative hunch.
4. Compare 2-3 installations I know well (especially my own Apps) against GitHub UI permissions.
5. Write sanitized findings to docs/FEEDBACK-<date>.md (gitignored) - no company name required in filename if I prefer FEEDBACK-local.md.

## Hard stops (never run without my explicit "yes, run it")

- gh app-check trace (not implemented; would hit Code Search)
- Installing/uninstalling Apps on production org repos
- Changing App permissions or repository_selection on any App I do not solely own
- Creating repos, workflows, or PATs for attack simulation
- Any gcloud/terraform/kubectl apply or billable cloud commands

## Read-only commands to run

Replace ORG with my org slug when I provide it.

### A. Preflight

gh auth status
gh auth refresh -s read:org
go build -o gh-app-check .
go test ./...
gh extension install . --force

### B. Core org scan (default)

gh app-check org ORG --platform auto

### C. Output formats

gh app-check org ORG --format json | jq '{scan_platform, excluded_ghes_rules, count: (.installations|length), risks: [.installations[] | {slug: .app_slug, risk: .risk_level, toxic: (.toxic_matches|length), near: (.near_misses|length)}]}'
gh app-check org ORG --format markdown > /tmp/gh-app-check-org.md
# Review markdown locally; do not commit

### D. Platform / GHES filtering

gh app-check org ORG --platform cloud
gh app-check org ORG --platform auto
# If on GHES host only:
# gh app-check org ORG --platform ghes

Compare: excluded_ghes_rules count, any ghes_scopes highlights on cloud vs ghes.

### E. Friendly names (P2)

gh app-check org ORG
gh app-check org ORG --no-enrich-names
# Confirm table shows "Name (slug)" vs slug-only when enrichment on

### F. Timeout / error handling

gh app-check org ORG --timeout 30s
# Expect success if org is small; note if timeout too aggressive

### G. Spot-check known Apps (manual cross-check)

For each App I name (especially ones I own):
- Open GitHub UI: Org Settings -> GitHub Apps -> Installed apps -> permissions + repo access
- Compare repository_selection (all vs selected), permission keys, read vs write
- Note any mapping mismatch between API output and UI

### H. Signal quality review (subjective)

For the full org scan, note:
- Apps flagged CRITICAL/HIGH that feel correct vs surprising
- Toxic combo names that fire and whether exploit_path matches intuition
- Near misses that are useful vs noisy
- All-repos + read-only WARN vs all-repos + write HIGH
- administration: write CRITICAL entries
- God-mode write count (HIGH when >5 writes)

### I. Optional bounded owned-App check (only if I approve)

If I confirm App SLUG is mine and sandbox-only:
- UI-only review of manifest vs scan row (still read-only), OR
- I may adjust a test App installation on a single sandbox repo in UI - you wait for
  me to do it manually, then re-run: gh app-check org ORG --format json | jq filter for SLUG

Do NOT click install/upgrade flows for me.

## Out of scope this session

- Audit log event verification
- Attack pattern desk traces (PRT exfil, ephemeral admin, etc.)
- Permission catalog YAML edits (already human-confirmed)
- Blog drafts
- Publishing or pushing results to git

## Deliverable

A short FEEDBACK doc with sections:
1. Environment (cloud vs GHES, gh version, extension version/commit)
2. Test matrix results (A-I)
3. False positives / false negatives (with App slug redacted if needed)
4. UX notes (table columns, JSON fields missing, performance)
5. Recommended backlog items for next sprint
```

---

## Test matrix (quick reference)

| ID | Thread | Command / action | API impact | Destructive? |
|---|---|---|---|---|
| A | Build + unit tests | `go test ./...` | None | No |
| B | Default org audit | `gh app-check org ORG` | Read installations + app metadata | No |
| C | JSON / markdown output | `--format json`, `--format markdown` | Same as B | No |
| D | GHES rule filtering | `--platform auto` vs `cloud` vs `ghes` | Same as B | No |
| E | App name enrichment | default vs `--no-enrich-names` | Extra GET /apps/{slug} per install | No |
| F | Timeout | `--timeout 30s` | Same as B | No |
| G | Known-App cross-check | Manual UI vs scan row | None (manual) | No |
| H | Signal quality | Review toxic/near-miss/all-repos rules | None | No |
| I | Owned-App delta | Manual UI change, re-scan one slug | Read only after you change UI | **You** change UI |

---

## Auth troubleshooting (read-only)

```bash
gh auth status
gh auth refresh -s read:org
```

403/404 on installations:

- Confirm org **admin** (not repo admin only)
- Re-auth with `read:org` at minimum

---

## What we learn without audit logs

| Question | Answerable this week? | How |
|---|---|---|
| Permission mapping correct? | Yes | UI cross-check (G, I) |
| Toxic combos fire appropriately? | Yes | H |
| GHES scopes filtered on cloud? | Yes | D |
| JSON good for automation? | Yes | C |
| Who approved install? | No | Skip |
| Permission drift over time? | No | Phase 4; needs snapshots or audit log |
| Workflow token misuse? | No | Phase 2 not built |

---

## After the session

1. Keep `docs/FEEDBACK-*.md` local (gitignored).
2. Open gh-app-check issues for concrete bugs (mapping, false positives).
3. If signal quality is good: tag `gh-app-graph`, pin version, drop `go.mod` replace.
4. Schedule theorizing session separately (attack pattern gap-check, not blocked on org run).

---

## Related docs

- [`INSTALLATION.md`](INSTALLATION.md) - auth and deployment
- [`improvement-backlog.md`](improvement-backlog.md) - known product backlog
- gh-app-graph [`docs/methodology.md`](https://github.com/wakeward/gh-app-graph/blob/main/docs/methodology.md) - what Phase 1 measures
