# Installation & Deployment Guide

gh-app-check can be run locally by human administrators or deployed as a continuous monitoring job in CI/CD (the recommended enterprise pattern).

## 1. Local CLI Installation

Install the extension directly via the GitHub CLI:

```bash
gh extension install wakeward/gh-app-check
```

### Local development (sibling checkout)

`go build ./...` compiles packages but does **not** emit the extension binary.
Build the command explicitly, then install from the repo root:

```bash
go build -o gh-app-check .
gh extension install . --force
```

Without `gh-app-check` on disk, `gh extension install .` fails with
`fork/exec .../gh-app-check: no such file or directory`.

### Prerequisites

Listing organization installations requires a token that can read organization
data. In practice, **`read:org`** on a user token belonging to an **organization
admin** has been enough to call `GET /orgs/{org}/installations`. GitHub may
tighten this; treat 403/404 as an auth or role problem first.

- **Organization role:** org **admin** (not merely a repository admin) is the
  practical minimum for org-wide installation listing.
- **OAuth scope:** prefer `read:org` at minimum. Some docs mention `admin:org`;
  re-authenticate if you see scope-related 403 errors:

```bash
gh auth login --scopes "read:org,repo"
```

Organization **owners** always have sufficient access; admins often do as well.

## 2. Enterprise CI/CD Deployment (Recommended)

For enterprise environments, relying on a human's highly privileged Personal Access Token is a security anti-pattern. We strongly recommend running gh-app-check continuously in a scheduled GitHub Actions workflow.

### Step 1: Create an Internal GitHub App

Create a custom GitHub App owned by your organization to act as the service account for this tool.

**Required Permissions:**

- **Organization Administration:** Read-only (required to list installed apps)
- **Metadata:** Read-only (mandatory baseline)
- **Contents:** Read-only (required only if you intend to use the `--trace` command to scan codebase workflows)

### Step 2: Install the App

Install this new App on the repository where your monitoring Actions workflow will live.

### Step 3: Configure the Workflow

Use an action like `actions/create-github-app-token` to generate a short-lived token for your app, and pass it to the CLI. Pin any third-party action you add here by full commit SHA, respecting a cooldown window after release, per this repo's own supply-chain policy.

```yaml
name: Nightly App Audit
on:
  schedule:
    - cron: '0 2 * * *' # Run at 2 AM daily

permissions:
  contents: read

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - name: Generate Token
        id: generate-token
        uses: actions/create-github-app-token@v2
        with:
          app-id: ${{ secrets.APP_ID }}
          private-key: ${{ secrets.APP_PRIVATE_KEY }}

      - name: Install gh-app-check
        run: gh extension install wakeward/gh-app-check
        env:
          GH_TOKEN: ${{ steps.generate-token.outputs.token }}

      - name: Run Organization Audit
        run: gh app-check org my-organization --format markdown >> "$GITHUB_STEP_SUMMARY"
        env:
          GH_TOKEN: ${{ steps.generate-token.outputs.token }}
```

## 3. Troubleshooting & Common Errors

If you are running into issues executing gh-app-check, review the common errors below.

### 404 Not Found or 403 Forbidden on the Org Command

**Symptom:**

```
Error: failed to fetch installations: GET https://api.github.com/orgs/my-org/installations: 404 Not Found []
```

**Cause:**

The GitHub API endpoint for listing organization installations requires
**organization admin or owner** membership. This is an **org role** check, not
merely an OAuth scope check.

Common cases:

1. **Org member, not admin.** `read:org` is present but the account is only a
   member on that org. The same token may work on a different org where you are
   admin.
2. **Wrong org slug** or no access to the org at all.

**False lead:** `gh` may suggest `gh auth refresh -s admin:org` on 404. Adding
`admin:org` does **not** grant org admin role. Verify membership role first:

```bash
gh api orgs/my-org/memberships/$(gh api user -q .login) -q .role
```

You need `admin` or `owner`. Re-authenticate with `read:org` if scopes are missing:

```bash
gh auth login --scopes "read:org,repo"
```

### Secondary Rate Limit Exceeded during `--trace` (Phase 2 - not implemented)

> **Note:** `gh app-check trace` is not implemented yet. The guidance below applies
> to the planned Phase 2 command only.

**Symptom:**

```
Error: you have exceeded a secondary rate limit. Please wait a few minutes before you try again.
```

**Cause:**

The `--trace` command relies heavily on the GitHub Code Search API (`GET /search/code`) to hunt for workflow usage and `.pem` keys. The Code Search API has aggressive secondary rate limits, often capping at 30 requests per minute depending on the query complexity.

**Resolution:**

The CLI implements an exponential backoff strategy, but in massive organizations, it may still fail if the limits are exhausted globally by other users.

- Wait 5-10 minutes and try again.
- If running in CI/CD, ensure you are authenticating with a GitHub App token, which has its own dedicated rate limit bucket separate from human users.

### "No valid Actions workflows found" during Trace (Phase 2 - not implemented)

> **Note:** `gh app-check trace` is not implemented yet.

**Cause:**

The tracer looks for specific `actions/create-github-app-token` syntax. If your organization relies on a custom, internally built Docker action or a heavy wrapper script to generate tokens instead of the standard community actions, the AST parser will not detect the generation steps.

**Resolution:**

Currently, gh-app-check is tuned for the standard ecosystem. If you use custom wrappers, you may need to rely strictly on the Control Plane audit (the default `org` command).
