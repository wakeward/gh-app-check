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
