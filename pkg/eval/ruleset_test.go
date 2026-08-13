package eval

import (
	"testing"

	graphmodel "github.com/wakeward/gh-app-graph/pkg/model"
	"github.com/wakeward/gh-app-check/pkg/rules"
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

func TestEvaluate(t *testing.T) {
	cases := []struct {
		name          string
		inst          rules.Installation
		toxic         []graphmodel.ToxicCombination
		wantRisk      string
		wantViolCount int
		wantToxic     int
	}{
		{
			name:          "least privilege installation passes",
			inst:          rules.Installation{RepositorySelection: "selected", Permissions: map[string]string{"metadata": "read"}},
			wantRisk:      "PASS",
			wantViolCount: 0,
		},
		{
			name:          "all repos flagged high",
			inst:          rules.Installation{RepositorySelection: "all", Permissions: map[string]string{"metadata": "read"}},
			wantRisk:      "HIGH",
			wantViolCount: 1,
		},
		{
			name:          "administration write flagged critical",
			inst:          rules.Installation{RepositorySelection: "selected", Permissions: map[string]string{"administration": "write"}},
			wantRisk:      "CRITICAL",
			wantViolCount: 1,
		},
		{
			name: "combined violations, critical wins over high",
			inst: rules.Installation{
				RepositorySelection: "all",
				Permissions:         map[string]string{"administration": "write"},
			},
			wantRisk:      "CRITICAL",
			wantViolCount: 2,
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
			wantViolCount: 2, // administration write rule + toxic combo
			wantToxic:     1,
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
			if result.AppSlug != "test-app" || result.Owner != "test-org" {
				t.Errorf("AppSlug/Owner not passed through: got %q/%q", result.AppSlug, result.Owner)
			}
		})
	}
}
