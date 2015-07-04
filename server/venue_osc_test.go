package main

import (
	"testing"
)

func TestMultiPosition(t *testing.T) {
	tests := []struct {
		x, y, dx, dy, bank int
		pos                int
	}{
		// 1x1 control
		{1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 2, 2},
		// 2x2 control
		{1, 1, 2, 2, 1, 1},
		{2, 1, 2, 2, 1, 2},
		{1, 2, 2, 2, 1, 3},
		{2, 2, 2, 2, 1, 4},
		{2, 2, 2, 2, 2, 8},
		// 3x2 control
		{1, 1, 3, 2, 1, 1},
		{3, 1, 3, 2, 1, 3},
		{1, 2, 3, 2, 1, 4},
		{3, 2, 3, 2, 1, 6},
		{3, 2, 3, 2, 2, 12},
	}

	for tnum, tt := range tests {
		if got, want := multiPosition(tt.x, tt.y, tt.dx, tt.dy, tt.bank), tt.pos; got != want {
			t.Errorf("multiPosition(%v): got = %v, want = %v", tnum, got, want)
		}
	}
}

func TestMultiRotate(t *testing.T) {
	// Note: dx is never used, but it is given to make it easier for a human to
	// understand what's going on.
	tests := []struct {
		x, y, dx, dy int
		xx, yy       int
	}{
		// 1x1 control
		{1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1},
		// 2x2 control
		{1, 1, 2, 2, 1, 2},
		{2, 1, 2, 2, 1, 1},
		{1, 2, 2, 2, 2, 2},
		{2, 2, 2, 2, 2, 1},
		// 3x2 control
		{1, 1, 3, 2, 1, 2},
		{2, 1, 3, 2, 1, 1},
		{1, 3, 3, 2, 3, 2},
		{2, 3, 3, 2, 3, 1},
	}

	for tnum, tt := range tests {
		xx, yy := multiRotate(tt.x, tt.y, tt.dy)
		if xx != tt.xx || yy != tt.yy {
			t.Errorf("multiRotate(%v): got = x:%v y:%v, want = x:%v y:%v", tnum, xx, yy, tt.xx, tt.yy)
		}
	}
}

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
