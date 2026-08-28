package output

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// MarkdownWriter renders results as a Markdown table, for PR comments and issues.
type MarkdownWriter struct{}

// Write renders a Markdown table of audit results.
func (MarkdownWriter) Write(w io.Writer, results []eval.AppAuditResult) error {
	scan := eval.OrgScanResult{Installations: results}
	return MarkdownWriter{}.WriteOrgScan(w, scan)
}

// WriteOrgScan renders scan metadata and a Markdown results table.
func (MarkdownWriter) WriteOrgScan(w io.Writer, scan eval.OrgScanResult) error {
	if scan.ScanPlatform != "" {
		line := fmt.Sprintf("**Scan target:** %s", scan.ScanPlatform)
		if scan.ExcludedGHESRules > 0 {
			line += fmt.Sprintf(" (%d GHES-only toxic rules excluded)", scan.ExcludedGHESRules)
		}
		fmt.Fprintln(w, line)
		fmt.Fprintln(w)
	}
	SortResults(scan.Installations)
	fmt.Fprintln(w, "| App | Risk | Repos | Writes | Violations | Toxic | Near misses | GHES scopes |")
	fmt.Fprintln(w, "| --- | --- | --- | --- | --- | --- | --- | --- |")
	for _, result := range scan.Installations {
		fmt.Fprintf(w, "| %s | %s | %s | %s | %s | %s | %s | %s |\n",
			escapeCell(appLabel(result)),
			result.RiskLevel,
			result.RepoSelection,
			strconv.Itoa(result.WriteScopeCount),
			escapeCell(strings.Join(result.Violations, "; ")),
			escapeCell(formatToxicSummary(result.ToxicMatches)),
			escapeCell(formatNearMissSummary(result.NearMisses)),
			escapeCell(strings.Join(result.GHESScopes, "; ")),
		)
	}
	return nil
}

func escapeCell(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
