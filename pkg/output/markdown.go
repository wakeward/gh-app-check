package output

import (
	"fmt"
	"io"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// MarkdownWriter renders results as a Markdown table, for PR comments and issues.
type MarkdownWriter struct{}

// Write is a placeholder until Phase 1 wires in real installation data.
func (MarkdownWriter) Write(w io.Writer, results []eval.AppAuditResult) error {
	// TODO(Phase 1): render a Markdown table.
	_, err := fmt.Fprintf(w, "%d result(s) (markdown rendering not implemented yet)\n", len(results))
	return err
}
