package rules

// RepoSelectionAll flags installations with access to all repositories in
// the organization ("Flag if repository_selection == all", Design Spec §4.A).
func RepoSelectionAll(inst Installation) bool {
	return inst.RepositorySelection == "all"
}
