// Package ghclient wires up an authenticated go-github client shared by the
// Control Plane auditor (Phase 1) and Execution Plane tracer (Phase 2).
package ghclient

import "github.com/google/go-github/v89/github"

// New returns a go-github client. Pass an empty token to get an
// unauthenticated client (useful for tests against a local httptest server
// with go-github's BaseURL redirected via github.WithEnterpriseURLs).
func New(token string) (*github.Client, error) {
	opts := []github.ClientOptionsFunc{}
	if token != "" {
		opts = append(opts, github.WithAuthToken(token))
	}
	return github.NewClient(opts...)
}
