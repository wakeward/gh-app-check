package eval

import (
	"fmt"

	grapheval "github.com/wakeward/gh-app-graph/pkg/eval"
	graphmodel "github.com/wakeward/gh-app-graph/pkg/model"
	"github.com/wakeward/gh-app-check/pkg/rules"
)

// ToxicMatch records one toxic combination satisfied by an installation.
type ToxicMatch struct {
	ID        string `json:"id"`
	Technique string `json:"technique"`
	Blast     string `json:"blast_radius"`
}

// AppAuditResult is the outcome of evaluating one GitHub App installation
// against the least-privilege rules engine and toxic-combination catalog.
type AppAuditResult struct {
	AppSlug       string       `json:"app_slug"`
	AppName       string       `json:"app_name,omitempty"`
	Owner         string       `json:"owner"`
	RepoSelection string       `json:"repo_selection"` // "all" or "selected"
	RiskLevel     string       `json:"risk_level"`     // "CRITICAL", "HIGH", "WARN", "PASS"
	Violations    []string     `json:"violations"`
	ToxicMatches  []ToxicMatch `json:"toxic_matches,omitempty"`
}

// namedRule pairs a rule predicate with the violation message to record and
// the risk level it contributes when triggered.
type namedRule struct {
	Predicate rules.Predicate
	Message   string
	Risk      string
}

// riskRank orders risk levels so the highest-severity triggered rule wins.
var riskRank = map[string]int{"PASS": 0, "WARN": 1, "HIGH": 2, "CRITICAL": 3}

// ControlPlaneRuleset is the Phase 1 least-privilege rules engine described
// in Design Spec §4.A.
func ControlPlaneRuleset() []namedRule {
	return []namedRule{
		{rules.RepoSelectionAll, "installation has access to all repositories", "HIGH"},
		{rules.AdministrationWrite, "installation has write access to administration", "CRITICAL"},
		{rules.ExcessiveWriteScopes, "installation has more than 5 write-level scopes (god-mode)", "CRITICAL"},
	}
}

// Evaluate runs the control-plane ruleset and toxic-combination catalog
// against inst and returns the aggregated AppAuditResult.
func Evaluate(appSlug, appName, owner string, inst rules.Installation, toxic []graphmodel.ToxicCombination) AppAuditResult {
	result := AppAuditResult{
		AppSlug:       appSlug,
		AppName:       appName,
		Owner:         owner,
		RepoSelection: inst.RepositorySelection,
		RiskLevel:     "PASS",
		Violations:    []string{},
	}

	for _, rule := range ControlPlaneRuleset() {
		if rule.Predicate(inst) {
			result.Violations = append(result.Violations, rule.Message)
			raiseRisk(&result, rule.Risk)
		}
	}

	if len(toxic) > 0 {
		toxicResult := grapheval.Evaluate(graphmodel.AppPermissionSet{Permissions: inst.Permissions}, toxic)
		for _, match := range toxicResult.Matches {
			result.ToxicMatches = append(result.ToxicMatches, ToxicMatch{
				ID:        match.Combination.ID,
				Technique: match.Combination.Technique,
				Blast:     string(match.Combination.BlastRadius),
			})
			result.Violations = append(result.Violations, fmt.Sprintf(
				"toxic combination: %s (%s)",
				match.Combination.Technique,
				match.Combination.BlastRadius,
			))
			raiseRisk(&result, blastToRisk(match.Combination.BlastRadius))
		}
	}

	return result
}

func raiseRisk(result *AppAuditResult, risk string) {
	if riskRank[risk] > riskRank[result.RiskLevel] {
		result.RiskLevel = risk
	}
}

func blastToRisk(blast graphmodel.BlastRadius) string {
	switch blast {
	case graphmodel.BlastRadiusCritical:
		return "CRITICAL"
	case graphmodel.BlastRadiusHigh:
		return "HIGH"
	case graphmodel.BlastRadiusMedium:
		return "WARN"
	default:
		return "PASS"
	}
}
