package rules

// AdministrationWrite flags installations granted write access to
// organization/repository administration ("Flag if permissions.administration
// == write", Design Spec §4.A).
func AdministrationWrite(inst Installation) bool {
	return inst.Permissions["administration"] == "write"
}
