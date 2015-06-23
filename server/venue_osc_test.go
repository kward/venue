package main

import (
	"testing"
)

func TestCar(t *testing.T) {
	tests := []struct {
		addr, first string
	}{
		{"/0.0/foo/bar", "0.0"},
		{"/foo/bar", "foo"},
		{"/bar", "bar"},
		{"", ""},
	}

	for _, tt := range tests {
		if got, want := car(tt.addr), tt.first; got != want {
			t.Errorf("car() failed; got = %v, want = %v", got, want)
		}
	}
}
func TestCarInt(t *testing.T) {
	tests := []struct {
		addr, first int
	}{
		{"/2/4", 2},
		{"/4", 4},
		{"", -1},
	}

	for _, tt := range tests {
		if got, want := carInt(tt.addr), tt.first; got != want {
			t.Errorf("car() failed; got = %v, want = %v", got, want)
		}
	}
}
func TestCdr(t *testing.T) {
	tests := []struct {
		addr, first string
	}{
		{"/0.0/foo/bar", "/foo/bar"},
		{"/foo/bar", "/bar"},
		{"/bar", ""},
		{"", ""},
	}

	for _, tt := range tests {
		if got, want := cdr(tt.addr), tt.first; got != want {
			t.Errorf("cdr() failed; got = %v, want = %v", got, want)
		}
	}
}
