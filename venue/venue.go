/*
The Venue package exposes the Avid™ VENUE VNC interface as a programmatic API.
*/
package venue

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/kward/venue/router"
	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/router/signals"
	"github.com/kward/venue/venue/pages"
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
	currPage pages.Page
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

// Verify that expected interfaces are implemented properly.
var _ router.Endpoint = new(Venue)

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

	if glog.V(2) {
		glog.Infof("Clearing input solo.")
	}
	p, err := v.ui.selectPage(v.vnc, pages.Inputs)
	if err != nil {
		return err
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

// EndpointName implements router.Endpoint.
func (v *Venue) EndpointName() string { return "Venue" }

// Handle implements router.Endpoint.
func (v *Venue) Handle(pkt *router.Packet) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if pkt == nil { // Ignore nil packets.
		return
	}
	if glog.V(2) {
		glog.Infof("Handling %s packet.", pkt.Action)
	}

	switch pkt.Action {
	case actions.Ping:
		v.Ping()
	case actions.SelectInput:
		v.SelectInput(pkt.SignalNo)
	case actions.SelectOutput:
		v.SelectOutput(pkt.Signal, pkt.SignalNo)
	case actions.DropPacket: // Do nothing.
	default:
		glog.Errorf("%s action unimplemented.", pkt.Action)
	}
}

// Ping is deprecated. This should move to the OSC module, passed on a channel.
func (v *Venue) Ping() {
	v.vnc.DebugMetrics()
}

// SelectInput for adjustment.
func (v *Venue) SelectInput(input int) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if err := v.ui.selectInput(v.vnc, uint16(input)); err != nil {
		glog.Errorf("unable to select input; %s", err)
	}
}

// SelectOutput for adjustment.
func (v *Venue) SelectOutput(sig signals.Signal, sigNo int) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	var name string
	switch sig {
	case signals.Aux:
		name = WidgetAux
	case signals.Group:
		name = WidgetGroup
	}
	if err := v.ui.selectOutput(v.vnc, fmt.Sprintf("%s%d", name, sigNo)); err != nil {
		glog.Errorf("unable to select output; %s", err)
	}
}
