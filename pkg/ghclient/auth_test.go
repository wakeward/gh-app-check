package ghclient

import (
	"strings"
	"testing"
)

func TestResolveScanPlatform(t *testing.T) {
	cases := []struct {
		flag   string
		detect ScanPlatform
		want   ScanPlatform
	}{
		{"auto", ScanPlatformCloud, ScanPlatformCloud},
		{"", ScanPlatformGHES, ScanPlatformGHES},
		{"cloud", ScanPlatformGHES, ScanPlatformCloud},
		{"ghes", ScanPlatformCloud, ScanPlatformGHES},
	}
	for _, tc := range cases {
		got, err := ResolveScanPlatform(tc.flag, tc.detect)
		if err != nil {
			t.Fatalf("flag=%q: %v", tc.flag, err)
		}
		if got != tc.want {
			t.Errorf("flag=%q detect=%q: got %q want %q", tc.flag, tc.detect, got, tc.want)
		}
	}
}

func TestScanPlatformLabel(t *testing.T) {
	if got := ScanPlatformCloud.Label("github.com"); got != "cloud (github.com)" {
		t.Errorf("got %q", got)
	}
	if got := ScanPlatformGHES.Label("github.example.com"); got != "ghes (github.example.com)" {
		t.Errorf("got %q", got)
	}
}

func TestNewForHost_GHESEnterpriseURLs(t *testing.T) {
	client, err := NewForHost("", "github.example.com")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(client.BaseURL(), "github.example.com") {
		t.Fatalf("BaseURL = %q", client.BaseURL())
	}
}
