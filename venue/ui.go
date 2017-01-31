package venue

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/venue/router/controls"
	"github.com/kward/venue/router/signals"
	"github.com/kward/venue/venue/pages"
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
	if _, err := ui.selectPage(v, pages.Inputs); err != nil {
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
func (ui *UI) selectOutput(v *vnc.VNC, sig signals.Signal, sigNo signals.SignalNo) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	ctrlName := signalControlName(sig, sigNo)
	if ctrlName == "Invalid" {
		return fmt.Errorf("invalid control name for %s %d signal combination", sig, sigNo)
	}

	// Select outputs page.
	page, err := ui.selectPage(v, pages.Outputs)
	if err != nil {
		return err
	}

	// Clear output solo.
	if glog.V(2) {
		glog.Infof("Clearing output solo.")
	}
	widget, err := page.Widget("SoloClear")
	if err != nil {
		return err
	}
	if err := widget.Press(v); err != nil {
		return err
	}

	// Solo output.
	if glog.V(2) {
		glog.Infof("Soloing %s output.", ctrlName)
	}
	widget, err = page.Widget(ctrlName + " Solo")
	if err != nil {
		return err
	}
	if err := widget.Press(v); err != nil {
		return err
	}

	return nil
}

func (ui *UI) setOutputLevel(v *vnc.VNC, sig signals.Signal, sigNo signals.SignalNo, val int) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	ctrlName := signalControlName(sig, sigNo)
	if ctrlName == "Invalid" {
		return fmt.Errorf("invalid control name for %s %d signal combination", sig, sigNo)
	}

	if err := ui.selectOutput(v, sig, sigNo); err != nil {
		return err
	}
	page, err := ui.selectPage(v, pages.Inputs)
	if err != nil {
		return err
	}
	ctrl, err := page.Widget(ctrlName)
	if err != nil {
		return err
	}
	if err := ctrl.(*Encoder).Adjust(v, val); err != nil {
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

// signalControlName returns a control name for a `signal` and `signalNo`
// combination.
func signalControlName(sig signals.Signal, sigNo signals.SignalNo) string {
	switch sig {
	case signals.Input, signals.FXReturn:
		return controls.Fader.String()
	case signals.Direct:
		return "Invalid"
	case signals.Aux:
		return fmt.Sprintf("%s %d", controls.Aux, sigNo)
	case signals.Group:
		return fmt.Sprintf("%s %d", controls.Group, sigNo)
	default:
		return controls.Unknown.String()
	}
}
