package venue

import (
	"context"
	"fmt"
	"image"
	"time"

	"github.com/golang/glog"
	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/venue/codes"
	"github.com/kward/venue/internal/router"
	"github.com/kward/venue/internal/router/actions"
	"github.com/kward/venue/internal/router/controls"
	"github.com/kward/venue/internal/router/signals"
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
		{Action: actions.Noop, Handler: Noop},
		{Action: actions.Ping, Handler: Ping},
		{Action: actions.SelectInput, Handler: SelectInput},
		{Action: actions.InputGain, Handler: InputGain},
		//router.HandlerSpec{actions.InputGuess, InputGuess},
		{Action: actions.InputMute, Handler: InputMute},
		{Action: actions.InputPad, Handler: InputPad},
		{Action: actions.InputPhantom, Handler: InputPhantom},
		{Action: actions.InputSolo, Handler: InputSolo},
		{Action: actions.SelectOutput, Handler: SelectOutput},
		{Action: actions.OutputLevel, Handler: OutputLevel},
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
	if glog.V(3) {
		glog.Infof("Venue.%s", venuelib.FnName())
	}
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
		glog.Infof("Venue.%s", venuelib.FnName())
	}

	v.ui = NewUI()

	// Initialize inputs.
	if glog.V(2) {
		glog.Info("Initializing inputs.")
	}
	for sigNo := 0; sigNo < numInputs; sigNo++ {
		input := NewInput(signals.Input, signals.SignalNo(sigNo+1))
		v.inputs[sigNo] = input
	}

	// Choose output before input so that later when the Inputs page is selected,
	// it shows first bank of channels.
	// TODO(kward:20170207) Remove once the console state can be determined.
	if glog.V(2) {
		glog.Info("Selecting I/O.")
	}
	if err := SelectOutput(v, &router.Packet{Signal: signals.Aux, SignalNo: 1}); err != nil {
		return err
	}
	if err := SelectInput(v, &router.Packet{SignalNo: 1}); err != nil {
		return err
	}

	wf := vnc.NewWorkflow(v.vnc.ClientConn())
	if glog.V(2) {
		glog.Infof("Clearing input solo.")
	}
	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}
	w, err := p.Widget("SoloClear")
	if err != nil {
		return err
	}
	if err := w.Press(wf); err != nil {
		return err
	}
	return wf.Execute()
}

// ListenAndHandle connections and incoming requests.
// ListenAndHandle maintains backward compatibility; prefer ListenAndHandleCtx.
func (v *Venue) ListenAndHandle() { v.ListenAndHandleCtx(context.Background()) }

// ListenAndHandleCtx starts goroutines to listen for messages and refresh the
// framebuffer until the context is cancelled.
func (v *Venue) ListenAndHandleCtx(ctx context.Context) {
	go v.vnc.ListenAndHandleCtx(ctx)
	go v.vnc.FramebufferRefreshCtx(ctx, v.opts.refresh)
}

// EndpointName implements router.Endpoint.
func (v *Venue) EndpointName() string { return "Venue" }

// Handle implements router.Endpoint.
func (v *Venue) Handle(pkt *router.Packet) {
	if glog.V(3) {
		glog.Infof("Venue.%s", venuelib.FnName())
	}
	if err := router.Handle(v, pkt, handlers); err != nil {
		glog.Errorf("Error handling %s packet; %s", pkt.Action, err)
	}
}

//-----------------------------------------------------------------------------
// router.Handler functions

// Noop is a noop packet.
func Noop(_ router.Endpoint, _ *router.Packet) error { return nil }

// Ping is deprecated. This should move to the OSC module, passed on a channel.
func Ping(ep router.Endpoint, _ *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	ep.(*Venue).vnc.DebugMetrics()
	return nil
}

// SelectInput for adjustment.
func SelectInput(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Selecting input #%d.", pkt.SignalNo)
	}

	if pkt.SignalNo > maxInputs {
		return venuelib.Errorf(codes.InvalidArgument, "signal number %d exceeds maximum number of inputs %d", pkt.SignalNo, maxInputs)
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	// Select the INPUTS page.
	_, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}

	// Select channels 1-48.
	// TODO(kward:20170226) Handle this with UI element.
	wf.MouseClick(buttons.Left, image.Point{932, 524})

	// Type the channel number.
	ks := keys.Keys{}
	if pkt.SignalNo < 10 {
		ks = append(ks, keys.Digit0)
	}
	ks = append(ks, keys.IntToKeys(int(pkt.SignalNo))...)
	for _, k := range ks {
		wf.KeyPress(k)
	}
	// TODO(kward:20161126): Start a timer that expires after 1750ms. Additional
	// key presses aren't allowed until the time expires, but mouse input is.
	wf.Sleep(inputWait)

	return wf.Execute()
}

