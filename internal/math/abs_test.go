package math

import "testing"

func TestAbs(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		in, out int
	}{
		{"positive", 1, 1},
		{"zero", 0, 0},
		{"negative zero", -0, 0},
		{"negative", -1, 1},
	} {
		if got, want := Abs(tt.in), tt.out; got != want {
			t.Errorf("%s: Abs() = %d, want %d", tt.desc, got, want)
		}
	}
}
