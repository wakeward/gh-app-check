package rules

// WriteScopeCount returns how many permission grants are at write level.
func WriteScopeCount(inst Installation) int {
	count := 0
	for _, level := range inst.Permissions {
		if level == "write" {
			count++
		}
	}
	return count
}

// HasAnyWriteScope reports whether the installation has at least one write grant.
func HasAnyWriteScope(inst Installation) bool {
	return WriteScopeCount(inst) > 0
}
