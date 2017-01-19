/*
The Venue package exposes the Avid™ VENUE VNC interface as a programmatic API.
*/
package venue

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/kward/venue/vnc"
)

const (
	refresh   = 1000 * time.Millisecond
	numInputs = 48
)

// Venue holds information representing the state of the VENUE backend.
type Venue struct {
	opts *options

	VNC *vnc.VNC

	inputs    [numInputs]*Input
	currInput *Input

	outputs    map[string]*Output
	currOutput *Output

	currPage int
}

// New returns a populated Venue struct.
func New(opts ...func(*options) error) (*Venue, error) {
	o := &options{}
	o.setInputs(numInputs)
	o.setRefresh(refresh)
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	return &Venue{opts: o}, nil
}

// Close a Venue session.
func (v *Venue) Close() error {
	return v.VNC.Close()
}

// Connect to a VENUE VNC server.
func (v *Venue) Connect(ctx context.Context, h string, p uint, pw string) error {
	// Establish a connection to the VENUE VNC server.
	handle, err := vnc.New(vnc.Host(h), vnc.Port(p), vnc.Password(pw))
	if err != nil {
		return err
	}
	if err := handle.Connect(ctx); err != nil {
		return err
	}
	v.VNC = handle
	return nil
}

// Initialize the in-memory state representation of a VENUE console.
func (v *Venue) Initialize() error {
	glog.Info("Initialize()")

	// Initialize inputs.
	glog.Info("Initializing inputs.")
	for ch := 0; ch < numInputs; ch++ {
		input := NewInput(v, ch+1, Ichannel)
		v.inputs[ch] = input
	}

	// Choose output before input so that later when the Inputs page is selected,
	// it shows first bank of channels.
	glog.Info("Selecting I/O.")
	if err := v.VNC.SelectOutput("aux1"); err != nil {
		return err
	}
	if err := v.VNC.SelectInput(1); err != nil {
		return err
	}

	// Clear solo.
	glog.Info("Clearing solo.")
	// TODO(kward:20170120) use something generated instead of solo_clear string.
	widget := v.VNC.Widget(vnc.InputsPage, "solo_clear")
	if err := v.VNC.Update(widget, vnc.SwitchOff); err != nil {
		return err
	}

	return nil
}

// ListenAndHandle connections and incoming requests.
func (v *Venue) ListenAndHandle() {
	go v.VNC.ListenAndHandle()
	go v.VNC.FramebufferRefresh(v.opts.refresh)
}

// Ping is deprecated. This should move to the OSC module, passed on a channel.
func (v *Venue) Ping() {
	v.VNC.DebugMetrics()
}
