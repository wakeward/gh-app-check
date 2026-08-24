package output

import (
	"fmt"
	"io"
	"strings"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// ExplainWriter renders narrative findings for engineers triaging risk levels.
type ExplainWriter struct {
	ExplainAll bool
}

// WriteOrgScan renders explain output for each installation in the scan.
func (ew ExplainWriter) WriteOrgScan(w io.Writer, scan eval.OrgScanResult) error {
	if scan.ScanPlatform != "" {
		line := fmt.Sprintf("# scan: %s", scan.ScanPlatform)
		if scan.ExcludedGHESRules > 0 {
			line += fmt.Sprintf("; excluded %d GHES-only toxic rule(s)", scan.ExcludedGHESRules)
		}
		fmt.Fprintln(w, line)
		fmt.Fprintln(w)
	}

	SortResults(scan.Installations)
	shown := 0
	for _, result := range scan.Installations {
		if !ew.shouldShow(result) {
			continue
		}
		if shown > 0 {
			fmt.Fprintln(w)
		}
		writeExplainInstallation(w, result)
		shown++
	}
	if shown == 0 {
		fmt.Fprintln(w, "No installations matched explain filters (CRITICAL/HIGH; use --explain-all for PASS/WARN).")
	}
	return nil
}

func (ew ExplainWriter) shouldShow(result eval.AppAuditResult) bool {
	if ew.ExplainAll {
		return true
	}
	return result.RiskLevel == "CRITICAL" || result.RiskLevel == "HIGH"
}

func writeExplainInstallation(w io.Writer, result eval.AppAuditResult) {
	fmt.Fprintf(w, "== %s — %s ==\n", appLabel(result), result.RiskLevel)
	fmt.Fprintf(w, "Repository access: %s (%d write scopes)\n", result.RepoSelection, result.WriteScopeCount)
	if result.HTMLURL != "" {
		fmt.Fprintf(w, "Installation: %s\n", result.HTMLURL)
	}

	if len(result.ToxicMatches) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Toxic combinations (all required grants matched):")
		for _, match := range result.ToxicMatches {
			label := match.Technique
			if match.ID != "" {
				label = fmt.Sprintf("%s [%s]", match.Technique, match.ID)
			}
			fmt.Fprintf(w, "  [%s] %s\n", match.Blast, label)
			if len(match.MatchedGrants) > 0 {
				fmt.Fprintln(w, "    Matched grants:")
				for _, grant := range match.MatchedGrants {
					fmt.Fprintf(w, "      • %s: %s\n", grant.APIKey, grant.Access)
				}
			}
			if match.ExploitPath != "" {
				fmt.Fprintln(w, "    What this enables:")
				writeWrapped(w, "      ", match.ExploitPath)
			}
		}
	}

	if len(result.ControlPlaneFindings) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Structural findings:")
		for _, finding := range result.ControlPlaneFindings {
			fmt.Fprintf(w, "  [%s] %s\n", finding.Risk, finding.Message)
			if finding.Rationale != "" {
				writeWrapped(w, "    ", finding.Rationale)
			}
		}
	}

	if len(result.NearMisses) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Near misses (one grant away; not counted in risk level):")
		seen := map[string]struct{}{}
		for _, near := range result.NearMisses {
			key := near.Technique + "\x00" + near.MissingGrant
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			fmt.Fprintf(w, "  • %s — missing %s\n", near.Technique, near.MissingGrant)
			if near.ExploitPath != "" {
				writeWrapped(w, "    Would enable: ", near.ExploitPath)
			}
		}
	}

	if len(result.NotableGrants) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "Notable standalone permissions (catalog severity; no toxic combo matched):")
		for _, grant := range result.NotableGrants {
			fmt.Fprintf(w, "  [%s] %s: %s\n", grant.Severity, grant.APIKey, grant.Access)
			if grant.SecurityNotes != "" {
				writeWrapped(w, "    ", grant.SecurityNotes)
			}
		}
	}

	if len(result.ToxicMatches) == 0 && len(result.ControlPlaneFindings) == 0 && len(result.NotableGrants) == 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "No toxic combinations or structural findings. Review near misses or granted permissions if this App still looks over-privileged.")
	}

	if len(result.GHESScopes) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintf(w, "GHES-only scopes granted: %s\n", strings.Join(result.GHESScopes, ", "))
	}
}

func writeWrapped(w io.Writer, prefix, text string) {
	text = strings.Join(strings.Fields(text), " ")
	const width = 88
	line := prefix
	for _, word := range strings.Fields(text) {
		if len(line)+len(word)+1 > width && line != prefix {
			fmt.Fprintln(w, line)
			line = prefix + word
		} else if line == prefix {
			line += word
		} else {
			line += " " + word
		}
	}
	if line != "" {
		fmt.Fprintln(w, line)
	}
}
