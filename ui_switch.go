package venue

import "image"

// Switch size examples:
// - Tiny: Channel solo or mute on bank (13x13 px)
// - Small: Encoder ON (18x14 px)
// - Medium: 48V, Pad, Guess (26x16 px)
// - Large: Channel solo or mute (32x18)
const (
	tinySwitch = iota
	smallSwitch
	mediumSwitch
	largeSwitch
)

// Switch is a UIElement representing a switch.
type Switch struct {
	pos    image.Point // Position of UI element.
	size   int         // tiny..large
	toggle bool        // Toggle switch? true = toggle, false = hold-to-enable
	def    bool        // Default value; true = on
	state  bool        // Current state; true = on (pressed)
}

func newPushButton(x, y int, size int) *Switch {
	return &Switch{
		pos:  image.Point{x, y},
		size: size,
	}
}

func newToggle(x, y int, size int, def bool) *Switch {
	return &Switch{
		pos:    image.Point{x, y},
		size:   size,
		toggle: true,
		def:    def,
		state:  def,
	}
}

func (e *Switch) Read(v *Venue) error {
	return nil
}

func (e *Switch) Update(v *Venue) error {
	// Move mouse pointer to switch.
	v.MouseLeftClick(e.clickOffset())

	// Update local state.
	if e.toggle {
		e.state = !e.state
	}

	return nil
}

func (e *Switch) clickOffset() image.Point {
	switch e.size {
	case tinySwitch:
		return e.pos.Add(image.Point{7, 7})
	case smallSwitch:
		return e.pos.Add(image.Point{9, 7})
	case mediumSwitch:
		return e.pos.Add(image.Point{13, 8})
	case largeSwitch:
		return e.pos.Add(image.Point{16, 9})
	}
	return e.pos
}
