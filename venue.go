/*
The Venue package exposes the Avid™ VENUE VNC interface as a programmatic API.
*/
package venue

import (
	"context"
	"log"
	"time"

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

	inputs    [numInputs]*Input
	currInput *Input

	outputs    map[string]*Output
	currOutput *Output

	Pages    VenuePages
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
func (v *Venue) Initialize() {
	// Initialize pages.
	v.Pages = VenuePages{}
	v.Pages[InputsPage] = NewInputsPage()
	v.Pages[OutputsPage] = NewOutputsPage()
	// Initialize inputs.
	for ch := 0; ch < numInputs; ch++ {
		input := NewInput(v, ch+1, Ichannel)
		v.inputs[ch] = input
	}

	// Choose output before input so that later when the Inputs page is selected,
	// it shows first bank of channels.
	v.SetOutput("aux1")
	v.SetInput(1)

	// Clear solo.
	log.Println("Clearing solo.")
	vp := v.Pages[InputsPage]
	e := vp.Elements["solo_clear"]
	e.(*Switch).Update(v)
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

// abs returns the absolute value of an int.
func abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0
	}
	return x
}
