package vnc

import "fmt"

const (
	// Max number of consecutive arrow key presses.
	// WiFi connections have enough latency to not need more.
	maxArrowKeys = 4
)

type UI struct {
	inputs  *Page
	outputs *Page
}

func NewUI() *UI {
	return &UI{NewInputsPage(), NewOutputsPage()}
}

func (ui *UI) Inputs() *Page  { return ui.inputs }
func (ui *UI) Outputs() *Page { return ui.outputs }

// The Widget interface provides functionality for interacting with VNC widgets.
type Widget interface {
	// Read the value of a widget.
	Read(v *VNC) (interface{}, error)

	// Update the value of a widget.
	Update(v *VNC, val interface{}) error

	// Press the widget (if possible).
	Press(v *VNC) error
}

type Widgets map[string]Widget

// Note, although similar, VNC does *not* meet the Widget interface.

func (v *VNC) Read(w Widget) (interface{}, error) {
	if w == nil {
		return nil, fmt.Errorf("Widget is nil")
	}
	return w.Read(v)
}

func (v *VNC) Update(w Widget, val interface{}) error {
	if w == nil {
		return fmt.Errorf("Widget is nil")
	}
	return w.Update(v, val)
}

func (v *VNC) Press(w Widget) error {
	if w == nil {
		return fmt.Errorf("Widget is nil")
	}
	return w.Press(v)
}
