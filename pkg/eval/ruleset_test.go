package eval

import (
	"testing"

	"github.com/wakeward/gh-app-check/pkg/rules"
	graphmodel "github.com/wakeward/gh-app-graph/pkg/model"
)

func stealthBackdoorCombo() graphmodel.ToxicCombination {
	return graphmodel.ToxicCombination{
		ID:          "stealth-backdoor",
		Technique:   "Stealth Backdoor",
		BlastRadius: graphmodel.BlastRadiusCritical,
		Permissions: []graphmodel.PermissionGrant{
			{APIKey: "administration", Access: graphmodel.AccessWrite},
			{APIKey: "contents", Access: graphmodel.AccessWrite},
		},
	}
}

func aceCombo() graphmodel.ToxicCombination {
	return graphmodel.ToxicCombination{
		ID:          "arbitrary-code-execution",
		Technique:   "Arbitrary Code Execution",
		BlastRadius: graphmodel.BlastRadiusCritical,
		Permissions: []graphmodel.PermissionGrant{
			{APIKey: "workflows", Access: graphmodel.AccessWrite},
			{APIKey: "contents", Access: graphmodel.AccessWrite},
		},
	}
}

func TestEvaluate(t *testing.T) {
	cases := []struct {
		name          string
		inst          rules.Installation
		toxic         []graphmodel.ToxicCombination
		wantRisk      string
		wantViolCount int
		wantToxic     int
		wantNearMiss  int
		wantWrites    int
	}{
		{
			name:          "least privilege installation passes",
			inst:          rules.Installation{RepositorySelection: "selected", Permissions: map[string]string{"metadata": "read"}},
			wantRisk:      "PASS",
			wantViolCount: 0,
		},
		{
			name:          "all repos read-only flagged warn",
			inst:          rules.Installation{RepositorySelection: "all", Permissions: map[string]string{"metadata": "read", "contents": "read"}},
			wantRisk:      "WARN",
			wantViolCount: 1,
		},
		{
			name:          "all repos with write flagged high",
			inst:          rules.Installation{RepositorySelection: "all", Permissions: map[string]string{"metadata": "read", "contents": "write"}},
			wantRisk:      "HIGH",
			wantViolCount: 1,
			wantWrites:    1,
		},
		{
			name:          "administration write flagged critical",
			inst:          rules.Installation{RepositorySelection: "selected", Permissions: map[string]string{"administration": "write"}},
			wantRisk:      "CRITICAL",
			wantViolCount: 1,
			wantWrites:    1,
		},
		{
			name: "combined violations, critical wins over high",
			inst: rules.Installation{
				RepositorySelection: "all",
				Permissions:         map[string]string{"administration": "write"},
			},
			wantRisk:      "CRITICAL",
			wantViolCount: 2,
			wantWrites:    1,
		},
		{
			name: "toxic combination match elevates risk",
			inst: rules.Installation{
				RepositorySelection: "selected",
				Permissions: map[string]string{
					"administration": "write",
					"contents":       "write",
				},
			},
			toxic:         []graphmodel.ToxicCombination{stealthBackdoorCombo()},
			wantRisk:      "CRITICAL",
			wantViolCount: 2,
			wantToxic:     1,
			wantWrites:    2,
		},
		{
			name: "near miss does not raise risk to toxic blast radius",
			inst: rules.Installation{
				RepositorySelection: "selected",
				Permissions: map[string]string{
					"workflows": "write",
					"contents":  "read",
				},
			},
			toxic:        []graphmodel.ToxicCombination{aceCombo()},
			wantRisk:     "PASS",
			wantNearMiss: 1,
			wantWrites:   1,
		},
		{
			name: "god-mode is high not critical without toxic match",
			inst: rules.Installation{
				RepositorySelection: "selected",
				Permissions: map[string]string{
					"a": "write", "b": "write", "c": "write",
					"d": "write", "e": "write", "f": "write",
				},
			},
			wantRisk:      "HIGH",
			wantViolCount: 1,
			wantWrites:    6,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Evaluate("test-app", "Test App", "test-org", tc.inst, tc.toxic)
			if result.RiskLevel != tc.wantRisk {
				t.Errorf("RiskLevel = %q, want %q", result.RiskLevel, tc.wantRisk)
			}
			if len(result.Violations) != tc.wantViolCount {
				t.Errorf("len(Violations) = %d, want %d (%v)", len(result.Violations), tc.wantViolCount, result.Violations)
			}
			if len(result.ToxicMatches) != tc.wantToxic {
				t.Errorf("len(ToxicMatches) = %d, want %d", len(result.ToxicMatches), tc.wantToxic)
			}
			if len(result.NearMisses) != tc.wantNearMiss {
				t.Errorf("len(NearMisses) = %d, want %d (%v)", len(result.NearMisses), tc.wantNearMiss, result.NearMisses)
			}
			if tc.wantWrites > 0 && result.WriteScopeCount != tc.wantWrites {
				t.Errorf("WriteScopeCount = %d, want %d", result.WriteScopeCount, tc.wantWrites)
			}
			if result.AppSlug != "test-app" || result.Owner != "test-org" {
				t.Errorf("AppSlug/Owner not passed through: got %q/%q", result.AppSlug, result.Owner)
			}
		})
	}
}

func TestRiskRank(t *testing.T) {
	if RiskRank("CRITICAL") <= RiskRank("HIGH") {
		t.Fatal("CRITICAL should rank above HIGH")
	}
}
