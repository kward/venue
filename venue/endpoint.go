package venue

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/venue/router"
	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/router/controls"
	"github.com/kward/venue/router/signals"
	"github.com/kward/venue/venue/pages"
	"github.com/kward/venue/venuelib"
	"github.com/kward/venue/vnc"
)

const (
	refresh   = 1000 * time.Millisecond
	numInputs = 48
	// Maximum number of consecutive arrow key presses.
	// WiFi connections have low enough latency to not need more.
	maxArrowKeys = 4
	// Maximum number of signal inputs the code can handle.
	maxInputs = 96
	// The amount of time to delay after a keyboard input was made. It takes this
	// long for the VENUE UI to stop waiting for additional input.
	inputWait = 1750 * time.Millisecond
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
	if err := selectOutput(v, &router.Packet{Signal: signals.Aux, SignalNo: 1}); err != nil {
		return err
	}
	if err := selectInput(v, &router.Packet{SignalNo: 1}); err != nil {
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
	if err := router.Handle(v, pkt, handlers); err != nil {
		glog.Errorf("Error handling %s packet; %s", pkt.Action, err)
	}
}

// noop is a noop packet.
func noop(_ router.Endpoint, _ *router.Packet) error { return nil }

// Ping is deprecated. This should move to the OSC module, passed on a channel.
func ping(ep router.Endpoint, _ *router.Packet) error {
	ep.(*Venue).vnc.DebugMetrics()
	return nil
}

// selectInput for adjustment.
func selectInput(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Selecting input #%d.", pkt.SignalNo)
	}

	if pkt.SignalNo > maxInputs {
		return fmt.Errorf("signal number %d exceeds maximum number of inputs %d", pkt.SignalNo, maxInputs)
	}

	v := ep.(*Venue)

	// Select INPUTS page.
	if _, err := v.ui.selectPage(v.vnc, pages.Inputs); err != nil {
		return err
	}

	// Type the channel number.
	ks := keys.Keys{}
	if pkt.SignalNo < 10 {
		ks = append(ks, keys.Digit0)
	}
	ks = append(ks, keys.IntToKeys(int(pkt.SignalNo))...)
	for _, k := range ks {
		if err := v.vnc.KeyPress(k); err != nil {
			return err
		}
	}

	// TODO(kward:20161126): Start a timer that expires after 1750ms. Additional
	// key presses aren't allowed until the time expires, but mouse input is.
	time.Sleep(inputWait)
	return nil
}

// selectOutput for adjustment.
func selectOutput(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Selecting %s %d output.", pkt.Signal, pkt.SignalNo)
	}

	ctrlName := signalControlName(pkt.Signal, pkt.SignalNo)
	if ctrlName == "Invalid" {
		return fmt.Errorf("invalid control name for %s %d signal combination", pkt.Signal, pkt.SignalNo)
	}

	v := ep.(*Venue)

	// Select the OUTPUTS page.
	page, err := v.ui.selectPage(v.vnc, pages.Outputs)
	if err != nil {
		return err
	}

	// Clear the output solo.
	if glog.V(2) {
		glog.Infof("Clearing output solo.")
	}
	widget, err := page.Widget("SoloClear")
	if err != nil {
		return err
	}
	if err := widget.Press(v.vnc); err != nil {
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
	if err := widget.Press(v.vnc); err != nil {
		return err
	}

	return nil
}

// setOutputLevel for the specified output. This handler operates on the
// currently selected input.
func setOutputLevel(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Adjusting %s %d output level by %d dB.", pkt.Signal, pkt.SignalNo, pkt.Value)
	}

	ctrlName := signalControlName(pkt.Signal, pkt.SignalNo)
	if ctrlName == "Invalid" {
		return fmt.Errorf("invalid control name for %s %d signal combination", pkt.Signal, pkt.SignalNo)
	}

	v := ep.(*Venue)

	// Select the INPUTS page.
	page, err := v.ui.selectPage(v.vnc, pages.Inputs)
	if err != nil {
		return err
	}

	// Adjust the Aux/Group knob.
	ctrl, err := page.Widget(ctrlName)
	if err != nil {
		return err
	}
	if err := ctrl.(*Encoder).Adjust(v.vnc, pkt.Value.(int)); err != nil {
		return err
	}

	return nil
}

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
