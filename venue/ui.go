package venue

import (
	"fmt"
	"image"

	"github.com/golang/glog"
	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/venue/codes"
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
	return &UI{Pages{
		pages.Inputs:  NewInputsPage(),
		pages.Outputs: NewOutputsPage(),
	}}
}

// selectPage changes the VENUE page.
func (ui *UI) selectPage(wf *vnc.Workflow, p pages.Page) (*Page, error) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	w, ok := ui.pages[p]
	if !ok {
		return nil, venuelib.Errorf(codes.Unimplemented, "support for %q page unimplemented")
	}
	if p == pages.Inputs {
		// To ensure we start on inputs bank 1-48, select another page first.
		wf.KeyPress(keys.F2) // OUTPUTS
	}
	if err := w.Press(wf); err != nil {
		return nil, err
	}
	return w, nil
}

// The Widget interface provides functionality for interacting with VNC widgets.
type Widget interface {
	// Press the widget (if possible).
	Press(wf *vnc.Workflow) error
	// Read the value of a widget.
	Read(wf *vnc.Workflow) (interface{}, error)
	// Update the value of a widget.
	Update(wf *vnc.Workflow, val interface{}) error
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
func (w *Encoder) Press(wf *vnc.Workflow) error {
	wf.MouseClick(buttons.Left, w.clickPoint())
	return nil
}

// Read implements the Widget interface.
func (w *Encoder) Read(wf *vnc.Workflow) (interface{}, error) {
	return nil, venuelib.Errorf(codes.Unimplemented, "Encoder.Read() unimplemented")
}

// Update implements the Widget interface.
func (w *Encoder) Update(wf *vnc.Workflow, val interface{}) error {
	w.Press(wf)
	for _, key := range keys.IntToKeys(val.(int)) {
		wf.KeyPress(key)
	}
	wf.KeyPress(keys.Return)
	return nil
}

// Adjust the value of an encoder with cursor keys.
func (w *Encoder) Adjust(wf *vnc.Workflow, val int) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Adjusting encoder by %d steps.", val)
	}
	if val == 0 {
		return nil
	}

	if err := w.Press(wf); err != nil {
		return err
	}
	key := keys.Up
	if val < 0 {
		key = keys.Down
	}
	amount := math.Abs(val)
	for i := 0; i < amount; i++ {
		wf.KeyPress(key)
	}
	wf.KeyPress(keys.Return)

	return nil
}

// Increment the value of an encoder.
func (w *Encoder) Increment(wf *vnc.Workflow) error { return w.Adjust(wf, 1) }

// Decrement the value of an encoder.
func (w *Encoder) Decrement(wf *vnc.Workflow) error { return w.Adjust(wf, -1) }

