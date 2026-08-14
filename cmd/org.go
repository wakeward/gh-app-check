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
)

var orgTimeout time.Duration

var orgCmd = &cobra.Command{
	Use:   "org <org-name>",
	Short: "Audit the control plane (installed GitHub Apps) for an organization",
	Long: `Fetches all GitHub App installations for the given organization via
GET /orgs/{org}/installations and evaluates them against the least-privilege
rules engine (blast radius, toxic permissions).`,
	Args: cobra.ExactArgs(1),
	RunE: runOrg,
}

func init() {
	orgCmd.Flags().DurationVar(&orgTimeout, "timeout", 2*time.Minute, "Maximum time to wait for GitHub API responses")
}

func runOrg(_ *cobra.Command, args []string) error {
	org := args[0]

	token, err := resolveToken()
	if err != nil {
		return err
	}
	client, err := ghclient.New(token)
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
	toxic, err := graphdata.LoadToxicCombinations()
	if err != nil {
		return fmt.Errorf("load toxic combinations: %w", err)
	}

	results := make([]eval.AppAuditResult, 0, len(installations))
	for _, inst := range installations {
		results = append(results, eval.Evaluate(
			inst.Slug,
			inst.Name,
			org,
			rules.Installation{
				RepositorySelection: inst.RepositorySelection,
				Permissions:         inst.Permissions,
			},
			toxic,
		))
	}

	writer, err := output.ForFormat(format)
	if err != nil {
		return err
	}
	return writer.Write(os.Stdout, results)
}
