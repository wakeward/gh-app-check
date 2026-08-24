package eval

import (
	"fmt"

	"github.com/wakeward/gh-app-check/pkg/rules"
	grapheval "github.com/wakeward/gh-app-graph/pkg/eval"
	graphplatform "github.com/wakeward/gh-app-graph/pkg/platform"
	graphmodel "github.com/wakeward/gh-app-graph/pkg/model"
)

// ToxicMatch records one toxic combination satisfied by an installation.
type ToxicMatch struct {
	ID        string `json:"id"`
	Technique string `json:"technique"`
	Blast     string `json:"blast_radius"`
}

// NearMiss records a toxic combination one grant away from being satisfied.
type NearMiss struct {
	ID           string `json:"id"`
	Technique    string `json:"technique"`
	MissingGrant string `json:"missing_grant"`
}

// AppAuditResult is the outcome of evaluating one GitHub App installation
// against the least-privilege rules engine and toxic-combination catalog.
type AppAuditResult struct {
	InstallationID  int64             `json:"installation_id,omitempty"`
	AppID           int64             `json:"app_id,omitempty"`
	AppSlug         string            `json:"app_slug"`
	AppName         string            `json:"app_name,omitempty"`
	HTMLURL         string            `json:"html_url,omitempty"`
	Owner           string            `json:"owner"`
	RepoSelection   string            `json:"repo_selection"` // "all" or "selected"
	Permissions     map[string]string `json:"permissions,omitempty"`
	WriteScopeCount int               `json:"write_scope_count"`
	RiskLevel       string            `json:"risk_level"` // "CRITICAL", "HIGH", "WARN", "PASS"
	Violations      []string          `json:"violations"`
	ToxicMatches    []ToxicMatch      `json:"toxic_matches"`
	NearMisses      []NearMiss        `json:"near_misses"`
	GHESScopes      []string          `json:"ghes_scopes,omitempty"`
}

// OrgScanResult is the full output of an organization audit including scan metadata.
type OrgScanResult struct {
	ScanPlatform      string           `json:"scan_platform"`
	ScanHost            string           `json:"scan_host"`
	ExcludedGHESRules   int              `json:"excluded_ghes_rules"`
	Installations       []AppAuditResult `json:"installations"`
}

// ScanContext configures platform-aware evaluation.
type ScanContext struct {
	IncludeGHESRules bool
	GHESOnlyKeys     map[string]struct{}
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

// RiskRank returns the sort order for a risk level string.
func RiskRank(level string) int {
	return riskRank[level]
}

// ControlPlaneRuleset is the Phase 1 least-privilege rules engine described
// in Design Spec §4.A.
func ControlPlaneRuleset() []namedRule {
	return []namedRule{
		{rules.RepoSelectionAllWithWrite, "installation has access to all repositories", "HIGH"},
		{rules.RepoSelectionAllReadOnly, "installation has access to all repositories (read-only grants)", "WARN"},
		{rules.AdministrationWrite, "installation has write access to administration", "CRITICAL"},
		{rules.ExcessiveWriteScopes, "installation has more than 5 write-level scopes (god-mode)", "HIGH"},
	}
}

// Evaluate runs the control-plane ruleset and toxic-combination catalog
// against inst and returns the aggregated AppAuditResult.
func Evaluate(appSlug, appName, owner string, inst rules.Installation, toxic []graphmodel.ToxicCombination) AppAuditResult {
	return EvaluateWithContext(appSlug, appName, owner, inst, toxic, ScanContext{IncludeGHESRules: true})
}

// EvaluateWithContext applies platform-aware toxic rule filtering and highlights
// GHES-only scopes when present.
func EvaluateWithContext(appSlug, appName, owner string, inst rules.Installation, toxic []graphmodel.ToxicCombination, scan ScanContext) AppAuditResult {
	if !scan.IncludeGHESRules && len(scan.GHESOnlyKeys) > 0 {
		toxic = graphplatform.FilterToxicCombinations(toxic, scan.GHESOnlyKeys, false)
	}

	result := AppAuditResult{
		AppSlug:         appSlug,
		AppName:         appName,
		Owner:           owner,
		RepoSelection:   inst.RepositorySelection,
		WriteScopeCount: rules.WriteScopeCount(inst),
		RiskLevel:       "PASS",
		Violations:      []string{},
		ToxicMatches:    []ToxicMatch{},
		NearMisses:      []NearMiss{},
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
			msg := fmt.Sprintf(
				"toxic combination: %s (%s)",
				match.Combination.Technique,
				match.Combination.BlastRadius,
			)
			if !stringSliceContains(result.Violations, msg) {
				result.Violations = append(result.Violations, msg)
			}
			raiseRisk(&result, blastToRisk(match.Combination.BlastRadius))
		}
		for _, near := range toxicResult.NearMisses {
			result.NearMisses = append(result.NearMisses, NearMiss{
				ID:           near.Combination.ID,
				Technique:    near.Combination.Technique,
				MissingGrant: formatMissingGrant(near.Missing[0]),
			})
		}
	}

	if ghesScopes := graphplatform.GrantedGHESOnly(inst.Permissions, scan.GHESOnlyKeys); len(ghesScopes) > 0 {
		result.GHESScopes = ghesScopes
		for _, scope := range ghesScopes {
			result.Violations = append(result.Violations, "GHES-only scope granted: "+scope)
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

func formatMissingGrant(grant graphmodel.PermissionGrant) string {
	return fmt.Sprintf("%s:%s", grant.APIKey, grant.Access)
}

func stringSliceContains(items []string, target string) bool {
	for _, item := range items {
		if item == target {
			return true
		}
	}
	return false
}
