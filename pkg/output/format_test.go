package output

import (
	"testing"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

func TestAppLabel_NameAndSlug(t *testing.T) {
	label := appLabel(eval.AppAuditResult{
		AppName: "Dependabot",
		AppSlug: "dependabot",
	})
	if label != "Dependabot (dependabot)" {
		t.Fatalf("got %q", label)
	}
}

func TestAppLabel_SlugOnly(t *testing.T) {
	label := appLabel(eval.AppAuditResult{AppSlug: "my-app"})
	if label != "my-app" {
		t.Fatalf("got %q", label)
	}
}

func TestFormatToxicSummary_DeduplicatesTechniqueWithID(t *testing.T) {
	got := formatToxicSummary([]eval.ToxicMatch{
		{ID: "combo-a", Technique: "CI/CD Pipeline Takeover"},
		{ID: "combo-b", Technique: "CI/CD Pipeline Takeover"},
	})
	if got != "CI/CD Pipeline Takeover [combo-a]; CI/CD Pipeline Takeover [combo-b]" {
		t.Fatalf("got %q", got)
	}
}

func TestFormatNearMissSummary_Deduplicates(t *testing.T) {
	got := formatNearMissSummary([]eval.NearMiss{
		{Technique: "Real-Time Data Tap", MissingGrant: "organization_hooks:write"},
		{Technique: "Real-Time Data Tap", MissingGrant: "organization_hooks:write"},
	})
	if got != "Real-Time Data Tap (needs organization_hooks:write)" {
		t.Fatalf("got %q", got)
	}
}
