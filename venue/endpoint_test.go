package venue

import (
	"testing"

	"github.com/kward/venue/codes"
	"github.com/kward/venue/router"
	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/router/signals"
	"github.com/kward/venue/venuelib"
	"github.com/kward/venue/vnc"
)

func TestSignalControlName(t *testing.T) {
	for _, tt := range []struct {
		sig   signals.Signal
		sigNo signals.SignalNo
		name  string
	}{
		{signals.Aux, 1, "Aux 1"},
		{signals.Group, 2, "Group 2"},
		{signals.Direct, 3, "Invalid"},
	} {
		if got, want := signalControlName(tt.sig, tt.sigNo), tt.name; got != want {
			t.Errorf("signalControlName(%s, %d) = %s, want %s", tt.sig, tt.sigNo, got, want)
			continue
		}
	}
}

func TestSignalControlName_MoreMappings(t *testing.T) {
	for _, tt := range []struct {
		sig   signals.Signal
		sigNo signals.SignalNo
		name  string
	}{
		{signals.Input, 5, "Fader"},
		{signals.FXReturn, 9, "Fader"},
		{signals.Unknown, 0, "Unknown"},
		{signals.Matrix, 1, "Unknown"}, // default branch
		{signals.Aux, 96, "Aux 96"},
		{signals.Group, 12, "Group 12"},
	} {
		if got, want := signalControlName(tt.sig, tt.sigNo), tt.name; got != want {
			t.Errorf("signalControlName(%s, %d) = %q, want %q", tt.sig, tt.sigNo, got, want)
		}
	}
}

func TestSelectInput_TooManyInputs_ReturnsInvalidArgument(t *testing.T) {
	v := &Venue{}
	pkt := &router.Packet{SignalNo: signals.SignalNo(maxInputs + 1)}
	err := SelectInput(v, pkt)
	if err == nil {
		t.Fatalf("expected error for too many inputs, got nil")
	}
	if code := venuelib.Code(err); code != codes.InvalidArgument {
		t.Fatalf("unexpected error code: got %v, want %v (err=%v)", code, codes.InvalidArgument, err)
	}
}

func TestOutputLevel_InvalidSignal_ReturnsInvalidArgument(t *testing.T) {
	v := &Venue{}
	pkt := &router.Packet{Signal: signals.Direct, SignalNo: 1, Value: 1}
	err := OutputLevel(v, pkt)
	if err == nil {
		t.Fatalf("expected error for invalid signal, got nil")
	}
	if code := venuelib.Code(err); code != codes.InvalidArgument {
		t.Fatalf("unexpected error code: got %v, want %v (err=%v)", code, codes.InvalidArgument, err)
	}
}

func TestOutputLevel_Table(t *testing.T) {
	// Build a minimally-initialized Venue to avoid panics and return predictable errors.
	baseVenue := func() *Venue {
		return &Venue{
			ui:  NewUI(),    // UI available so page/widget lookups work
			vnc: &vnc.VNC{}, // zero-value VNC so Workflow.Execute() returns Internal
		}
	}

	tests := []struct {
		name  string
		sig   signals.Signal
		sigNo signals.SignalNo
		value int
		want  codes.Code
	}{
		{name: "Direct invalid", sig: signals.Direct, sigNo: 1, value: 1, want: codes.InvalidArgument},
		{name: "Aux valid path (no conn)", sig: signals.Aux, sigNo: 1, value: 2, want: codes.Internal},
		{name: "Group valid path (no conn)", sig: signals.Group, sigNo: 3, value: -1, want: codes.Internal},
		{name: "Unknown maps to Internal", sig: signals.Unknown, sigNo: 1, value: 1, want: codes.Internal},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := baseVenue()
			pkt := &router.Packet{Signal: tt.sig, SignalNo: tt.sigNo, Value: tt.value}
			err := OutputLevel(v, pkt)
			if err == nil {
				t.Fatalf("OutputLevel() error = nil, want code %v", tt.want)
			}
			if got := venuelib.Code(err); got != tt.want {
				t.Fatalf("OutputLevel() code = %v, want %v; err=%v", got, tt.want, err)
			}
		})
	}
}

func TestEndpointName(t *testing.T) {
	v, err := New()
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}
	if got, want := v.EndpointName(), "Venue"; got != want {
		t.Fatalf("EndpointName() = %q, want %q", got, want)
	}
}

func TestNoopHandler(t *testing.T) {
	if err := Noop(&Venue{}, &router.Packet{Action: actions.Noop}); err != nil {
		t.Fatalf("Noop() error: %v", err)
	}
}

func TestHandlersInitialized(t *testing.T) {
	// Ensure expected actions are present and have non-nil handlers.
	required := []actions.Action{
		actions.Noop,
		actions.Ping,
		actions.SelectInput,
		actions.InputGain,
		actions.InputMute,
		actions.InputPad,
		actions.InputPhantom,
		actions.InputSolo,
		actions.SelectOutput,
		actions.OutputLevel,
	}
	for _, a := range required {
		spec, ok := handlers[a]
		if !ok {
			t.Fatalf("handler for %s not initialized", a)
		}
		if spec.Handler == nil {
			t.Fatalf("handler for %s is nil", a)
		}
	}
}
