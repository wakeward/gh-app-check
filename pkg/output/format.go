package output

import (
	"fmt"
	"strings"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

func appLabel(result eval.AppAuditResult) string {
	if result.AppName != "" && result.AppSlug != "" && result.AppName != result.AppSlug {
		return fmt.Sprintf("%s (%s)", result.AppName, result.AppSlug)
	}
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
	techniqueCount := map[string]int{}
	for _, match := range matches {
		techniqueCount[match.Technique]++
	}
	parts := make([]string, 0, len(matches))
	for _, match := range matches {
		label := match.Technique
		if techniqueCount[match.Technique] > 1 && match.ID != "" {
			label = fmt.Sprintf("%s [%s]", match.Technique, match.ID)
		}
		parts = append(parts, label)
	}
	return strings.Join(parts, "; ")
}

func formatNearMissSummary(nearMisses []eval.NearMiss) string {
	if len(nearMisses) == 0 {
		return ""
	}
	seen := map[string]struct{}{}
	parts := make([]string, 0, len(nearMisses))
	for _, near := range nearMisses {
		key := near.Technique + "\x00" + near.MissingGrant
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		parts = append(parts, fmt.Sprintf("%s (needs %s)", near.Technique, near.MissingGrant))
	}
	return strings.Join(parts, "; ")
}
