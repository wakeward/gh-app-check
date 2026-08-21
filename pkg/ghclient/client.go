// Package ghclient wires up an authenticated go-github client shared by the
// Control Plane auditor (Phase 1) and Execution Plane tracer (Phase 2).
package ghclient

import (
	"fmt"
	"strings"

	"github.com/google/go-github/v89/github"
)

// New returns a go-github client for github.com. Prefer NewForHost when the
// target host may be GitHub Enterprise Server.
func New(token string) (*github.Client, error) {
	return NewForHost(token, "github.com")
}

// NewForHost returns an authenticated client for host (github.com or a GHES
// hostname). Pass an empty token to get an unauthenticated client (useful for
// tests against httptest servers via github.WithEnterpriseURLs).
func NewForHost(token, host string) (*github.Client, error) {
	opts := []github.ClientOptionsFunc{}
	if token != "" {
		opts = append(opts, github.WithAuthToken(token))
	}
	host = strings.TrimSpace(host)
	if host == "" {
		host = "github.com"
	}
	if strings.EqualFold(host, "github.com") {
		return github.NewClient(opts...)
	}
	baseURL := fmt.Sprintf("https://%s/api/v3/", strings.TrimSuffix(host, "/"))
	uploadURL := fmt.Sprintf("https://%s/api/uploads/", strings.TrimSuffix(host, "/"))
	opts = append(opts, github.WithEnterpriseURLs(baseURL, uploadURL))
	return github.NewClient(opts...)
}
