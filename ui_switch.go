package venue

import "image"

// Switch size examples:
// - Large: Channel solo or mute
// - Medium: 48V, Pad, Guess
// - Small: Encoder ON
// - Tiny: Channel solo or mute on bank
const (
	tinySwitch = iota
	smallSwitch
	mediumSwitch
	largeSwitch
)

// Switch is a UIElement representing a switch.
type Switch struct {
	rect   image.Rectangle // Rectangle representing the UI element.
	toggle bool            // Toggle switch? true = toggle, false = hold-to-enable
	// state  bool            // Current state; true = on (pressed)
	// def    bool            // Default value; true = on
	size   int // tiny..large
	center image.Point
}

func newSwitch(x, y int, size int, toggle, def bool) *Switch {
	var (
		dx, dy int
		center image.Point
	)
	switch size {
	case tinySwitch:
		dx, dy = 13, 13
		center = image.Point{x + 7, y + 7}
	}
	return &Switch{
		rect:   image.Rect(x, y, x+dx, y+dy),
		toggle: toggle,
		// state:  def,
		// def:    def,
		size:   size,
		center: center,
	}
}

func (e *Switch) Read(v *Venue) error {
	return nil
}

func (e *Switch) Select(v *Venue) {
	// Move mouse pointer to switch.
	v.MouseLeftClick(e.center)
}

func (e *Switch) Update(v *Venue) error {
	return nil
}
