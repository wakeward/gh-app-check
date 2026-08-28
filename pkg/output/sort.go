package output

import (
	"slices"
	"strings"

	"github.com/wakeward/gh-app-check/pkg/eval"
)

// SortResults orders audit results by descending risk, then app slug.
func SortResults(results []eval.AppAuditResult) {
	slices.SortFunc(results, func(a, b eval.AppAuditResult) int {
		if ra, rb := eval.RiskRank(a.RiskLevel), eval.RiskRank(b.RiskLevel); ra != rb {
			return rb - ra
		}
		return strings.Compare(a.AppSlug, b.AppSlug)
	})
}
