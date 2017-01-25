package vnc

import (
	"fmt"
	"reflect"
	"testing"

	vnclib "github.com/kward/go-vnc"
)

func TestIntToKeys(t *testing.T) {
	for _, tt := range []struct {
		desc string
		v    int      // Value.
		keys []uint32 // Necessary keys.
	}{
		{"zero", 0, []uint32{vnclib.Key0}},
		{"negative", -123, []uint32{vnclib.KeyMinus, vnclib.Key1, vnclib.Key2, vnclib.Key3}},
		{"positive", 456, []uint32{vnclib.Key4, vnclib.Key5, vnclib.Key6}},
	} {
		if got, want := IntToKeys(tt.v), tt.keys; !equalSlices(got, want) {
			t.Errorf("%s: intToKeys() = %v, want = %v", tt.desc, got, want)
		}
	}
}

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
	case []uint32:
		if len(x.([]uint32)) != len(y.([]uint32)) {
			return false
		}
		for i, v := range x.([]uint32) {
			if v != y.([]uint32)[i] {
				return false
			}
		}
	default:
		panic(fmt.Sprintf("unrecognized type %v", reflect.TypeOf(x)))
	}
	return true
}
