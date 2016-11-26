package vnc

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
	return w.Read(v)
}

func (v *VNC) Update(w Widget, val interface{}) error {
	return w.Update(v, val)
}

func (v *VNC) Press(w Widget) error {
	return w.Press(v)
}
