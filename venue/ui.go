package venue

import (
	"fmt"
	"image"

	"github.com/golang/glog"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/venue/math"
	"github.com/kward/venue/router/controls"
	"github.com/kward/venue/venue/encoders"
	"github.com/kward/venue/venue/meters"
	"github.com/kward/venue/venue/pages"
	"github.com/kward/venue/venue/switches"
	"github.com/kward/venue/venuelib"
	"github.com/kward/venue/vnc"
)

// UI holds references to the Venue UI pages.
// TODO(kward:20170201) Can I get rid of the UI struct?
type UI struct {
	pages Pages
}

// NewUI returns a populated UI struct.
func NewUI() *UI {
	return &UI{
		Pages{
			pages.Inputs:  NewInputsPage(),
			pages.Outputs: NewOutputsPage(),
		},
	}
}

// selectPage changes the VENUE page.
func (ui *UI) selectPage(v *vnc.VNC, p pages.Page) (*Page, error) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	w := ui.pages[p]
	if err := w.Press(v); err != nil {
		return nil, err
	}
	return w, nil
}

// The Widget interface provides functionality for interacting with VNC widgets.
type Widget interface {
	// Press the widget (if possible).
	Press(v *vnc.VNC) error
	// Read the value of a widget.
	Read(v *vnc.VNC) (interface{}, error)
	// Update the value of a widget.
	Update(v *vnc.VNC, val interface{}) error
}

// Widgets holds references to a grouping of UI widgets.
type Widgets map[string]Widget

//-----------------------------------------------------------------------------
// Encoder

type Encoder struct {
	center   image.Point
	window   encoders.Encoder // Position of value window
	hasOnOff bool             // Has an on/off switch
}

// Verify that the expected interface is implemented properly.
var _ Widget = new(Encoder)

// Press implements the Widget interface.
func (w *Encoder) Press(v *vnc.VNC) error {
	return v.MouseLeftClick(w.clickOffset())
}

// Read implements the Widget interface.
func (w *Encoder) Read(v *vnc.VNC) (interface{}, error) { return nil, nil }

// Update implements the Widget interface.
func (w *Encoder) Update(v *vnc.VNC, val interface{}) error {
	w.Press(v)
	for _, key := range keys.IntToKeys(val.(int)) {
		v.KeyPress(key)
	}
	v.KeyPress(keys.Return)
	return nil
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
	key := keys.Up
	if val < 0 {
		key = keys.Down
	}
	amount := math.Abs(val)
	for i := 0; i < amount; i++ {
		if err := v.KeyPress(key); err != nil {
			return err
		}
	}
	return v.KeyPress(keys.Return)
}

// Increment the value of an encoder.
func (w *Encoder) Increment(v *vnc.VNC) error { return w.Adjust(v, 1) }

// Decrement the value of an encoder.
func (w *Encoder) Decrement(v *vnc.VNC) error { return w.Adjust(v, -1) }

func (w *Encoder) clickOffset() image.Point {
	var dx, dy int

	// Horizontal position.
	switch w.window {
	case encoders.TopLeft, encoders.MiddleLeft, encoders.BottomLeft: // Left
		dx = -38
	case encoders.TopRight, encoders.MiddleRight, encoders.BottomRight: // Right
		dx = 38
	default: // Middle
		dx = 0
	}

	// Vertical position.
	switch w.window {
	case encoders.TopLeft, encoders.TopRight: // Top
		dy = -8
	case encoders.MiddleLeft, encoders.MiddleRight: // Middle
		dy = 0
	case encoders.BottomLeft, encoders.BottomRight: // Bottom Left/Right
		dy = 8
	default: // Center
		dy = 28
	}

	return w.center.Add(image.Point{dx, dy})
}

//-----------------------------------------------------------------------------
// Meter

type Meter struct {
	pos      image.Point  // Position of UI element.
	size     meters.Meter // Meter size (small..large).
	isStereo bool         // True if the meter is stereo.
}

// Verify that the expected interface is implemented properly.
var _ Widget = new(Meter)

// Press implements the Widget interface.
func (w *Meter) Press(v *vnc.VNC) error {
	return v.MouseLeftClick(w.clickOffset())
}

// Read implements the Widget interface.
func (w *Meter) Read(v *vnc.VNC) (interface{}, error) { return nil, nil }

// Update implements the Widget interface.
func (w *Meter) Update(v *vnc.VNC, val interface{}) error { return nil }

// IsMono returns true if this a mono meter.
func (w *Meter) IsMono() bool { return !w.isStereo }

// IsMono returns true if this a stereo meter.
func (w *Meter) IsStereo() bool { return w.isStereo }

func (w *Meter) clickOffset() image.Point {
	switch w.size {
	case meters.SmallVertical:
		return w.pos.Add(image.Point{7, 25})
	case meters.MediumHorizontal:
		return w.pos.Add(image.Point{0, 0})
	case meters.LargeVertical:
		return w.pos.Add(image.Point{0, 0})
	}
	return w.pos
}

