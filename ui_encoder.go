package venue

import (
	"fmt"
	"image"

	vnclib "github.com/kward/go-vnc"
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
func (e *Encoder) Select(v *Venue)     { v.vnc.MouseLeftClick(e.clickOffset()) }

func (e *Encoder) Set(v *Venue, val int) {
	e.Select(v)
	for _, key := range intToKeys(val) {
		v.vnc.KeyPress(key)
	}
	v.vnc.KeyPress(vnclib.KeyReturn)
}

func (e *Encoder) Update(v *Venue) error { return nil }

func (e *Encoder) Adjust(v *Venue, c int) {
	v.vnc.MouseLeftClick(e.clickOffset())
	for i := 0; i < abs(c); i++ {
		if c > 0 {
			v.vnc.KeyPress(vnclib.KeyUp)
		} else {
			v.vnc.KeyPress(vnclib.KeyDown)
		}
	}
	v.vnc.KeyPress(vnclib.KeyReturn)
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
		'-': vnclib.KeyMinus,
		'0': vnclib.Key0,
		'1': vnclib.Key1,
		'2': vnclib.Key2,
		'3': vnclib.Key3,
		'4': vnclib.Key4,
		'5': vnclib.Key5,
		'6': vnclib.Key6,
		'7': vnclib.Key7,
		'8': vnclib.Key8,
		'9': vnclib.Key9,
	}
	k := []uint32{}
	s := fmt.Sprintf("%d", v)
	for _, c := range s {
		k = append(k, keys[c])
	}
	return k
}
