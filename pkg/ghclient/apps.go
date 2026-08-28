package ghclient

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v89/github"
)

// EnrichInstallationNames fills OrgInstallation.Name from GET /apps/{slug}
// when the installation payload only includes app_slug. Slugs are deduplicated
// so each app is fetched at most once. Lookup failures leave the slug in place.
func EnrichInstallationNames(ctx context.Context, client *github.Client, installations []OrgInstallation) error {
	if client == nil {
		return fmt.Errorf("github client is nil")
	}
	if len(installations) == 0 {
		return nil
	}

	names, err := FetchAppNames(ctx, client, installationSlugs(installations))
	if err != nil {
		return err
	}
	for i := range installations {
		if name, ok := names[installations[i].Slug]; ok && name != "" {
			installations[i].Name = name
		}
	}
	return nil
}

// FetchAppNames resolves display names for the given app slugs via the public
// app metadata endpoint.
func FetchAppNames(ctx context.Context, client *github.Client, slugs []string) (map[string]string, error) {
	out := make(map[string]string, len(slugs))
	for _, slug := range slugs {
		slug = strings.TrimSpace(slug)
		if slug == "" {
			continue
		}
		if _, ok := out[slug]; ok {
			continue
		}
		app, _, err := client.Apps.Get(ctx, slug)
		if err != nil {
			continue
		}
		if app == nil {
			continue
		}
		name := strings.TrimSpace(app.GetName())
		if name != "" {
			out[slug] = name
		}
	}
	return out, nil
}

func installationSlugs(installations []OrgInstallation) []string {
	seen := make(map[string]struct{}, len(installations))
	slugs := make([]string, 0, len(installations))
	for _, inst := range installations {
		if inst.Slug == "" {
			continue
		}
		if _, ok := seen[inst.Slug]; ok {
			continue
		}
		seen[inst.Slug] = struct{}{}
		slugs = append(slugs, inst.Slug)
	}
	return slugs
}