//-----------------------------------------------------------------------------
// Switch

const ( // Switch enabled.
	SwitchOff = false
	SwitchOn  = true
)

// Switch is a UIElement representing a switch.
type Switch struct {
	pos        image.Point         // Position of UI element.
	size       switches.SwitchSize // Switch size (tiny..large).
	kind       switches.SwitchKind // true == toggle, false == push-button
	defEnabled bool                // Default switch state.
	isEnabled  bool                // True if the switch state is enabled.
}

// Verify that the expected interface is implemented properly.
var _ Widget = new(Switch)

// Press implements the Widget interface.
func (w *Switch) Press(v *vnc.VNC) error {
	return v.MouseLeftClick(w.clickOffset())
}

// Read implements the Widget interface.
func (w *Switch) Read(v *vnc.VNC) (interface{}, error) { return nil, nil }

// Update implements the Widget interface.
func (w *Switch) Update(v *vnc.VNC, val interface{}) error {
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

// IsEnabled returns true if the switch is enabled.
func (w *Switch) IsEnabled() bool { return w.isEnabled }

// NewPushButton returns a new push-button switch.
func NewPushButton(x, y int, size switches.SwitchSize) *Switch {
	return &Switch{
		pos:  image.Point{x, y},
		size: size,
		kind: switches.PushButton,
	}
}

// IsPushButton returns true if this is a push-button switch.
func (w *Switch) IsPushButton() bool { return w.kind == switches.PushButton }

// NewToggle returns a new toggle switch.
func NewToggle(x, y int, size switches.SwitchSize, def bool) *Switch {
	return &Switch{
		pos:        image.Point{x, y},
		size:       size,
		kind:       switches.Toggle,
		defEnabled: def,
		isEnabled:  def,
	}
}

// IsToggle returns true if this is a toggle switch.
func (w *Switch) IsToggle() bool { return w.kind == switches.Toggle }

func (e *Switch) clickOffset() image.Point {
	switch e.size {
	case switches.Tiny:
		return e.pos.Add(image.Point{7, 7})
	case switches.Small:
		return e.pos.Add(image.Point{9, 7})
	case switches.Medium:
		return e.pos.Add(image.Point{13, 8})
	case switches.Large:
		return e.pos.Add(image.Point{16, 9})
	}
	return e.pos
}

//-----------------------------------------------------------------------------
// Page

// Page holds references to the UI elements on a VENUE page.
type Page struct {
	page    pages.Page
	widgets map[string]Widget
}

// Pages is a map of Page keyed on pages.Page.
type Pages map[pages.Page]*Page

// Verify that the expected interface is implemented properly.
var _ Widget = new(Page)

// Press implements the Widget interface.
func (w *Page) Press(v *vnc.VNC) error {
	var key keys.Key
	switch w.page {
	case pages.Inputs:
		key = keys.F1
	case pages.Outputs:
		key = keys.F2
	}
	if err := v.KeyPress(key); err != nil {
		return err
	}
	return nil
}

// Read implements the Widget interface.
func (w *Page) Read(v *vnc.VNC) (interface{}, error) {
	return nil, fmt.Errorf("page.Read() is unsupported")
}

// Update implements the Widget interface.
func (w *Page) Update(v *vnc.VNC, val interface{}) error {
	return fmt.Errorf("Page.Update() is unsupported")
}

const (
	bankX  = 8   // X position of 1st bank.
	bankDX = 131 // dX between banks.
	chanDX = 15  // dX between channels in a bank.

	// Inputs
	auxOddX  = 316
	auxPanX  = 473
	aux12Y   = 95
	aux34Y   = 146
	aux56Y   = 197
	aux78Y   = 248
	aux910Y  = 299
	aux1112Y = 350
	aux1314Y = 401
	aux1516Y = 452

	// Outputs
	meterY = 512
	muteY  = 588
	soloY  = 573
)

// NewInputsPage returns a populated Inputs page.
func NewInputsPage() *Page {
	return &Page{
		pages.Inputs,
		Widgets{
			"Gain":         &Encoder{image.Point{167, 279}, encoders.BottomLeft, true},
			"Delay":        &Encoder{image.Point{168, 387}, encoders.BottomLeft, false},
			"HPF":          &Encoder{image.Point{168, 454}, encoders.BottomLeft, true},
			"Pan":          &Encoder{image.Point{239, 443}, encoders.BottomCenter, false},
			"VarGroups":    NewPushButton(226, 299, switches.Medium),
			"Aux 1":        &Encoder{image.Point{auxOddX, aux12Y}, encoders.TopRight, true},
			"AuxPan 1/2":   &Encoder{image.Point{auxPanX, aux12Y}, encoders.TopLeft, false},
			"Aux 3":        &Encoder{image.Point{auxOddX, aux34Y}, encoders.TopRight, true},
			"AuxPan 3/4":   &Encoder{image.Point{auxPanX, aux34Y}, encoders.TopLeft, false},
			"Aux 5":        &Encoder{image.Point{auxOddX, aux56Y}, encoders.TopRight, true},
			"AuxPan 5/6":   &Encoder{image.Point{auxPanX, aux56Y}, encoders.TopLeft, false},
			"Aux 7":        &Encoder{image.Point{auxOddX, aux78Y}, encoders.TopRight, true},
			"AuxPan 7/8":   &Encoder{image.Point{auxPanX, aux78Y}, encoders.TopLeft, false},
			"Aux 9":        &Encoder{image.Point{auxOddX, aux910Y}, encoders.TopRight, true},
			"AuxPan 9/10":  &Encoder{image.Point{auxPanX, aux910Y}, encoders.TopLeft, false},
			"Aux 11":       &Encoder{image.Point{auxOddX, aux1112Y}, encoders.TopRight, true},
			"AuxPan 11/12": &Encoder{image.Point{auxPanX, aux1112Y}, encoders.TopLeft, false},
			"Aux 13":       &Encoder{image.Point{auxOddX, aux1314Y}, encoders.TopRight, true},
			"AuxPan 13/14": &Encoder{image.Point{auxPanX, aux1314Y}, encoders.TopLeft, false},
			"Aux 15":       &Encoder{image.Point{auxOddX, aux1516Y}, encoders.TopRight, true},
			"AuxPan 15/16": &Encoder{image.Point{auxPanX, aux1516Y}, encoders.TopLeft, false},
			"Group 1":      &Encoder{image.Point{auxOddX, aux12Y}, encoders.TopRight, true},
			"GroupPan 1/2": &Encoder{image.Point{auxPanX, aux12Y}, encoders.TopLeft, false},
			"Group 3":      &Encoder{image.Point{auxOddX, aux34Y}, encoders.TopRight, true},
			"GroupPan 3/4": &Encoder{image.Point{auxPanX, aux34Y}, encoders.TopLeft, false},
			"Group 5":      &Encoder{image.Point{auxOddX, aux56Y}, encoders.TopRight, true},
			"GroupPan 5/6": &Encoder{image.Point{auxPanX, aux56Y}, encoders.TopLeft, false},
			"Group 7":      &Encoder{image.Point{auxOddX, aux78Y}, encoders.TopRight, true},
			"GroupPan 7/8": &Encoder{image.Point{auxPanX, aux78Y}, encoders.TopLeft, false},
			"SoloClear":    NewPushButton(979, 493, switches.Medium),
		}}
}

// NewOutputsPage returns a populated Outputs page.
func NewOutputsPage() *Page {
	widgets := Widgets{
		"SoloClear": NewPushButton(980, 490, switches.Medium),
	}

	// Auxes
	for _, b := range []int{1, 2} { // Bank.
		pre := controls.Aux.String()
		for c := 1; c <= 8; c++ { // Bank channel.
			ch, x := (b-1)*8+c, bankX+(b-1)*bankDX+(c-1)*chanDX

			n := fmt.Sprintf("%s %d Solo", pre, ch)
			if glog.V(4) {
				glog.Infof("NewOutput() element[%v]:", n)
			}
			widgets[n] = NewToggle(x, soloY, switches.Tiny, false)

			n = fmt.Sprintf("%s %d Value", pre, ch)
			if glog.V(4) {
				glog.Infof("NewOutput() element[%v]:", n)
			}
			widgets[n] = &Meter{
				pos:  image.Point{x, meterY},
				size: meters.SmallVertical,
			}
		}
	}

	// Groups
	b := 5 // bank
	pre := controls.Group.String()
	for c := 1; c <= 8; c++ { // bank channel
		ch, x := c, bankX+(b-1)*bankDX+(c-1)*chanDX

		n := fmt.Sprintf("%s %d Solo", pre, ch)
		if glog.V(4) {
			glog.Infof("NewOutput() element[%v]:", n)
		}
		widgets[n] = NewToggle(x, soloY, switches.Tiny, false)

		n = fmt.Sprintf("%s %d Meter", pre, ch)
		if glog.V(4) {
			glog.Infof("NewOutput() element[%v]:", n)
		}
		widgets[n] = &Meter{
			pos:  image.Point{x, meterY},
			size: meters.SmallVertical,
		}
	}

	return &Page{pages.Outputs, widgets}
}

// Widget returns the named widget.
func (w *Page) Widget(n string) (Widget, error) {
	v, ok := w.widgets[n]
	if !ok {
		return nil, fmt.Errorf("invalid page widget %q", n)
	}
	return v, nil
}
