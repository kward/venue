package vnc

import "image"

// Switch size examples:
// - Tiny: Channel solo or mute on bank (13x13 px)
// - Small: Encoder ON (18x14 px)
// - Medium: 48V, Pad, Guess (26x16 px)
// - Large: Channel solo or mute (32x18)
const ( // Switch size.
	tinySwitch = iota
	smallSwitch
	mediumSwitch
	largeSwitch
)
const ( // Switch kind.
	toggleSwitch = iota
	pushButtonSwitch
)
const ( // Switch enabled.
	SwitchOff = false
	SwitchOn  = true
)

// Switch is a UIElement representing a switch.
type Switch struct {
	pos        image.Point // Position of UI element.
	size       int         // Switch size (tiny..large).
	kind       int         // true == toggle, false == push-button
	defEnabled bool        // Default switch state.
	isEnabled  bool        // True if the switch state is enabled.
}

// Verify that the Widget interface is honored.
var _ Widget = new(Switch)

func (w *Switch) Read(v *VNC) (interface{}, error) { return nil, nil }

func (w *Switch) Update(v *VNC, val interface{}) error {
	if w.IsPushButton() {
		// It doesn't make sense to update a push-button.
		return nil
	}

	val, err := w.Read(v)
	if err != nil {
		return err
	}
	if w.isEnabled != val.(bool) {
		return w.Press(v)
	}
	return nil
}

func (w *Switch) Press(v *VNC) error {
	return v.MouseLeftClick(w.clickOffset())
}

// IsEnabled returns true if the switch is enabled.
func (w *Switch) IsEnabled() bool { return w.isEnabled }

// NewPushButton returns a new push-button switch.
func NewPushButton(x, y int, size int) *Switch {
	return &Switch{
		pos:  image.Point{x, y},
		size: size,
		kind: pushButtonSwitch,
	}
}

// IsPushButton returns true if this is a push-button switch.
func (w *Switch) IsPushButton() bool { return w.kind == pushButtonSwitch }

// NewToggle returns a new toggle switch.
func NewToggle(x, y int, size int, def bool) *Switch {
	return &Switch{
		pos:        image.Point{x, y},
		size:       size,
		kind:       toggleSwitch,
		defEnabled: def,
		isEnabled:  def,
	}
}

// IsToggle returns true if this is a toggle switch.
func (w *Switch) IsToggle() bool { return w.kind == toggleSwitch }

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
