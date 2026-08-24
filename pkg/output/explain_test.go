package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

func TestExplainWriter_ShowsCriticalWithRationale(t *testing.T) {
	var buf bytes.Buffer
	err := (ExplainWriter{}).WriteOrgScan(&buf, eval.OrgScanResult{
		ScanPlatform: "cloud (github.com)",
		Installations: []eval.AppAuditResult{{
			AppSlug:       "ci-app",
			RiskLevel:     "CRITICAL",
			RepoSelection: "all",
			WriteScopeCount: 7,
			ToxicMatches: []eval.ToxicMatch{{
				ID:        "ace",
				Technique: "Arbitrary Code Execution",
				Blast:     "Critical",
				ExploitPath: "Runs malicious workflow steps.",
				MatchedGrants: []eval.PermissionGrant{
					{APIKey: "workflows", Access: "write"},
					{APIKey: "contents", Access: "write"},
				},
			}},
			ControlPlaneFindings: []eval.ControlPlaneFinding{{
				Risk:      "HIGH",
				Message:   "installation has access to all repositories",
				Rationale: "Org-wide reach.",
			}},
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	for _, want := range []string{
		"ci-app — CRITICAL",
		"Arbitrary Code Execution",
		"Matched grants:",
		"workflows: write",
		"What this enables:",
		"Runs malicious workflow steps.",
		"Structural findings:",
		"Org-wide reach.",
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in output:\n%s", want, out)
		}
	}
}

func TestExplainWriter_SkipsPassUnlessExplainAll(t *testing.T) {
	scan := eval.OrgScanResult{
		Installations: []eval.AppAuditResult{{
			AppSlug:   "quiet",
			RiskLevel: "PASS",
		}},
	}
	var buf bytes.Buffer
	if err := (ExplainWriter{}).WriteOrgScan(&buf, scan); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(buf.String(), "quiet") {
		t.Fatalf("expected PASS app hidden: %s", buf.String())
	}

	buf.Reset()
	if err := (ExplainWriter{ExplainAll: true}).WriteOrgScan(&buf, scan); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "quiet") {
		t.Fatal("expected PASS app with --explain-all")
	}
}

func TestWriteWrapped_NearMissContinuationDoesNotRepeatLabel(t *testing.T) {
	var buf strings.Builder
	long := strings.Repeat("word ", 30)
	writeWrapped(&buf, "    ", "Would enable: ", long)
	out := buf.String()
	if strings.Count(out, "Would enable:") != 1 {
		t.Fatalf("label repeated in wrapped near-miss text:\n%s", out)
	}
}
