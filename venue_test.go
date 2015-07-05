package venue

import "testing"

// "image"
// "testing"

// func TestOverlay(t *testing.T) {
// 	tests := []struct {
// 		a, b, c *image.RGBA
// 	}{
// 		{ // 0, 0 is painted; overlay at 1, 1.
// 			&image.RGBA{
// 				[]uint8{255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
// 				8, image.Rectangle{image.Point{0, 0}, image.Point{2, 2}}},
// 			&image.RGBA{
// 				[]uint8{255, 255, 255, 255},
// 				4, image.Rectangle{image.Point{1, 1}, image.Point{2, 2}}},
// 			&image.RGBA{
// 				[]uint8{255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255},
// 				8, image.Rectangle{image.Point{0, 0}, image.Point{2, 2}}},
// 		},
// 		{ // 0, 0 and 1, 1 are painted; overlay new 1, 1.
// 			&image.RGBA{
// 				[]uint8{255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 255, 255, 255, 255},
// 				8, image.Rectangle{image.Point{0, 0}, image.Point{2, 2}}},
// 			&image.RGBA{
// 				[]uint8{127, 127, 127, 255},
// 				4, image.Rectangle{image.Point{1, 1}, image.Point{2, 2}}},
// 			&image.RGBA{
// 				[]uint8{255, 255, 255, 255, 0, 0, 0, 0, 0, 0, 0, 0, 127, 127, 127, 255},
// 				8, image.Rectangle{image.Point{0, 0}, image.Point{2, 2}}},
// 		},
// 	}

// 	for _, tt := range tests {
// 		Overlay(tt.a, tt.b)
// 		if got, want := tt.a.Pix, tt.c.Pix; !EqualSlices(got, want) {
// 			t.Errorf("Overlay() failed; Pix got = %v, want = %v", got, want)
// 		}
// 	}
// }

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

func equalSlices(x, y interface{}) bool {
	// Special cases.
	switch {
	case x == nil && y == nil:
		return true
	case x == nil || y == nil:
		return false
	}

	switch x.(type) {
	case []uint8:
		if len(x.([]uint8)) != len(y.([]uint8)) {
			return false
		}
		for i, v := range x.([]uint8) {
			if v != y.([]uint8)[i] {
				return false
			}
		}
	}
	return true
}
