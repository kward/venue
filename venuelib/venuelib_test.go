package venuelib

import (
	"testing"

	"github.com/golang/glog"
)

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
			t.Errorf("ToInt(%s) = %d, want %d", tt.in, got, want)
		}
	}
}

func TestFnName(t *testing.T) {
	if got, want := FnName(), "TestFnName()"; got != want {
		t.Errorf("FnName() = %s, want %s", got, want)
	}
}

func ExampleFnName() {
	if glog.V(3) {
		glog.Info(FnName())
	}
}
