package eval

import (
	"testing"

	"github.com/wakeward/gh-app-check/pkg/rules"
	graphmodel "github.com/wakeward/gh-app-graph/pkg/model"
)

func TestNotableGrants_IncludesHighWriteGrants(t *testing.T) {
	catalog := []graphmodel.Permission{{
		APIKey: "contents",
		AccessLevels: []graphmodel.AccessLevelDetail{{
			Access:        graphmodel.AccessWrite,
			Severity:      graphmodel.SeverityHigh,
			SecurityNotes: "Can modify repository content.",
		}},
	}, {
		APIKey: "metadata",
		AccessLevels: []graphmodel.AccessLevelDetail{{
			Access:   graphmodel.AccessRead,
			Severity: graphmodel.SeverityLow,
		}},
	}}

	got := NotableGrants(rules.Installation{
		Permissions: map[string]string{"contents": "write", "metadata": "read"},
	}, catalog, graphmodel.SeverityHigh)

	if len(got) != 1 || got[0].APIKey != "contents" {
		t.Fatalf("unexpected notable grants: %+v", got)
	}
}

func TestEvaluateWithContext_IncludesExploitPath(t *testing.T) {
	combo := graphmodel.ToxicCombination{
		ID:          "arbitrary-code-execution",
		Technique:   "Arbitrary Code Execution",
		BlastRadius: graphmodel.BlastRadiusCritical,
		ExploitPath: "Runs malicious workflow steps.",
		Permissions: []graphmodel.PermissionGrant{
			{APIKey: "workflows", Access: graphmodel.AccessWrite},
			{APIKey: "contents", Access: graphmodel.AccessWrite},
		},
	}
	result := EvaluateWithContext("app", "App", "org", rules.Installation{
		Permissions: map[string]string{"workflows": "write", "contents": "write"},
	}, []graphmodel.ToxicCombination{combo}, ScanContext{IncludeGHESRules: true})

	if len(result.ToxicMatches) != 1 {
		t.Fatalf("expected one toxic match, got %+v", result.ToxicMatches)
	}
	if result.ToxicMatches[0].ExploitPath == "" {
		t.Fatal("expected exploit_path on toxic match")
	}
	if len(result.ToxicMatches[0].MatchedGrants) != 2 {
		t.Fatalf("expected matched grants, got %+v", result.ToxicMatches[0].MatchedGrants)
	}
	if result.RiskLevel != "CRITICAL" {
		t.Fatalf("risk=%s", result.RiskLevel)
	}
}
