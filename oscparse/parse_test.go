package oscparse

import "testing"

type parseTest struct {
	name string
	addr string
	val  []interface{}
	pkt  *Packet
	ok   bool
}

func TestIsTablet(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		layout  string
		isHoriz bool
	}{
		{"tablet/horiz", "th", true},
		{"tablet/vert", "tv", true},
		{"phone/horiz", "ph", false},
		{"phone/vert", "pv", false},
		{"unknown", "unknown", false},
	} {
		req := request{layout: tt.layout}
		if got, want := req.isTablet(), tt.isHoriz; got != want {
			t.Errorf("%s: isTablet() = %v, want %v", tt.desc, got, want)
			continue
		}
	}
}

func TestIsHorizontal(t *testing.T) {
	for _, tt := range []struct {
		desc    string
		layout  string
		isHoriz bool
	}{
		{"tablet/horiz", "th", true},
		{"phone/horiz", "ph", true},
		{"tablet/vert", "tv", false},
		{"phone/vert", "pv", false},
		{"unknown", "unknown", false},
	} {
		req := request{layout: tt.layout}
		if got, want := req.isHorizontal(), tt.isHoriz; got != want {
			t.Errorf("%s: isHorizontal() = %v, want %v", tt.desc, got, want)
			continue
		}
	}
}

var multiTests = []struct {
	desc       string
	dx, dy     int // dy is unused.
	x, y       int
	xx, yy     int
	vPos, hPos int
}{
	// 1x1 control
	{"1x1", 1, 1, 1, 1, 1, 1, 1, 1},
	// 2x2 control
	{"2x2-1,1", 2, 2, 1, 1, 1, 2, 1, 3},
	{"2x2-2,1", 2, 2, 2, 1, 1, 1, 2, 1},
	{"2x2-1,2", 2, 2, 1, 2, 2, 2, 3, 4},
	{"2x2-2,2", 2, 2, 2, 2, 2, 1, 4, 2},
	// 3x2 control
	{"3x2-1,1", 3, 2, 1, 1, 1, 3, 1, 5},
	{"3x2-2,1", 3, 2, 2, 1, 1, 2, 2, 3},
	{"3x2-3,1", 3, 2, 3, 1, 1, 1, 3, 1},
	{"3x2-1,2", 3, 2, 1, 2, 2, 3, 4, 6},
	{"3x2-2,2", 3, 2, 2, 2, 2, 2, 5, 4},
	{"3x2-3,2", 3, 2, 3, 2, 2, 1, 6, 2},
}

func TestMultiPosition(t *testing.T) {
	for _, tt := range multiTests {
		req := request{x: tt.x, y: tt.y}

		// Vertical layout.
		req.layout = "pv"
		if got, want := req.multiPosition(tt.dx, tt.dy), tt.vPos; got != want {
			t.Errorf("%s-%s: multiPosition() = %d, want %d", tt.desc, req.layout, got, want)
		}

		// Horizontal layout.
		req.layout = "th"
		if got, want := req.multiPosition(tt.dx, tt.dy), tt.hPos; got != want {
			t.Errorf("%s-%s: multiPosition() = %d, want %d", tt.desc, req.layout, got, want)
		}
	}
}

func TestMultiRotate(t *testing.T) {
	for _, tt := range multiTests {
		req := request{x: tt.x, y: tt.y}

		// Vertical layout, expect no coordinate rotation.
		req.layout = "pv"
		xx, yy := req.multiRotate(tt.dx)
		if xx != tt.x || yy != tt.y {
			t.Errorf("%s-%s: multiRotate() = %d/%d, want %d/%d", tt.desc, req.layout, xx, yy, tt.x, tt.y)
		}

		// Horizontal layout, expect coordinate rotation.
		req.layout = "th"
		xx, yy = req.multiRotate(tt.dx)
		if xx != tt.xx || yy != tt.yy {
			t.Errorf("%s-%s: multiRotate() = %d/%d, want %d/%d", tt.desc, req.layout, xx, yy, tt.xx, tt.yy)
		}
	}
}
