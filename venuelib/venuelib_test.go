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
	for _, tt := range []struct {
		code codes.Code
		desc string
	}{
		// Test that Errorf returns nil for OK code
		{codes.OK, "test message"},
		// Test that Errorf creates proper error for non-OK code
		{codes.NotFound, "file not found"},
		{codes.Internal, "internal error"},
		{codes.InvalidArgument, "invalid argument"},
		{codes.DeadlineExceeded, "deadline exceeded"},
	} {
		err := Errorf(tt.code, "%s", tt.desc)
		if tt.code == codes.OK {
			if err != nil {
				t.Errorf("Errorf(%v, %q) should return nil, got %v", tt.code, tt.desc, err)
			}
		} else {
			if err == nil {
				t.Errorf("Errorf(%v, %q) should return an error, got nil", tt.code, tt.desc)
			}

			// Test that Code function extracts the correct error code
			code := Code(err)
			if code != tt.code {
				t.Errorf("Code(Errorf(%v, %q)) should return %v, got %v", tt.code, tt.desc, tt.code, code)
			}

			// Test that ErrorDesc function extracts the correct error description
			desc := ErrorDesc(err)
			if desc != tt.desc {
				t.Errorf("ErrorDesc(Errorf(%v, %q)) should return %q, got %q", tt.code, tt.desc, tt.desc, desc)
			}
		}
	}
}

func ExampleFnName() {
	if glog.V(3) {
		glog.Info(FnName())
	}
}
