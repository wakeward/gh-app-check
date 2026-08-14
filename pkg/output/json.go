package output

import (
	"encoding/json"
	"io"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// JSONWriter renders results as structured JSON, for piping to jq or SIEM ingestion.
type JSONWriter struct{}

// Write encodes results as indented JSON.
func (JSONWriter) Write(w io.Writer, results []eval.AppAuditResult) error {
	SortResults(results)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}
