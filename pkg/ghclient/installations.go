package ghclient

import (
	"context"
	"fmt"

	"github.com/google/go-github/v89/github"
)

// OrgInstallation is the subset of a GitHub App installation needed for
// Phase 1 control-plane auditing.
type OrgInstallation struct {
	Slug                string
	Name                string
	RepositorySelection string
	Permissions         map[string]string
}

// ListOrgInstallations pages through GET /orgs/{org}/installations.
func ListOrgInstallations(ctx context.Context, client *github.Client, org string) ([]OrgInstallation, error) {
	if client == nil {
		return nil, fmt.Errorf("github client is nil")
	}
	var out []OrgInstallation
	opt := &github.ListOptions{PerPage: 100}
	for {
		result, resp, err := client.Organizations.ListInstallations(ctx, org, opt)
		if err != nil {
			return nil, fmt.Errorf("list installations for org %q: %w", org, err)
		}
		if result == nil {
			break
		}
		for _, inst := range result.Installations {
			if inst == nil {
				continue
			}
			slug := inst.GetAppSlug()
			name := slug
			if name == "" && inst.AppID != nil {
				name = fmt.Sprintf("app-%d", inst.GetAppID())
			}
			perms, err := PermissionsMap(inst.Permissions)
			if err != nil {
				return nil, fmt.Errorf("installation %q: %w", slug, err)
			}
			out = append(out, OrgInstallation{
				Slug:                slug,
				Name:                name,
				RepositorySelection: inst.GetRepositorySelection(),
				Permissions:         perms,
			})
		}
		if resp == nil || resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return out, nil
}
