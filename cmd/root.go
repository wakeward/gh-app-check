package cmd

import (
	"fmt"
	"os"

	ghauth "github.com/cli/go-gh/v2/pkg/auth"
	"github.com/spf13/cobra"
)

// format holds the persistent --format flag value shared by all subcommands.
var format string

var rootCmd = &cobra.Command{
	Use:   "gh-app-check",
	Short: "Audit GitHub App installations for least-privilege violations",
	Long: `gh-app-check is a security auditing tool for enterprise organizations to
evaluate installed GitHub Apps against the principle of least privilege,
identify toxic permission combinations, and trace the actual execution of
app credentials within the codebase.`,
	PersistentPreRunE: validateFormat,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "table", "Output format: table, json, or markdown")
	rootCmd.AddCommand(orgCmd)
	rootCmd.AddCommand(traceCmd)
}

// Execute runs the root command and exits non-zero on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// validateFormat rejects anything other than the three supported --format values
// before any subcommand logic runs.
func validateFormat(_ *cobra.Command, _ []string) error {
	switch format {
	case "table", "json", "markdown":
		return nil
	default:
		return fmt.Errorf("invalid --format %q: must be one of table, json, markdown", format)
	}
}

// resolveToken returns a GitHub token for api.github.com, preferring an
// explicit GH_TOKEN/GITHUB_TOKEN env var, then falling back to the gh CLI's
// stored credentials via github.com/cli/go-gh's auth package.
func resolveToken() (string, error) {
	token, _ := ghauth.TokenForHost("github.com")
	if token == "" {
		return "", fmt.Errorf("no GitHub token found: set GH_TOKEN/GITHUB_TOKEN or run `gh auth login`")
	}
	return token, nil
}
