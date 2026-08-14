package output

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

func TestForFormat(t *testing.T) {
	cases := []struct {
		format  string
		wantErr bool
	}{
		{"table", false},
		{"json", false},
		{"markdown", false},
		{"yaml", true},
	}

	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			_, err := ForFormat(tc.format)
			if (err != nil) != tc.wantErr {
				t.Errorf("ForFormat(%q) err = %v, wantErr = %v", tc.format, err, tc.wantErr)
			}
		})
	}
}

func TestJSONWriterRoundTrip(t *testing.T) {
	results := []eval.AppAuditResult{
		{
			AppSlug:         "test-app",
			Owner:           "test-org",
			RepoSelection:   "all",
			WriteScopeCount: 2,
			RiskLevel:       "HIGH",
			Violations:      []string{"installation has access to all repositories"},
			ToxicMatches:    []eval.ToxicMatch{},
			NearMisses:      []eval.NearMiss{},
		},
	}

	var buf bytes.Buffer
	if err := (JSONWriter{}).Write(&buf, results); err != nil {
		t.Fatalf("Write() error = %v", err)
	}

	var decoded []eval.AppAuditResult
	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("output is not valid JSON: %v\n%s", err, buf.String())
	}

	if len(decoded) != 1 || decoded[0].AppSlug != "test-app" || decoded[0].RiskLevel != "HIGH" {
		t.Errorf("round-tripped result mismatch: %+v", decoded)
	}
	if decoded[0].WriteScopeCount != 2 {
		t.Errorf("WriteScopeCount = %d, want 2", decoded[0].WriteScopeCount)
	}
	if decoded[0].ToxicMatches == nil || decoded[0].NearMisses == nil {
		t.Errorf("expected empty toxic_matches and near_misses arrays, got %+v", decoded[0])
	}
}

func TestSortResultsOrdersByRisk(t *testing.T) {
	results := []eval.AppAuditResult{
		{AppSlug: "beta", RiskLevel: "PASS"},
		{AppSlug: "alpha", RiskLevel: "CRITICAL"},
		{AppSlug: "gamma", RiskLevel: "HIGH"},
	}
	SortResults(results)
	if results[0].AppSlug != "alpha" || results[1].AppSlug != "gamma" || results[2].AppSlug != "beta" {
		t.Fatalf("unexpected order: %+v", results)
	}
}

func TestTableWriterIncludesNearMissColumn(t *testing.T) {
	var buf bytes.Buffer
	err := (TableWriter{}).Write(&buf, []eval.AppAuditResult{{
		AppSlug:   "demo",
		RiskLevel: "PASS",
		NearMisses: []eval.NearMiss{{
			ID:           "arbitrary-code-execution",
			Technique:    "Arbitrary Code Execution",
			MissingGrant: "contents:write",
		}},
	}})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "NEAR_MISSES") {
		t.Fatalf("table missing NEAR_MISSES header: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "contents:write") {
		t.Fatalf("table missing near-miss detail: %s", buf.String())
	}
}
