package eval

import (
	"testing"

	"github.com/wakeward/gh-app-check/pkg/rules"
)

func TestEvaluate(t *testing.T) {
	cases := []struct {
		name          string
		inst          rules.Installation
		wantRisk      string
		wantViolCount int
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
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result := Evaluate("test-app", "test-org", tc.inst)
			if result.RiskLevel != tc.wantRisk {
				t.Errorf("RiskLevel = %q, want %q", result.RiskLevel, tc.wantRisk)
			}
			if len(result.Violations) != tc.wantViolCount {
				t.Errorf("len(Violations) = %d, want %d (%v)", len(result.Violations), tc.wantViolCount, result.Violations)
			}
			if result.AppSlug != "test-app" || result.Owner != "test-org" {
				t.Errorf("AppSlug/Owner not passed through: got %q/%q", result.AppSlug, result.Owner)
			}
		})
	}
}
