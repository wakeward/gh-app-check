package eval

import "github.com/wakeward/gh-app-check/pkg/rules"

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

// Evaluate runs every rule in the ruleset against inst and returns the
// aggregated AppAuditResult. The highest risk level among triggered rules
// wins; the result is "PASS" if none trigger.
func Evaluate(appSlug, owner string, inst rules.Installation) AppAuditResult {
	result := AppAuditResult{
		AppSlug:       appSlug,
		Owner:         owner,
		RepoSelection: inst.RepositorySelection,
		RiskLevel:     "PASS",
		Violations:    []string{},
	}

	for _, rule := range ControlPlaneRuleset() {
		if rule.Predicate(inst) {
			result.Violations = append(result.Violations, rule.Message)
			if riskRank[rule.Risk] > riskRank[result.RiskLevel] {
				result.RiskLevel = rule.Risk
			}
		}
	}

	return result
}
