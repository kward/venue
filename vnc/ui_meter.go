package vnc

import "image"

type Meter struct {
	pos      image.Point // Position of UI element.
	size     int         // Meter size (small..large).
	isStereo bool        // True if the meter is stereo.
}

// Verify that the Widget interface is honored.
var _ Widget = new(Meter)

func (w *Meter) Read(v *VNC) (interface{}, error)     { return nil, nil }
func (w *Meter) Update(v *VNC, val interface{}) error { return nil }
func (w *Meter) Press(v *VNC) error {
	return v.MouseLeftClick(w.clickOffset())
}

const (
	smallVMeter  = iota // Channel (13x50 px)
	mediumHMeter        // Comp/Lim or Exp/Gate ()
	largeVMeter         // Input ()
)

// IsMono returns true if this a mono meter.
func (w *Meter) IsMono() bool { return !w.isStereo }

// IsMono returns true if this a stereo meter.
func (w *Meter) IsStereo() bool { return w.isStereo }

func (w *Meter) clickOffset() image.Point {
	switch w.size {
	case smallVMeter:
		return w.pos.Add(image.Point{7, 25})
	case mediumHMeter:
		return w.pos.Add(image.Point{0, 0})
	case largeVMeter:
		return w.pos.Add(image.Point{0, 0})
	}
	return w.pos
}
