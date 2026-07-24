package cmd

import "testing"

func TestValidateFormat(t *testing.T) {
	cases := []struct {
		format  string
		wantErr bool
	}{
		{"table", false},
		{"json", false},
		{"markdown", false},
		{"yaml", true},
		{"", true},
	}

	original := format
	defer func() { format = original }()

	for _, tc := range cases {
		t.Run(tc.format, func(t *testing.T) {
			format = tc.format
			err := validateFormat(rootCmd, nil)
			if (err != nil) != tc.wantErr {
				t.Errorf("validateFormat() with format=%q, err = %v, wantErr = %v", tc.format, err, tc.wantErr)
			}
		})
	}
}
