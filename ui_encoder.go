package venue

import (
	"fmt"
	"image"

	vnc "github.com/kward/go-vnc"
)

const (
	encoderTL = iota // Top left
	encoderML        // Middle left
	encoderBL        // Bottom left
	encoderBC        // Bottom center
	encoderTR        // Top right
	encoderMR        // Middle right
	encoderBR        // Bottom right
)

type Encoder struct {
	center   image.Point
	window   int  // Position of value window
	hasOnOff bool // Has an on/off switch
}

func (e *Encoder) Read(v *Venue) error { return nil }
func (e *Encoder) Select(v *Venue)     { v.MouseLeftClick(e.clickOffset()) }

func (e *Encoder) Set(v *Venue, val int) {
	e.Select(v)
	for _, key := range intToKeys(val) {
		v.KeyPress(key)
	}
	v.KeyPress(vnc.KeyReturn)
}

func (e *Encoder) Update(v *Venue) error { return nil }

func (e *Encoder) Adjust(v *Venue, c int) {
	v.MouseLeftClick(e.clickOffset())
	for i := 0; i < abs(c); i++ {
		if c > 0 {
			v.KeyPress(vnc.KeyUp)
		} else {
			v.KeyPress(vnc.KeyDown)
		}
	}
	v.KeyPress(vnc.KeyReturn)
}

func (e *Encoder) Increment(v *Venue) {
	e.Adjust(v, 1)
}

func (e *Encoder) Decrement(v *Venue) {
	e.Adjust(v, -1)
}

func (e *Encoder) Refresh(v *Venue) {}

func (e *Encoder) clickOffset() image.Point {
	var dx, dy int

	if e.window == encoderTL || e.window == encoderML || e.window == encoderBL { // Left
		dx = -38
	} else if e.window == encoderTR || e.window == encoderMR || e.window == encoderBR { // Right
		dx = 38
	} else { // Middle
		dx = 0
	}

	if e.window == encoderTL || e.window == encoderTR { // Top
		dy = -8
	} else if e.window == encoderML || e.window == encoderMR { // Middle
		dy = 0
	} else if e.window == encoderBL || e.window == encoderBR { // Bottom
		dy = 8
	} else { // Center
		dy = 28
	}

	return e.center.Add(image.Point{dx, dy})
}

func intToKeys(v int) []uint32 {
	keys := map[rune]uint32{
		'-': vnc.KeyMinus,
		'0': vnc.Key0,
		'1': vnc.Key1,
		'2': vnc.Key2,
		'3': vnc.Key3,
		'4': vnc.Key4,
		'5': vnc.Key5,
		'6': vnc.Key6,
		'7': vnc.Key7,
		'8': vnc.Key8,
		'9': vnc.Key9,
	}
	k := []uint32{}
	s := fmt.Sprintf("%d", v)
	for _, c := range s {
		k = append(k, keys[c])
	}
	return k
}
