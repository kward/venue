package venue

import (
	"image"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/venue/math"
	"github.com/kward/venue/venuelib"
	"github.com/kward/venue/vnc"
)

type encoderEnum int

const (
	encoderTL encoderEnum = iota // Top left
	encoderML                    // Middle left
	encoderBL                    // Bottom left
	encoderBC                    // Bottom center
	encoderTR                    // Top right
	encoderMR                    // Middle right
	encoderBR                    // Bottom right
)

type Encoder struct {
	center   image.Point
	window   encoderEnum // Position of value window
	hasOnOff bool        // Has an on/off switch
}

// Verify that the Widget interface is honored.
var _ Widget = new(Encoder)

func (w *Encoder) Read(v *vnc.VNC) (interface{}, error) { return nil, nil }

func (w *Encoder) Update(v *vnc.VNC, val interface{}) error {
	w.Press(v)
	for _, key := range vnc.IntToKeys(val.(int)) {
		v.KeyPress(key)
	}
	v.KeyPress(vnclib.KeyReturn)
	return nil
}

func (w *Encoder) Press(v *vnc.VNC) error {
	return v.MouseLeftClick(w.clickOffset())
}

// Adjust the value of an encoder with cursor keys.
func (w *Encoder) Adjust(v *vnc.VNC, val int) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Adjusting encoder by %d steps.", val)
	}
	if val == 0 {
		return nil
	}

	if err := w.Press(v); err != nil {
		return err
	}
	key := vnclib.KeyUp
	if val < 0 {
		key = vnclib.KeyDown
	}
	amount := math.Abs(val)
	for i := 0; i < amount; i++ {
		if err := v.KeyPress(key); err != nil {
			return err
		}
	}
	return v.KeyPress(vnclib.KeyReturn)
}

// Increment the value of an encoder.
func (w *Encoder) Increment(v *vnc.VNC) error { return w.Adjust(v, 1) }

// Decrement the value of an encoder.
func (w *Encoder) Decrement(v *vnc.VNC) error { return w.Adjust(v, -1) }

func (w *Encoder) clickOffset() image.Point {
	var dx, dy int

	// Horizontal position.
	switch w.window {
	case encoderTL, encoderML, encoderBL: // Left
		dx = -38
	case encoderTR, encoderMR, encoderBR: // Right
		dx = 38
	default: // Middle
		dx = 0
	}

	// Vertical position.
	switch w.window {
	case encoderTL, encoderTR: // Top
		dy = -8
	case encoderML, encoderMR: // Middle
		dy = 0
	case encoderBL, encoderBR: // Bottom
		dy = 8
	default: // Center
		dy = 28
	}

	return w.center.Add(image.Point{dx, dy})
}
