package venuelib

import "testing"

func TestToInt(t *testing.T) {
	for _, tt := range []struct {
		in   string
		want int
	}{
		// valid
		{"123", 123},
		{"0", 0},
		{"-123", -123},
		// invalid
		{"abc", 0},
	} {
		if got, want := ToInt(tt.in), tt.want; got != want {
			t.Errorf("ToInt(%s) = %d, want = %d", tt.in, got, want)
		}
	}
}
