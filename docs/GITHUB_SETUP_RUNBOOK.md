# GitHub repository setup runbook

Configure GitHub settings before the first public push. Apply to **both**
`wakeward/gh-app-check` and `wakeward/gh-app-graph` unless noted.

## Prerequisite: authenticate `gh`

```bash
gh auth login -h github.com
gh auth status
```

You need admin access on each repository to configure branch protection.

### gh-app-check CI (graph module fetch)

Both repos are **public**. CI downloads `github.com/wakeward/gh-app-graph` from
the public Go module proxy (`go mod download`). No extra GitHub secret is required.

If you previously set `GH_GRAPH_READ_TOKEN` while `gh-app-graph` was private, delete
it from **Settings → Secrets and variables → Actions** (and Dependabot secrets).

---

## Branch protection (solo maintainer)

Apply **after at least one successful workflow run** on `main` (or a PR) so
GitHub knows the required status check names.

### Settings to enable (GitHub UI)

Repository → **Settings** → **Branches** → **Add branch protection rule**
(or **Rules** → ruleset) for `main`:

| Setting | Value | Why |
|---|---|---|
| Require a pull request before merging | **On** | No direct pushes to `main` |
| Required approvals | **0** | Solo maintainer: PR gate + CI, merge without a second human. Setting `1` blocks you unless you use a bot/second account. |
| Dismiss stale pull request approvals | **On** | Scorecard branch-protection tier |
| Require review from Code Owners | **On** | Requires `.github/CODEOWNERS` (both repos) |
| Require status checks to pass | **On** | Blocking CI |
| Require branches to be up to date | **On** (`strict: true`) | Scorecard tier 3 |
| Status checks (exact job names) | See below | Must match workflow `jobs.*.name` |
| Require conversation resolution | **On** | Low-cost hygiene |
| Require linear history | **On** | Prevents merge commits on `main` |
| Do not allow bypassing | **On** / enforce admins | Includes admins |
| Allow force pushes | **Off** | Scorecard tier 1 |
| Allow deletions | **Off** | Scorecard tier 1 |

### Required status check names

These are the **job** `name:` fields (not the workflow file title):

| Check name | Workflow file |
|---|---|
| `golangci-lint` | `.github/workflows/lint_tests.yml` |
| `Go Build, Vet & Unit Tests` | `.github/workflows/unit_tests.yml` |
| `govulncheck & gosec` | `.github/workflows/security_scan.yml` |

Optional: add `Scorecard analysis` once Scorecard has run once (informational;
not required for merge unless you want it blocking).

**gh-app-graph only:** `refresh.yml` runs on schedule; do not require it for
every PR (it only runs weekly / manual dispatch).

### Solo maintainer workflow

1. Create a branch: `git checkout -b feat/my-change`
2. Push and open a PR to `main`
3. Wait for the three required checks to pass
4. Merge the PR (no second reviewer needed with approvals = 0)

Scorecard's separate **Code-Review** check will still score low on a solo
project. That is expected; do not fake reviews with a second account.

### Expected Scorecard branch-protection score

With the settings above: roughly **5-6/10** (tier 3: status checks + no force
push). **9/10** requires two human reviewers, which is not practical solo.

### API equivalent (gh-app-check example)

Run once checks exist on GitHub:

```bash
REPO=wakeward/gh-app-check   # or wakeward/gh-app-graph

cat <<'EOF' | gh api --method PUT "repos/${REPO}/branches/main/protection" --input -
{
  "required_status_checks": {
    "strict": true,
    "contexts": [
      "golangci-lint",
      "Go Build, Vet & Unit Tests",
      "govulncheck & gosec"
    ]
  },
  "enforce_admins": true,
  "required_pull_request_reviews": {
    "required_approving_review_count": 0,
    "dismiss_stale_reviews": true,
    "require_code_owner_reviews": true
  },
  "restrictions": null,
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
EOF
```

If the API rejects unknown context names, push once, let workflows run, then
re-run with the exact names shown under the PR **Checks** tab.

---

## Repo-level security settings

```bash
REPO=wakeward/gh-app-check   # repeat for gh-app-graph

gh api --method PATCH "repos/${REPO}" \
  -F security_and_analysis[secret_scanning][status]=enabled \
  -F security_and_analysis[secret_scanning_push_protection][status]=enabled \
  -F security_and_analysis[dependabot_security_updates][status]=enabled
```

CodeQL default setup (Go):

```bash
gh api --method PATCH "repos/${REPO}/code-scanning/default-setup" \
  -f state=configured \
  -f query_suite=default \
  -f 'languages[]=go'
```

---

## Order of operations

1. Ensure workflows exist on `main` (at least one green run)
2. Apply branch protection with real check names
3. Enable secret scanning + push protection + Dependabot security updates
4. Enable CodeQL default setup
5. Continue Phase B (releases) before flipping public

---

## Verify OpenSSF Scorecard

After `scorecard.yml` runs on `main`:

```bash
gh api "repos/wakeward/gh-app-check/actions/workflows" \
  | jq '.workflows[] | select(.name=="Scorecard supply-chain security")'
```

Badge: `https://scorecard.dev/viewer/?uri=github.com/wakeward/gh-app-check`

---

## Release engineering (Phase B - not part of branch protection)

Signed **git tags** (GPG), signed release artifacts (Cosign), and SBOM (syft) are
Phase B. See [`PUBLISH-READINESS.md`](PUBLISH-READINESS.md).

### Signed git tags (required)

All release tags (`v*`) must be **GPG-signed** before push. Unsigned tags fail
OpenSSF Scorecard **Signed-Releases** expectations and weaken provenance.

**One-time setup:**

```bash
# Confirm signing key (must match a key uploaded to https://github.com/settings/keys)
gpg --list-secret-keys --keyid-format=long
git config user.signingkey <KEYID>
git config tag.gpgSign true   # optional: sign all tags by default
```

**Cut a release tag:**

```bash
git tag -s v0.1.0 -m "Initial public release v0.1.0."
git tag -v v0.1.0             # must show "Good signature"
git push origin v0.1.0
```

To move a tag after orphan squash (Phase C):

```bash
git tag -d v0.1.0
git tag -s v0.1.0 -m "Initial public release v0.1.0." <new-commit>
git push --force origin v0.1.0
```

**gh-app-graph:** tag catalog releases (e.g. `v0.1.0`).  
**gh-app-check:** signed tag push triggers `.github/workflows/release.yml`
(GoReleaser + Cosign on release artifacts).
