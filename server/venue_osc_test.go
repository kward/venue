package main

import (
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
		val, abs int
	}{
		{1, 1},
		{0, 0},
		{-1, 1},
	}

	for _, tt := range tests {
		if got, want := abs(tt.val), tt.abs; got != want {
			t.Errorf("abs(): got = %v, want = %v", got, want)
		}
	}
}

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

func TestToInt(t *testing.T) {
	tests := []struct {
		s string
		i int
	}{
		{"1", 1},
		{"0", 0},
		{"-1", -1},
		{"foo", -1},
	}

	for _, tt := range tests {
		if got, want := toInt(tt.s), tt.i; got != want {
			t.Errorf("toInt() failed; got = %v, want = %v", got, want)
		}
	}
}
