package output

import (
	"bytes"
	"encoding/json"
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
		{AppSlug: "test-app", Owner: "test-org", RepoSelection: "all", RiskLevel: "HIGH", Violations: []string{"installation has access to all repositories"}},
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
}
