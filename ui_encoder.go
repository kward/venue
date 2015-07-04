package venue

import (
	"image"

	vnc "github.com/kward/go-vnc"
)

type EncoderWindow int

const (
	EncoderTL EncoderWindow = iota // Top left
	EncoderML                      // Middle left
	EncoderBL                      // Bottom left
	EncoderBC                      // Bottom center
	EncoderTR                      // Top right
	EncoderMR                      // Middle right
	EncoderBR                      // Bottom right
)

type Encoder struct {
	center   image.Point
	window   EncoderWindow // Position of value window
	hasOnOff bool          // Has an on/off switch
}

func (e *Encoder) Read(v *Venue) error {
	// TODO(kward): select Inputs page

	// Give the window focus.
	p := e.clickOffset()
	v.MouseLeftClick(p)

	// Cut the selected text.

	return nil
}

func (e *Encoder) Update(v *Venue) error {
	// TODO(kward): select Inputs page

	// Give window focus.
	v.MouseLeftClick(e.clickOffset())

	// Move mouse pointer center of Encoder.
	v.MouseMove(e.center)

	// Update
	v.KeyPress(vnc.Key3)
	v.KeyPress(vnc.Key4)
	v.KeyPress(vnc.KeyReturn)
	for i := 0; i < 5; i++ {
		e.Increment(v)
	}

	return nil
}

func (e *Encoder) Adjust(v *Venue, c int) {
	v.MouseDrag(e.center, image.Point{0, -c})
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

	if e.window == EncoderTL || e.window == EncoderML || e.window == EncoderBL { // Left
		dx = -38
	} else if e.window == EncoderTR || e.window == EncoderMR || e.window == EncoderBR { // Right
		dx = 38
	} else { // Middle
		dx = 0
	}

	if e.window == EncoderTL || e.window == EncoderTR { // Top
		dy = -8
	} else if e.window == EncoderML || e.window == EncoderMR { // Middle
		dy = 0
	} else if e.window == EncoderBL || e.window == EncoderBR { // Bottom
		dy = 8
	} else { // Center
		dy = 28
	}

	return e.center.Add(image.Point{dx, dy})
}
