package output

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// TableWriter renders results as a human-readable terminal table.
type TableWriter struct{}

// Write renders one row per installation with risk and violation summary.
func (TableWriter) Write(w io.Writer, results []eval.AppAuditResult) error {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "APP\tRISK\tREPOS\tVIOLATIONS")
	for _, result := range results {
		app := result.AppSlug
		if app == "" {
			app = result.AppName
		}
		if app == "" {
			app = "(unknown app)"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			app,
			result.RiskLevel,
			result.RepoSelection,
			strings.Join(result.Violations, "; "),
		)
	}
	return tw.Flush()
}
