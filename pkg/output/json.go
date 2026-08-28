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
	scan := eval.OrgScanResult{Installations: results}
	return JSONWriter{}.WriteOrgScan(w, scan)
}

// WriteOrgScan encodes the full scan envelope as indented JSON.
func (JSONWriter) WriteOrgScan(w io.Writer, scan eval.OrgScanResult) error {
	SortResults(scan.Installations)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(scan)
}
