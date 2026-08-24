package output

// Options configures org scan output behavior.
type Options struct {
	// Explain emits narrative rationale (table/markdown) or enriches JSON with
	// notable grants when combined with ExplainAll.
	Explain bool
	// ExplainAll includes PASS/WARN installations in explain output. By default
	// explain mode shows CRITICAL and HIGH only.
	ExplainAll bool
}
