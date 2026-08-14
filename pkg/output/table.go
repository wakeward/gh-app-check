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
	SortResults(results)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "APP\tRISK\tREPOS\tWRITES\tVIOLATIONS\tTOXIC\tNEAR_MISSES")
	for _, result := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%s\t%s\t%s\n",
			appLabel(result),
			result.RiskLevel,
			result.RepoSelection,
			result.WriteScopeCount,
			strings.Join(result.Violations, "; "),
			formatToxicSummary(result.ToxicMatches),
			formatNearMissSummary(result.NearMisses),
		)
	}
	return tw.Flush()
}
