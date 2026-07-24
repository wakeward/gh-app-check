# GitHub setup runbook (not yet executed)

This is a runbook, not a script that has been run. Nothing here has been executed against GitHub. Review each step, then run them yourself (or ask me to run a specific step once you've confirmed it).

## Prerequisite: re-authenticate `gh`

The local `gh` CLI token is currently invalid (checked during this session: `gh auth status` reported "The token in default is invalid"). Re-authenticate before anything else:

```bash
gh auth login -h github.com
```

## Step 1: Create the GitHub repository (no push yet)

```bash
cd /home/wakeward/src/gh-app-check

gh repo create wakeward/gh-app-check \
  --public \
  --source=. \
  --remote=origin \
  --description "GitHub CLI extension to audit GitHub App installations for least-privilege violations"
```

This creates the empty remote repository and wires up `origin` locally. It does **not** push any commits (no `--push` flag). Verify afterward with:

```bash
git remote -v
```

## Step 2: Configure branch protection on `main`

Solo-maintainer-tuned: PR required with **0** required approvals (you merge your own PRs once checks pass - GitHub blocks self-approval regardless of this number, so `1` would just lock you out), everything else Scorecard's `Branch-Protection` check rewards enabled, admins included in enforcement (no bypass list).

```bash
cat <<'EOF' | gh api --method PUT repos/wakeward/gh-app-check/branches/main/protection --input -
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
    "dismiss_stale_reviews": true
  },
  "restrictions": null,
  "required_linear_history": true,
  "allow_force_pushes": false,
  "allow_deletions": false,
  "required_conversation_resolution": true
}
EOF
```

Notes:
- The three `contexts` entries must match the job `name:` fields exactly as GitHub reports them after the first workflow run on `main` or on a PR - the API call may need the contexts added/corrected once you've seen the actual check names appear (they won't exist as valid contexts until the workflows have run at least once). It's fine to run this command with an empty `contexts: []` first, push, let the workflows run once, then re-run this command with the real context names.
- `required_conversation_resolution: true` is a low-cost addition beyond the original plan (mirrors gh-branch-auditor's `GH-BP-007` control) - remove if you don't want it.
- The separate `Code-Review` Scorecard check will score low/zero as a genuinely solo project. Not addressed here - not something to fake with a second account.

## Step 3: Repo-level security settings

```bash
gh api --method PATCH repos/wakeward/gh-app-check \
  -F security_and_analysis[secret_scanning][status]=enabled \
  -F security_and_analysis[secret_scanning_push_protection][status]=enabled \
  -F security_and_analysis[dependabot_security_updates][status]=enabled
```

CodeQL default setup (separate endpoint; Go is a supported CodeQL language):

```bash
gh api --method PUT repos/wakeward/gh-app-check/code-scanning/default-setup \
  -f state=configured \
  -f query_suite=default \
  -f 'languages[]=go'
```

Secret scanning itself is typically already on by default for new public repos - the command above is a belt-and-braces confirmation, and explicitly turns on push protection and Dependabot security updates, which are not on by default.

## Step 4: Push, only after Steps 1-3 are confirmed in place

```bash
git push -u origin main
```

Then go back and re-run the Step 2 command with the real check-run context names once the workflows have executed at least once (see note above).

## Step 5 (separate follow-up pass, not now)

Release engineering - `.goreleaser.yaml`, `.github/workflows/release.yml`, Cosign signing, syft SBOM, SLSA L3 provenance - adapted from your existing [gh-branch-auditor](https://github.com/wakeward/gh-branch-auditor) templates. Not required for a strong pre-release Scorecard score; do this deliberately once you're ready to cut `v0.1.0`, not squeezed into this bootstrap pass.

## Verifying the OpenSSF Scorecard result

Once pushed and the `scorecard.yml` workflow has run at least once on `main`:

```bash
gh api repos/wakeward/gh-app-check/actions/workflows | jq '.workflows[] | select(.name=="Scorecard supply-chain security")'
```

Or check the badge directly (may take a few minutes to populate after first run): `https://securityscorecards.dev/viewer/?uri=github.com/wakeward/gh-app-check`
