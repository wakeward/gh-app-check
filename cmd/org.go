package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var orgCmd = &cobra.Command{
	Use:   "org <org-name>",
	Short: "Audit the control plane (installed GitHub Apps) for an organization",
	Long: `Fetches all GitHub App installations for the given organization via
GET /orgs/{org}/installations and evaluates them against the least-privilege
rules engine (blast radius, toxic permissions).`,
	Args: cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		org := args[0]

		if _, err := resolveToken(); err != nil {
			return err
		}

		// TODO(Phase 1): fetch installations via pkg/ghclient, evaluate via
		// pkg/eval, and render via pkg/output using the --format flag.
		return fmt.Errorf("org audit for %q is not implemented yet (Phase 1); see docs/INSTALLATION.md", org)
	},
}
