package rules

import "testing"

func TestRepoSelectionAll(t *testing.T) {
	cases := []struct {
		name string
		inst Installation
		want bool
	}{
		{"all repos flagged", Installation{RepositorySelection: "all"}, true},
		{"selected repos not flagged", Installation{RepositorySelection: "selected"}, false},
		{"empty selection not flagged", Installation{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := RepoSelectionAll(tc.inst); got != tc.want {
				t.Errorf("RepoSelectionAll(%+v) = %v, want %v", tc.inst, got, tc.want)
			}
		})
	}
}

func TestRepoSelectionAllWithWrite(t *testing.T) {
	if !RepoSelectionAllWithWrite(Installation{
		RepositorySelection: "all",
		Permissions:         map[string]string{"contents": "write"},
	}) {
		t.Fatal("expected all-repos with write to match")
	}
	if RepoSelectionAllWithWrite(Installation{
		RepositorySelection: "all",
		Permissions:         map[string]string{"contents": "read"},
	}) {
		t.Fatal("read-only all-repos should not match WithWrite")
	}
}

func TestRepoSelectionAllReadOnly(t *testing.T) {
	if !RepoSelectionAllReadOnly(Installation{
		RepositorySelection: "all",
		Permissions:         map[string]string{"contents": "read", "metadata": "read"},
	}) {
		t.Fatal("expected read-only all-repos to match")
	}
	if RepoSelectionAllReadOnly(Installation{
		RepositorySelection: "all",
		Permissions:         map[string]string{"contents": "write"},
	}) {
		t.Fatal("all-repos with write should not match ReadOnly")
	}
}

func TestWriteScopeCount(t *testing.T) {
	if got := WriteScopeCount(Installation{Permissions: map[string]string{
		"contents": "write",
		"metadata": "read",
	}}); got != 1 {
		t.Fatalf("WriteScopeCount = %d, want 1", got)
	}
}
