package output

import (
	"fmt"
	"strings"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

func appLabel(result eval.AppAuditResult) string {
	if result.AppSlug != "" {
		return result.AppSlug
	}
	if result.AppName != "" {
		return result.AppName
	}
	return "(unknown app)"
}

func formatToxicSummary(matches []eval.ToxicMatch) string {
	if len(matches) == 0 {
		return ""
	}
	parts := make([]string, 0, len(matches))
	for _, match := range matches {
		parts = append(parts, match.Technique)
	}
	return strings.Join(parts, "; ")
}

func formatNearMissSummary(nearMisses []eval.NearMiss) string {
	if len(nearMisses) == 0 {
		return ""
	}
	parts := make([]string, 0, len(nearMisses))
	for _, near := range nearMisses {
		parts = append(parts, fmt.Sprintf("%s (needs %s)", near.Technique, near.MissingGrant))
	}
	return strings.Join(parts, "; ")
}
