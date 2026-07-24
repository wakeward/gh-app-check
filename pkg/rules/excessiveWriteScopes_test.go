package rules

import "testing"

func TestExcessiveWriteScopes(t *testing.T) {
	writePerms := func(n int) map[string]string {
		perms := make(map[string]string, n)
		for i := 0; i < n; i++ {
			perms[string(rune('a'+i))] = "write"
		}
		return perms
	}

	cases := []struct {
		name string
		inst Installation
		want bool
	}{
		{"exactly at threshold not flagged", Installation{Permissions: writePerms(5)}, false},
		{"one over threshold flagged", Installation{Permissions: writePerms(6)}, true},
		{"mixed read/write below threshold not flagged", Installation{Permissions: map[string]string{
			"metadata":       "read",
			"contents":       "write",
			"issues":         "write",
			"pull_requests":  "read",
			"administration": "read",
		}}, false},
		{"no permissions not flagged", Installation{}, false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ExcessiveWriteScopes(tc.inst); got != tc.want {
				t.Errorf("ExcessiveWriteScopes(%+v) = %v, want %v", tc.inst, got, tc.want)
			}
		})
	}
}
