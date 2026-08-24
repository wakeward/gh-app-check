package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/wakeward/gh-app-check/pkg/eval"
	"github.com/wakeward/gh-app-check/pkg/ghclient"
	"github.com/wakeward/gh-app-check/pkg/output"
	"github.com/wakeward/gh-app-check/pkg/rules"
	graphdata "github.com/wakeward/gh-app-graph/pkg/data"
	graphplatform "github.com/wakeward/gh-app-graph/pkg/platform"
)

var (
	orgTimeout        time.Duration
	orgPlatformFlag   string
	orgNoEnrichNames  bool
	orgNoNearMisses   bool
)

var orgCmd = &cobra.Command{
	Use:   "org <org-name>",
	Short: "Audit the control plane (installed GitHub Apps) for an organization",
	Long: `Fetches all GitHub App installations for the given organization via
GET /orgs/{org}/installations and evaluates them against the least-privilege
rules engine (blast radius, toxic permissions).

GHES-only permissions and toxic rules are excluded automatically on GitHub.com
and Enterprise Cloud scans. Use --platform ghes when auditing self-hosted
Enterprise Server, or --platform auto (default) to follow gh auth host detection.`,
	Args: cobra.ExactArgs(1),
	RunE: runOrg,
}

func init() {
	orgCmd.Flags().DurationVar(&orgTimeout, "timeout", 2*time.Minute, "Maximum time to wait for GitHub API responses")
	orgCmd.Flags().StringVar(&orgPlatformFlag, "platform", "auto", "Scan target: auto (from gh auth), cloud (exclude GHES-only rules), or ghes")
	orgCmd.Flags().BoolVar(&orgNoEnrichNames, "no-enrich-names", false, "Skip GET /apps/{slug} lookups for friendly display names")
	orgCmd.Flags().BoolVar(&orgNoNearMisses, "no-near-misses", false, "Omit near-miss toxic combinations from output")
}

func runOrg(_ *cobra.Command, args []string) error {
	org := args[0]

	auth, err := ghclient.ResolveAuth()
	if err != nil {
		return err
	}
	scanPlatform, err := ghclient.ResolveScanPlatform(orgPlatformFlag, auth.Platform)
	if err != nil {
		return err
	}
	if scanPlatform == ghclient.ScanPlatformCloud && auth.Platform == ghclient.ScanPlatformGHES {
		fmt.Fprintf(os.Stderr, "gh-app-check: warning: gh auth host %q is GHES but --platform cloud excludes GHES-only rules\n", auth.Host)
	}
	if scanPlatform == ghclient.ScanPlatformGHES && auth.Platform == ghclient.ScanPlatformCloud {
		fmt.Fprintf(os.Stderr, "gh-app-check: warning: gh auth host %q is GitHub.com but --platform ghes includes GHES-only rules (none may apply)\n", auth.Host)
	}

	client, err := ghclient.NewForHost(auth.Token, auth.Host)
	if err != nil {
		return fmt.Errorf("create github client: %w", err)
	}

	ctx := context.Background()
	if orgTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, orgTimeout)
		defer cancel()
	}

	installations, err := ghclient.ListOrgInstallations(ctx, client, org)
	if err != nil {
		return err
	}
	if !orgNoEnrichNames {
		if err := ghclient.EnrichInstallationNames(ctx, client, installations); err != nil {
			return fmt.Errorf("enrich app names: %w", err)
		}
	}
	toxic, err := graphdata.LoadToxicCombinations()
	if err != nil {
		return fmt.Errorf("load toxic combinations: %w", err)
	}
	ghesKeys, err := graphdata.LoadGHESOnlyAPIKeys()
	if err != nil {
		return fmt.Errorf("load GHES-only permission keys: %w", err)
	}

	includeGHES := scanPlatform.IncludesGHESRules()
	excludedRules := 0
	if !includeGHES {
		filtered := graphplatform.FilterToxicCombinations(toxic, ghesKeys, false)
		excludedRules = len(toxic) - len(filtered)
		toxic = filtered
	}

	scanCtx := eval.ScanContext{
		IncludeGHESRules: includeGHES,
		GHESOnlyKeys:     ghesKeys,
	}

	results := make([]eval.AppAuditResult, 0, len(installations))
	for _, inst := range installations {
		result := eval.EvaluateWithContext(
			inst.Slug,
			inst.Name,
			org,
			rules.Installation{
				RepositorySelection: inst.RepositorySelection,
				Permissions:         inst.Permissions,
			},
			toxic,
			scanCtx,
		)
		result.InstallationID = inst.InstallationID
		result.AppID = inst.AppID
		result.HTMLURL = inst.HTMLURL
		result.Permissions = inst.Permissions
		if orgNoNearMisses {
			result.NearMisses = nil
		}
		results = append(results, result)
	}

	report := eval.OrgScanResult{
		ScanPlatform:      scanPlatform.Label(auth.Host),
		ScanHost:          auth.Host,
		ExcludedGHESRules: excludedRules,
		Installations:     results,
	}
	return output.WriteOrgScan(os.Stdout, report, format)
}
