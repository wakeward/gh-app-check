package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// MarkdownWriter renders results as a Markdown table, for PR comments and issues.
type MarkdownWriter struct{}

// Write renders a Markdown table of audit results.
func (MarkdownWriter) Write(w io.Writer, results []eval.AppAuditResult) error {
	fmt.Fprintln(w, "| App | Risk | Repos | Violations |")
	fmt.Fprintln(w, "| --- | --- | --- | --- |")
	for _, result := range results {
		app := result.AppSlug
		if app == "" {
			app = result.AppName
		}
		if app == "" {
			app = "(unknown app)"
		}
		fmt.Fprintf(w, "| %s | %s | %s | %s |\n",
			escapeCell(app),
			result.RiskLevel,
			result.RepoSelection,
			escapeCell(strings.Join(result.Violations, "; ")),
		)
	}
	return nil
}

func escapeCell(s string) string {
	return strings.ReplaceAll(s, "|", "\\|")
}
