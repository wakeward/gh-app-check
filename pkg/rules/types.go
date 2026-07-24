// Package rules implements the Phase 1 least-privilege predicates described
// in the design spec's Control Plane Auditor rules engine (Design Spec §4.A).
// Each control lives in its own file with a matching _test.go, mirroring the
// gh-branch-auditor convention.
package rules

// Installation is the minimal installation data the rule predicates need.
// It intentionally mirrors the subset of go-github's Installation struct
// that Phase 1 will populate from GET /orgs/{org}/installations.
type Installation struct {
	RepositorySelection string            // "all" or "selected"
	Permissions         map[string]string // e.g. "administration": "write"
}

// Predicate evaluates one least-privilege rule against an installation.
// It returns true if the installation violates the rule.
type Predicate func(Installation) bool