// InputGain adjustment.
func InputGain(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Adjusting input gain by %d dB.", pkt.Value)
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	// Select the INPUTS page.
	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}
	w, err := p.Widget("Gain")
	if err != nil {
		return err
	}
	if err := w.(*Encoder).Adjust(wf, pkt.Value.(int)); err != nil {
		return err
	}

	return wf.Execute()
}

// InputGuess button push.
func InputGuess(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Info("Guessing input gain.")
	}

	// TODO(kward:20170209) Guess needs both press/release support as the user
	// will hold the button for some amount of time, O(seconds).
	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}
	if err := pressWidget(wf, p, "Guess"); err != nil {
		return err
	}

	return wf.Execute()
}

// InputMute button push.
func InputMute(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Info("Toggle the input mute.")
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}
	if err := pressWidget(wf, p, "Mute"); err != nil {
		return err
	}

	return wf.Execute()
}

// InputPad toggles the state of the input pad button.
func InputPad(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Info("Toggle the input pad.")
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}
	if err := pressWidget(wf, p, "Pad"); err != nil {
		return err
	}

	return wf.Execute()
}

// InputPhantom toggles the state of the input phantom button.
func InputPhantom(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Info("Toggle the input phantom.")
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}
	if err := pressWidget(wf, p, "Phantom"); err != nil {
		return err
	}

	return wf.Execute()
}

// InputSolo toggles the state of the input solo button.
func InputSolo(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Info("Toggle the input solo.")
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}
	if err := pressWidget(wf, p, "Solo"); err != nil {
		return err
	}

	return wf.Execute()
}

// SelectOutput for adjustment.
func SelectOutput(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Selecting %s %d output.", pkt.Signal, pkt.SignalNo)
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())
	if err := selectOutput(v, wf, pkt); err != nil {
		return err
	}
	return wf.Execute()
}

func selectOutput(v *Venue, wf *vnc.Workflow, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	ctrlName := signalControlName(pkt.Signal, pkt.SignalNo)
	if ctrlName == "Invalid" {
		return venuelib.Errorf(codes.Internal, "invalid control name for %s %d signal combination", pkt.Signal, pkt.SignalNo)
	}

	// Select the OUTPUTS page.
	p, err := v.ui.selectPage(wf, pages.Outputs)
	if err != nil {
		return err
	}

	// Select OUTPUTS tab (i.e. not USER).
	// TODO(kward:20170226) Handle this with UI element.
	wf.MouseClick(buttons.Left, image.Point{932, 524})

	// Clear the output solo.
	if glog.V(2) {
		glog.Infof("Clearing output solo.")
	}
	w, err := p.Widget("SoloClear")
	if err != nil {
		return err
	}
	if err := w.Press(wf); err != nil {
		return err
	}

	// Solo output.
	if glog.V(2) {
		glog.Infof("Soloing %q output.", ctrlName)
	}
	w, err = p.Widget(ctrlName + " Solo")
	if err != nil {
		return err
	}
	if err := w.Press(wf); err != nil {
		return err
	}

	return nil
}

// OutputLevel for the specified output. This handler operates on the
// currently selected input.
func OutputLevel(ep router.Endpoint, pkt *router.Packet) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Adjusting %s %d output level by %d dB.", pkt.Signal, pkt.SignalNo, pkt.Value)
	}

	ctrlName := signalControlName(pkt.Signal, pkt.SignalNo)
	if ctrlName == "Invalid" {
		return venuelib.Errorf(codes.InvalidArgument, "invalid control name for %s %d signal combination", pkt.Signal, pkt.SignalNo)
	}

	v := ep.(*Venue)
	wf := vnc.NewWorkflow(v.vnc.ClientConn())

	// Select output. Needed to select correct Aux or VarGroup.
	if err := selectOutput(v, wf, pkt); err != nil {
		return err
	}

	// Select the INPUTS page.
	p, err := v.ui.selectPage(wf, pages.Inputs)
	if err != nil {
		return err
	}

	// Adjust the Aux/Group knob.
	w, err := p.Widget(ctrlName)
	if err != nil {
		return err
	}
	if err := w.(*Encoder).Adjust(wf, pkt.Value.(int)); err != nil {
		return err
	}

	return wf.Execute()
}

//-----------------------------------------------------------------------------
// Misc

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

func pressWidget(wf *vnc.Workflow, page *Page, widget string) error {
	w, err := page.Widget(widget)
	if err != nil {
		return err
	}
	if err := w.Press(wf); err != nil {
		return err
	}
	return nil
}
