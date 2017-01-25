/*
The Venue package exposes the Avid™ VENUE VNC interface as a programmatic API.
*/
package venue

import (
	"context"
	"time"

	"github.com/golang/glog"
	"github.com/kward/venue/venuelib"
	"github.com/kward/venue/vnc"
)

const (
	refresh   = 1000 * time.Millisecond
	numInputs = 48
)

// Venue holds information representing the state of the VENUE backend.
type Venue struct {
	opts *options

	vnc *vnc.VNC

	ui       *UI
	currPage pageEnum
	inputs   [numInputs]*Input
	outputs  map[string]*Output
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
	return v.vnc.Close()
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
	v.vnc = handle
	return nil
}

// Initialize the in-memory state representation of a VENUE console.
func (v *Venue) Initialize() error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	v.ui = NewUI()

	// Initialize inputs.
	if glog.V(2) {
		glog.Info("Initializing inputs.")
	}
	for ch := 0; ch < numInputs; ch++ {
		input := NewInput(v, ch+1, signalChannel)
		v.inputs[ch] = input
	}

	// Choose output before input so that later when the Inputs page is selected,
	// it shows first bank of channels.
	if glog.V(2) {
		glog.Info("Selecting I/O.")
	}
	if err := v.ui.selectOutput(v.vnc, "aux1"); err != nil {
		return err
	}
	if err := v.ui.selectInput(v.vnc, 1); err != nil {
		return err
	}

	p, err := v.ui.selectPage(v.vnc, inputsPage)
	if err != nil {
		return err
	}

	if glog.V(2) {
		glog.Infof("Clearing input solo.")
	}
	w, err := p.Widget("solo_clear")
	if err != nil {
		return err
	}
	if err := w.Press(v.vnc); err != nil {
		return err
	}

	return nil
}

// ListenAndHandle connections and incoming requests.
func (v *Venue) ListenAndHandle() {
	go v.vnc.ListenAndHandle()
	go v.vnc.FramebufferRefresh(v.opts.refresh)
}

// Ping is deprecated. This should move to the OSC module, passed on a channel.
func (v *Venue) Ping() {
	v.vnc.DebugMetrics()
}

func (v *Venue) SelectInput(input uint16) error {
	return v.ui.selectInput(v.vnc, input)
}
