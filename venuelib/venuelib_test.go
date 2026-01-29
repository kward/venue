package venuelib

import (
	"testing"

	"github.com/golang/glog"
	"github.com/kward/venue/codes"
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

func TestErrorf(t *testing.T) {
	// Test that Errorf returns nil for OK code
	err := Errorf(codes.OK, "test message")
	if err != nil {
		t.Errorf("Errorf(codes.OK, ...) should return nil, got %v", err)
	}

	// Test that Errorf creates proper error for non-OK code
	err = Errorf(codes.NotFound, "file not found: %s", "test.txt")
	if err == nil {
		t.Error("Errorf(codes.NotFound, ...) should return an error")
	}

	// Test that Code function extracts the correct error code
	code := Code(err)
	if code != codes.NotFound {
		t.Errorf("Code(err) should return codes.NotFound, got %v", code)
	}

	// Test that ErrorDesc function extracts the correct error description
	desc := ErrorDesc(err)
	if desc != "file not found: test.txt" {
		t.Errorf("ErrorDesc(err) should return \"file not found: test.txt\", got %q", desc)
	}
}

func ExampleFnName() {
	if glog.V(3) {
		glog.Info(FnName())
	}
}
