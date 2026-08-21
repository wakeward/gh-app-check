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
	return TableWriter{}.writeResults(w, results, "")
}

// WriteOrgScan renders scan metadata and installation rows.
func (tw TableWriter) WriteOrgScan(w io.Writer, scan eval.OrgScanResult) error {
	header := fmt.Sprintf("# scan: %s", scan.ScanPlatform)
	if scan.ExcludedGHESRules > 0 {
		header += fmt.Sprintf("; excluded %d GHES-only toxic rule(s)", scan.ExcludedGHESRules)
	}
	return tw.writeResults(w, scan.Installations, header)
}

func (TableWriter) writeResults(w io.Writer, results []eval.AppAuditResult, header string) error {
	if header != "" {
		fmt.Fprintln(w, header)
	}
	SortResults(results)
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "APP\tRISK\tREPOS\tWRITES\tVIOLATIONS\tTOXIC\tNEAR_MISSES\tGHES")
	for _, result := range results {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%s\t%s\t%s\t%s\n",
			appLabel(result),
			result.RiskLevel,
			result.RepoSelection,
			result.WriteScopeCount,
			strings.Join(result.Violations, "; "),
			formatToxicSummary(result.ToxicMatches),
			formatNearMissSummary(result.NearMisses),
			strings.Join(result.GHESScopes, "; "),
		)
	}
	return tw.Flush()
}
