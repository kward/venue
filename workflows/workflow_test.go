package workflow

import "testing"

func TestMultiPosition(t *testing.T) {
	for tnum, tt := range []struct {
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
	} {
		if got, want := multiPosition(tt.x, tt.y, tt.dx, tt.dy, tt.bank), tt.pos; got != want {
			t.Errorf("%v: multiPosition() = %v, want = %v", tnum, got, want)
		}
	}
}

func TestMultiRotate(t *testing.T) {
	// Note: dx is never used, but it is given to make it easier for a human to
	// understand what's going on.
	for tnum, tt := range []struct {
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
	} {
		xx, yy := multiRotate(tt.x, tt.y, tt.dy)
		if xx != tt.xx || yy != tt.yy {
			t.Errorf("%v: multiRotate(): got = x:%v y:%v, want = x:%v y:%v", tnum, xx, yy, tt.xx, tt.yy)
		}
	}
}
