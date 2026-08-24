// Package output renders eval.AppAuditResult slices in the three formats the
// design spec requires: table, json, and markdown.
package output

import (
	"fmt"
	"io"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// Writer renders a set of AppAuditResults to an io.Writer.
type Writer interface {
	Write(io.Writer, []eval.AppAuditResult) error
}

// WriteOrgScan renders a full organization scan including platform metadata.
func WriteOrgScan(w io.Writer, scan eval.OrgScanResult, format string, opts Options) error {
	if opts.Explain && (format == "table" || format == "markdown") {
		return ExplainWriter{ExplainAll: opts.ExplainAll}.WriteOrgScan(w, scan)
	}
	switch format {
	case "table":
		return TableWriter{}.WriteOrgScan(w, scan)
	case "json":
		return JSONWriter{}.WriteOrgScan(w, scan)
	case "markdown":
		return MarkdownWriter{}.WriteOrgScan(w, scan)
	default:
		return fmt.Errorf("unknown format %q: must be table, json, or markdown", format)
	}
}

// ForFormat returns the Writer implementation for the given --format value.
func ForFormat(format string) (Writer, error) {
	switch format {
	case "table":
		return TableWriter{}, nil
	case "json":
		return JSONWriter{}, nil
	case "markdown":
		return MarkdownWriter{}, nil
	default:
		return nil, fmt.Errorf("unknown format %q: must be table, json, or markdown", format)
	}
}
