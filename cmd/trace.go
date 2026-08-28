package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var traceOrg string

var traceCmd = &cobra.Command{
	Use:   "trace <app-slug>",
	Short: "Audit the execution plane for a specific app's credential usage (Phase 2 - not implemented)",
	Long: `Searches the organization's codebase for usage of the given GitHub App
slug and inspects the discovered Actions workflows for insecure token
generation patterns (missing repositories/permissions scoping, leaked keys).`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		appSlug := args[0]

		if _, err := resolveToken(); err != nil {
			return err
		}

		// TODO(Phase 2): implement Code Search API integration and YAML AST
		// evaluation in pkg/ghclient / a future pkg/trace package.
		return fmt.Errorf("trace for app %q in org %q is not implemented yet (Phase 2); see docs/INSTALLATION.md", appSlug, traceOrg)
	},
}

func init() {
	traceCmd.Flags().StringVar(&traceOrg, "org", "", "GitHub organization to search (required)")
	if err := traceCmd.MarkFlagRequired("org"); err != nil {
		panic(err)
	}
}
