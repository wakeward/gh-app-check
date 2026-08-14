package ghclient

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-github/v89/github"
)

func TestPermissionsMap_FlattensStruct(t *testing.T) {
	read := "read"
	write := "write"
	got, err := PermissionsMap(&github.InstallationPermissions{
		Metadata:       &read,
		Administration: &write,
	})
	if err != nil {
		t.Fatal(err)
	}
	if got["metadata"] != "read" || got["administration"] != "write" {
		t.Fatalf("unexpected map: %#v", got)
	}
}

func TestListOrgInstallations_Paginates(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v3/orgs/acme/installations", func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Query().Get("page") {
		case "", "1":
			w.Header().Set("Link", `<https://example.com/orgs/acme/installations?page=2>; rel="next"`)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"total_count": 2,
				"installations": []map[string]any{{
					"id":                   1,
					"repository_selection": "selected",
					"app_slug":             "app-one",
					"permissions":          map[string]any{"metadata": "read"},
				}},
			})
		case "2":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"total_count": 2,
				"installations": []map[string]any{{
					"id":                   2,
					"repository_selection": "all",
					"app_slug":             "app-two",
					"permissions":          map[string]any{"administration": "write"},
				}},
			})
		default:
			http.NotFound(w, r)
		}
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

	installations, err := ListOrgInstallations(context.Background(), client, "acme")
	if err != nil {
		t.Fatal(err)
	}
	if len(installations) != 2 {
		t.Fatalf("expected 2 installations, got %d", len(installations))
	}
	if installations[0].Slug != "app-one" || installations[1].RepositorySelection != "all" {
		t.Fatalf("unexpected installations: %+v", installations)
	}
}