// clickPoint returns the point to click based on the window of the encoder.
func (w *Encoder) clickPoint() image.Point {
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
func (w *Meter) Press(wf *vnc.Workflow) error {
	wf.MouseClick(buttons.Left, w.clickOffset())
	return nil
}

// Read implements the Widget interface.
func (w *Meter) Read(wf *vnc.Workflow) (interface{}, error) {
	return nil, venuelib.Errorf(codes.Unimplemented, "Meter.Read() unimplemented")
}

// Update implements the Widget interface.
func (w *Meter) Update(wf *vnc.Workflow, val interface{}) error {
	return venuelib.Errorf(codes.Unimplemented, "Meter.Update() unimplemented")
}

// IsMono returns true if this a mono meter.
func (w *Meter) IsMono() bool { return !w.isStereo }

// IsMono returns true if this a stereo meter.
func (w *Meter) IsStereo() bool { return w.isStereo }

// clickPoint returns the point to click based on the size of the meter.
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

// Switch is a UIElement representing a switch.
type Switch struct {
	pos        image.Point   // Position of UI element.
	size       switches.Size // Switch size (Tiny..Large).
	kind       switches.Kind // Switch kind (Toggle, PushButton).
	defEnabled bool          // Default switch state (Enabled == true).
	isEnabled  bool          // True if the switch state is enabled.
}

// Verify that the expected interface is implemented properly.
var _ Widget = new(Switch)

// Press implements the Widget interface.
func (w *Switch) Press(wf *vnc.Workflow) error {
	wf.MouseClick(buttons.Left, w.clickOffset())
	return nil
}

// Read implements the Widget interface.
func (w *Switch) Read(wf *vnc.Workflow) (interface{}, error) {
	return nil, venuelib.Errorf(codes.Unimplemented, "Switch.Read() unimplemented")
}

// Update implements the Widget interface.
func (w *Switch) Update(wf *vnc.Workflow, val interface{}) error {
	if w.IsPushButton() {
		return venuelib.Errorf(codes.InvalidArgument, "switches cannot be updated")
	}

	val, err := w.Read(wf)
	if err != nil {
		return err
	}
	if w.isEnabled != val.(bool) {
		return w.Press(wf)
	}
	return nil
}

// IsEnabled returns true if the switch is enabled.
func (w *Switch) IsEnabled() bool { return w.isEnabled }

// NewPushButton returns a new push-button switch. x and y refer to the
// top-left corner of the switch.
func NewPushButton(x, y int, size switches.Size) *Switch {
	return &Switch{
		pos:  image.Point{x, y},
		size: size,
		kind: switches.PushButton,
	}
}

// IsPushButton returns true if this is a push-button switch.
func (w *Switch) IsPushButton() bool { return w.kind == switches.PushButton }

// NewToggle returns a new toggle switch. x and y refer to the top-left corner
// of the switch.
func NewToggle(x, y int, size switches.Size, def bool) *Switch {
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

type Pages map[pages.Page]*Page

// Verify that the expected interface is implemented properly.
var _ Widget = new(Page)

// Press implements the Widget interface.
func (w *Page) Press(wf *vnc.Workflow) error {
	var key keys.Key
	switch w.page {
	case pages.Inputs:
		key = keys.F1
	case pages.Outputs:
		key = keys.F2
	}
	wf.KeyPress(key)
	return nil
}

// Read implements the Widget interface.
func (w *Page) Read(wf *vnc.Workflow) (interface{}, error) {
	return nil, venuelib.Errorf(codes.Unimplemented, "Page.Read() unimplemented")
}

// Update implements the Widget interface.
func (w *Page) Update(wf *vnc.Workflow, val interface{}) error {
	return venuelib.Errorf(codes.Unimplemented, "Page.Update() unimplemented")
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
			// Input
			"Phantom": NewToggle(153, 171, switches.Medium, switches.Disabled),
			"Pad":     NewToggle(153, 196, switches.Medium, switches.Disabled),
			"Guess":   NewPushButton(153, 221, switches.Medium),
			"Gain":    &Encoder{image.Point{167, 279}, encoders.BottomLeft, true},
			"Phase":   NewToggle(12, 420, switches.Medium, switches.Disabled),
			"Solo":    NewToggle(12, 451, switches.Large, switches.Disabled),
			"Mute":    NewToggle(62, 451, switches.Large, switches.Disabled),
			"Delay":   &Encoder{image.Point{168, 387}, encoders.BottomLeft, false},
			"HPF":     &Encoder{image.Point{168, 454}, encoders.BottomLeft, true},
			// Bus Assign
			"VarGroups": NewPushButton(226, 299, switches.Medium),
			// Pan
			"Pan": &Encoder{image.Point{239, 443}, encoders.BottomCenter, false},
			//-- RightOffset
			//-- Balance
			// Aux Sends
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
			// EQ
			// Comp/Lim
			// Exp/Gate
			// Misc
			"SoloClear": NewPushButton(979, 493, switches.Medium),
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
		return nil, venuelib.Errorf(codes.Internal, "invalid %q page widget", n)
	}
	return v, nil
}
