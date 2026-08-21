package ghclient

import (
	"fmt"
	"os"
	"strings"

	ghauth "github.com/cli/go-gh/v2/pkg/auth"
)

// ScanPlatform identifies the GitHub deployment target for rule selection.
type ScanPlatform string

const (
	ScanPlatformCloud ScanPlatform = "cloud"
	ScanPlatformGHES  ScanPlatform = "ghes"
)

// AuthContext holds the resolved GitHub host and whether it is Enterprise Server.
type AuthContext struct {
	Host     string
	Token    string
	Platform ScanPlatform
}

// ResolveAuth returns credentials and scan platform from GH_HOST/GITHUB_TOKEN or
// the gh CLI configuration. Enterprise Cloud (github.com) maps to cloud; self-
// hosted Enterprise Server maps to ghes (see go-gh auth.IsEnterprise).
func ResolveAuth() (AuthContext, error) {
	host, _ := ghauth.DefaultHost()
	if host == "" {
		host = "github.com"
	}
	token, _ := ghauth.TokenForHost(host)
	if token == "" {
		if env := strings.TrimSpace(os.Getenv("GH_TOKEN")); env != "" {
			token = env
		} else if env := strings.TrimSpace(os.Getenv("GITHUB_TOKEN")); env != "" {
			token = env
		}
	}
	if token == "" {
		return AuthContext{}, fmt.Errorf("no GitHub token found for host %q: set GH_TOKEN/GITHUB_TOKEN or run `gh auth login`", host)
	}
	platform := ScanPlatformCloud
	if ghauth.IsEnterprise(host) {
		platform = ScanPlatformGHES
	}
	return AuthContext{
		Host:     host,
		Token:    token,
		Platform: platform,
	}, nil
}

// ResolveScanPlatform interprets an explicit --platform flag together with autodetected auth.
func ResolveScanPlatform(flagValue string, detected ScanPlatform) (ScanPlatform, error) {
	switch strings.ToLower(strings.TrimSpace(flagValue)) {
	case "", "auto":
		return detected, nil
	case "cloud", "github.com":
		return ScanPlatformCloud, nil
	case "ghes", "enterprise-server":
		return ScanPlatformGHES, nil
	default:
		return "", fmt.Errorf("invalid platform %q: must be auto, cloud, or ghes", flagValue)
	}
}

func (p ScanPlatform) IncludesGHESRules() bool {
	return p == ScanPlatformGHES
}

func (p ScanPlatform) Label(host string) string {
	switch p {
	case ScanPlatformGHES:
		return "ghes (" + host + ")"
	default:
		if host == "" || host == "github.com" {
			return "cloud (github.com)"
		}
		return "cloud (" + host + ")"
	}
}
