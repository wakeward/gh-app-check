# Backlog / future investigations

Ideas that are out of scope for the current implementation phases (see
[Implementation Phases](../README.md#implementation-phases) in the README) but
worth investigating once the Phase 0-3 foundation is in place.

## Permission drift guard

**Idea:** [google/capslock](https://github.com/google/capslock) does static
capability analysis for Go packages - it classifies which privileged
operations a package can reach transitively, so you can see *what a
dependency is capable of* rather than just scanning it for known CVEs. Applied
to GitHub Apps, the equivalent question is: *what is an installed app capable
of right now, and has that capability set grown since we last checked?*

Phase 1 (the Control Plane auditor) is a point-in-time snapshot: it evaluates
an installation's current permissions against the least-privilege ruleset
(`pkg/rules`, `pkg/eval`) and reports a risk level. A permission drift guard
would sit on top of that and answer a different question: **did anything
change, and was the change an escalation?**

### Why this matters

- GitHub App permission grants can be silently escalated by the app's
  publisher between installs, and organization owners must re-approve, but
  that approval can be a routine click-through with no security review.
- A compromised or malicious app update is a realistic supply-chain vector:
  the install looked fine on day one, then a later version requests
  `administration: write` or drops `repository_selection` from `selected` to
  `all`.
- Tracking drift turns "is this app currently safe" (Phase 1) into "did this
  app just become less safe, and who approved that."

### Open questions to investigate before designing this

1. **Event source for changes.** Does the GitHub organization audit log
   expose a reliable `integration_installation.*` (or similar) event for
   permission changes, repository-selection changes, and suspensions? If not
   granular enough, is polling `GET /orgs/{org}/installations` on a schedule
   and diffing snapshots against a local store (see below) the only option?
   Check both the REST audit-log endpoint and any relevant webhook events
   (`installation`, `installation_repositories`) an org-owned app could
   subscribe to.
2. **Snapshot storage.** Where do historical snapshots live so a diff is
   possible? Options: a local SQLite/JSON store committed nowhere (stateful
   CLI run), a GitHub-hosted artifact (e.g. a private repo or GitHub Actions
   cache) for CI-driven runs, or delegate entirely to whatever the org already
   uses for time-series data. Needs a decision before Phase 4 can be scoped.
3. **What counts as "high risk" drift.** Reuse `pkg/eval`'s existing risk
   ranking (PASS/WARN/HIGH/CRITICAL) but the *interesting* signal is the
   delta, not the absolute level - e.g. WARN -> HIGH is worth alerting on even
   if HIGH isn't the worst possible outcome. Needs a small diff-specific
   ruleset (`pkg/rules` predicates operate on a single `Installation`; drift
   predicates would need to operate on an `(old, new)` pair).
4. **Alerting/output.** Does this become a new `gh app-check watch` (or
   `diff`) subcommand that's run on a schedule via GitHub Actions, posting to
   Slack/Issues/SARIF on drift? Or a library function other tools call? Decide
   after (1)-(3) are answered.
5. **False-positive rate.** Legitimate apps do add scopes for new features.
   The guard needs to distinguish "routine, expected growth" from "suspicious
   escalation" - possibly via an allow-list of expected permission bumps per
   app, or just surfacing every change for human review initially and
   tightening later.

### Relationship to existing phases

This is not Phase 1-3 as currently scoped in the README; it depends on Phase 1
existing first (the ruleset it diffs against). Tentatively "Phase 4" once
Phase 1-3 ship and the open questions above have answers.
