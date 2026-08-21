package ghclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v89/github"
)

func TestEnrichInstallationNames(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/apps/dependabot", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"slug": "dependabot",
			"name": "Dependabot",
		})
	})
	mux.HandleFunc("/api/v3/apps/missing-app", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client, err := github.NewClient(
		github.WithHTTPClient(server.Client()),
		github.WithEnterpriseURLs(server.URL+"/", server.URL+"/"),
	)
	if err != nil {
		t.Fatal(err)
	}

	installations := []OrgInstallation{
		{Slug: "dependabot", Name: "dependabot"},
		{Slug: "dependabot", Name: "dependabot"},
		{Slug: "missing-app", Name: "missing-app"},
	}
	if err := EnrichInstallationNames(context.Background(), client, installations); err != nil {
		t.Fatal(err)
	}
	if installations[0].Name != "Dependabot" {
		t.Errorf("installation 0 name = %q, want Dependabot", installations[0].Name)
	}
	if installations[1].Name != "Dependabot" {
		t.Errorf("installation 1 name = %q, want Dependabot", installations[1].Name)
	}
	if installations[2].Name != "missing-app" {
		t.Errorf("installation 2 name = %q, want slug fallback", installations[2].Name)
	}
}
