package eval

import (
	"github.com/wakeward/gh-app-check/pkg/rules"
	graphmodel "github.com/wakeward/gh-app-graph/pkg/model"
)

// NotableGrant highlights a standalone permission grant whose catalog severity
// is material, even when no toxic combination fully matched.
type NotableGrant struct {
	APIKey        string `json:"api_key"`
	Access        string `json:"access"`
	Severity      string `json:"severity"`
	SecurityNotes string `json:"security_notes,omitempty"`
}

var severityRank = map[graphmodel.Severity]int{
	graphmodel.SeverityInformational: 0,
	graphmodel.SeverityLow:             1,
	graphmodel.SeverityMedium:          2,
	graphmodel.SeverityUnknown:       2,
	graphmodel.SeverityHigh:            3,
	graphmodel.SeverityCritical:        4,
}

// NotableGrants returns catalog entries for installed grants at or above minSeverity.
func NotableGrants(inst rules.Installation, catalog []graphmodel.Permission, minSeverity graphmodel.Severity) []NotableGrant {
	minRank := severityRank[minSeverity]
	byKey := map[string]graphmodel.Permission{}
	for _, perm := range catalog {
		byKey[perm.APIKey] = perm
	}

	var out []NotableGrant
	for apiKey, granted := range inst.Permissions {
		perm, ok := byKey[apiKey]
		if !ok {
			continue
		}
		for _, level := range perm.AccessLevels {
			if !accessSatisfies(granted, level.Access) {
				continue
			}
			// Write implies read in the catalog; surface the write row only.
			if level.Access == graphmodel.AccessRead && granted == string(graphmodel.AccessWrite) {
				continue
			}
			if severityRank[level.Severity] < minRank {
				continue
			}
			out = append(out, NotableGrant{
				APIKey:        apiKey,
				Access:        string(level.Access),
				Severity:      string(level.Severity),
				SecurityNotes: trimNotes(level.SecurityNotes),
			})
		}
	}
	return out
}

func accessSatisfies(granted string, required graphmodel.AccessLevel) bool {
	if required == graphmodel.AccessRead {
		return granted == string(graphmodel.AccessRead) || granted == string(graphmodel.AccessWrite)
	}
	return granted == string(graphmodel.AccessWrite)
}

func trimNotes(notes string) string {
	const maxLen = 500
	if len(notes) <= maxLen {
		return notes
	}
	return notes[:maxLen] + "..."
}
