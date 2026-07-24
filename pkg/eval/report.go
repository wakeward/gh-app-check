// Package eval aggregates pkg/rules predicates into per-installation audit
// results (Design Spec §4.A key struct).
package eval

// AppAuditResult is the outcome of evaluating one GitHub App installation
// against the least-privilege rules engine.
type AppAuditResult struct {
	AppSlug       string
	Owner         string
	RepoSelection string // "all" or "selected"
	RiskLevel     string // "CRITICAL", "HIGH", "WARN", "PASS"
	Violations    []string
}
