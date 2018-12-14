package decimal

import "testing"

func TestParse(t *testing.T) {
	testCases := []struct {
		input  string
		output string
	}{
		{input: "1.23", output: "1.23"},
		{input: "-12.34", output: "-12.34"},
		{input: "0.12", output: "0.12"},
		{input: ".12", output: "0.12"},
		{input: ".012", output: "0.012"},
	}
	for _, tc := range testCases {
		d, err := New(tc.input)
		if err != nil {
			t.Error(err)
			continue
		}
		if d.String() != tc.output {
			t.Errorf("%s: got %q, want %q", tc.input, d.String(), tc.output)
		}
	}
}
