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
	SortResults(results)
	fmt.Fprintln(w, "| App | Risk | Repos | Writes | Violations | Toxic | Near misses |")
	fmt.Fprintln(w, "| --- | --- | --- | --- | --- | --- | --- |")
	for _, result := range results {
		fmt.Fprintf(w, "| %s | %s | %s | %s | %s | %s | %s |\n",
			escapeCell(appLabel(result)),
			result.RiskLevel,
			result.RepoSelection,
			strconv.Itoa(result.WriteScopeCount),
			escapeCell(strings.Join(result.Violations, "; ")),
			escapeCell(formatToxicSummary(result.ToxicMatches)),
			escapeCell(formatNearMissSummary(result.NearMisses)),
		)
	}
	return nil
}

func escapeCell(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
