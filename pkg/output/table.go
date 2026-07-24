package output

import (
	"fmt"
	"io"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// TableWriter renders results as a human-readable terminal table.
type TableWriter struct{}

// Write is a placeholder until Phase 1 wires in real installation data;
// it currently emits a summary line so the CLI remains usable end-to-end.
func (TableWriter) Write(w io.Writer, results []eval.AppAuditResult) error {
	// TODO(Phase 1): render an aligned terminal table (e.g. via text/tabwriter).
	_, err := fmt.Fprintf(w, "%d result(s) (table rendering not implemented yet)\n", len(results))
	return err
}
