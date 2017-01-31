package venue

import (
	"context"
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

// Endpoint handlers.
var handlers router.Handlers

func init() {
	specs := []router.HandlerSpec{
		router.HandlerSpec{actions.Noop, noop},
		router.HandlerSpec{actions.Ping, ping},
		router.HandlerSpec{actions.SelectInput, selectInput},
		router.HandlerSpec{actions.SelectOutput, selectOutput},
		router.HandlerSpec{actions.SetOutputLevel, setOutputLevel},
	}
	handlers = make(router.Handlers, len(specs))
	for _, spec := range specs {
		handlers[spec.Action] = spec
	}
}

// Venue holds information representing the state of the VENUE backend.
type Venue struct {
	opts *options

	vnc *vnc.VNC

	ui       *UI
	currPage pages.Page
	inputs   [numInputs]*Input
	outputs  map[string]*Output
}

// Verify that the expected interface is implemented properly.
var _ router.Endpoint = new(Venue)

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
	if err := v.ui.selectOutput(v.vnc, signals.Aux, 1); err != nil {
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
	w, err := p.Widget("SoloClear")
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
	router.Handle(v, pkt, handlers)
}

// noop is a noop packet.
func noop(_ router.Endpoint, _ *router.Packet) {}

// Ping is deprecated. This should move to the OSC module, passed on a channel.
func ping(ep router.Endpoint, _ *router.Packet) {
	ep.(*Venue).vnc.DebugMetrics()
}

// selectInput for adjustment.
func selectInput(ep router.Endpoint, pkt *router.Packet) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	v := ep.(*Venue)
	if err := v.ui.selectInput(v.vnc, uint16(pkt.SignalNo)); err != nil {
		glog.Errorf("unable to select input; %s", err)
	}
}

// selectOutput for adjustment.
func selectOutput(ep router.Endpoint, pkt *router.Packet) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	v := ep.(*Venue)
	if err := v.ui.selectOutput(v.vnc, pkt.Signal, pkt.SignalNo); err != nil {
		glog.Errorf("unable to select output for %s %d; %s", pkt.Signal, pkt.SignalNo, err)
		return
	}
}

// setOutputLevel for the specified output. This handler operates on the
// currently selected input.
func setOutputLevel(ep router.Endpoint, pkt *router.Packet) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	v := ep.(*Venue)
	if err := v.ui.setOutputLevel(v.vnc, pkt.Signal, pkt.SignalNo, pkt.Value.(int)); err != nil {
		glog.Errorf("unable to set output level for %s %d; %s", pkt.Signal, pkt.SignalNo, err)
		return
	}
}