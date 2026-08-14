package rules

// RepoSelectionAll flags installations with access to all repositories in
// the organization ("Flag if repository_selection == all", Design Spec §4.A).
func RepoSelectionAll(inst Installation) bool {
	return inst.RepositorySelection == "all"
}

// RepoSelectionAllWithWrite flags all-repository access combined with at
// least one write-level permission grant.
func RepoSelectionAllWithWrite(inst Installation) bool {
	return inst.RepositorySelection == "all" && HasAnyWriteScope(inst)
}

// RepoSelectionAllReadOnly flags all-repository access when every grant is
// read (or the permission map has no write scopes).
func RepoSelectionAllReadOnly(inst Installation) bool {
	return inst.RepositorySelection == "all" && !HasAnyWriteScope(inst)
}
