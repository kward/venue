package venue

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/venue/venuelib"
	"github.com/kward/venue/vnc"
)

const (
	// Maximum number of consecutive arrow key presses.
	// WiFi connections have low enough latency to not need more.
	maxArrowKeys = 4
	// Maximum number of signal inputs the code can handle.
	maxInputs = 96
	// The amount of time to delay after a keyboard input was made. It takes this
	// long for the VENUE UI to stop waiting for additional input.
	inputWait = 1750 * time.Millisecond
)

// UI holds references to the Venue UI pages.
type UI struct {
	pages Pages
}

// NewUI returns a populated UI struct.
func NewUI() *UI {
	return &UI{
		Pages{
			inputsPage:  NewInputsPage(),
			outputsPage: NewOutputsPage(),
		},
	}
}

// selectPage changes the VENUE page.
func (ui *UI) selectPage(v *vnc.VNC, p pageEnum) (*Page, error) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	w := ui.pages[p]
	if err := w.Press(v); err != nil {
		return nil, err
	}
	return w, nil
}

// selectInput changes the input channel to operate on.
func (ui *UI) selectInput(v *vnc.VNC, input uint16) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Selecting input #%v.", input)
	}

	if input > maxInputs {
		return fmt.Errorf("input number %d exceeds maximum number of inputs %d", input, maxInputs)
	}
	if _, err := ui.selectPage(v, inputsPage); err != nil {
		return err
	}

	keys := []uint32{}
	if input < 10 {
		keys = append(keys, vnclib.Key0)
	}
	keys = append(keys, vnc.IntToKeys(int(input))...)
	for _, key := range keys {
		if err := v.KeyPress(key); err != nil {
			return err
		}
	}
	// TODO(kward:20161126): Start a timer that expires after 1750ms. Additional
	// key presses aren't allowed until the time expires, but mouse input is.
	time.Sleep(inputWait)

	return nil
}

// selectOutput changes the output channel to operate on.
func (ui *UI) selectOutput(v *vnc.VNC, output string) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	// Select outputs page.
	p, err := ui.selectPage(v, outputsPage)
	if err != nil {
		return err
	}

	// Clear output solo.
	if glog.V(2) {
		glog.Infof("Clearing output solo.")
	}
	w, err := p.Widget("solo_clear")
	if err != nil {
		return err
	}
	if err := w.Press(v); err != nil {
		return err
	}

	// Solo output.
	if glog.V(2) {
		glog.Infof("Soloing %v output.", output)
	}
	w, err = p.Widget(output + "solo")
	if err != nil {
		return err
	}
	if err := w.Press(v); err != nil {
		return err
	}

	return nil
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

const (
	WidgetAux   = "aux"
	WidgetGroup = "grp"
)
