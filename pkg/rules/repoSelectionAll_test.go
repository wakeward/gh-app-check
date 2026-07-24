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
