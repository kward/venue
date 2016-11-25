package venue

import "image"

const (
	smallVMeter  = iota // Channel (13x50 px)
	mediumHMeter        // Comp/Lim or Exp/Gate ()
	largeVMeter         // Input ()
	monoMeter    = false
	stereoMeter  = true
)

type Meter struct {
	pos    image.Point // Position of UI element.
	size   int         // small..large
	stereo bool        // Stereo? false = mono, true = stereo.
}

func (e *Meter) Read(v *Venue) error   { return nil }
func (e *Meter) Select(v *Venue)       { v.vnc.MouseLeftClick(e.clickOffset()) }
func (e *Meter) Set(v *Venue, val int) {}
func (e *Meter) Update(v *Venue) error { return nil }

func (e *Meter) clickOffset() image.Point {
	switch e.size {
	case smallVMeter:
		return e.pos.Add(image.Point{7, 25})
	case mediumHMeter:
		return e.pos.Add(image.Point{0, 0})
	case largeVMeter:
		return e.pos.Add(image.Point{0, 0})
	}
	return e.pos
}
