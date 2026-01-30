package vnc

// func TestOverlay(t *testing.T) {
// 	for _, tt := range []struct {
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
// 	} {
// 		Overlay(tt.a, tt.b)
// 		if got, want := tt.a.Pix, tt.c.Pix; !EqualSlices(got, want) {
// 			t.Errorf("Overlay() failed; Pix got = %v, want = %v", got, want)
// 		}
// 	}
// }
