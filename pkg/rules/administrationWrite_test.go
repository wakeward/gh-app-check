package rules

import "testing"

func TestAdministrationWrite(t *testing.T) {
	cases := []struct {
		name string
		inst Installation
		want bool
	}{
		{
			name: "administration write flagged",
			inst: Installation{Permissions: map[string]string{"administration": "write"}},
			want: true,
		},
		{
			name: "administration read not flagged",
			inst: Installation{Permissions: map[string]string{"administration": "read"}},
			want: false,
		},
		{
			name: "no permissions not flagged",
			inst: Installation{},
			want: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := AdministrationWrite(tc.inst); got != tc.want {
				t.Errorf("AdministrationWrite(%+v) = %v, want %v", tc.inst, got, tc.want)
			}
		})
	}
}
