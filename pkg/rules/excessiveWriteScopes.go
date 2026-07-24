package rules

// MaxWriteScopes is the threshold above which an installation is considered
// "god-mode" ("Calculate total write scopes; flag if > 5", Design Spec §4.A).
const MaxWriteScopes = 5

// ExcessiveWriteScopes flags installations with more than MaxWriteScopes
// write-level permission scopes.
func ExcessiveWriteScopes(inst Installation) bool {
	writeCount := 0
	for _, level := range inst.Permissions {
		if level == "write" {
			writeCount++
		}
	}
	return writeCount > MaxWriteScopes
}
