package touchosc

import (
	"flag"
	"testing"

	"github.com/kward/go-osc/osc"
	"github.com/kward/venue/router"
	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/router/controls"
	"github.com/kward/venue/router/signals"
)

func init() {
	if testing.Verbose() {
		flag.Set("alsologtostderr", "true")
		flag.Set("v", "5")
	}
}

func TestV01Parse(t *testing.T) {
	for _, tt := range []parseTest{
		{"thGain (press)",
			osc.NewMessage("/venue/0.1/th/soundcheck/input/gain/4/1", 1),
			&router.Packet{
				Source:  TouchOSC,
				Action:  actions.InputGain,
				Control: controls.Gain,
				Signal:  signals.Input,
				Value:   5,
			},
			true},
		{"thGain (release)",
			osc.NewMessage("/venue/0.1/th/soundcheck/input/gain/4/1", 0),
			&router.Packet{
				Source: TouchOSC,
				Action: actions.DropPacket,
			},
			true},
	} {
		pkt, err := Parse(tt.msg)
		if err != nil && tt.ok {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
		}
		if err == nil && !tt.ok {
			t.Errorf("%s: expected error", tt.name)
		}
		if !tt.ok {
			continue
		}
		if !pkt.Equal(tt.pkt) {
			t.Errorf("%s: packets not equal; got = %v, want = %v", tt.name, pkt, tt.pkt)
			continue
		}
	}
}

func TestV01PackGain(t *testing.T) {
	for _, tt := range []struct {
		desc   string
		layout string
		x, y   int
		want   int
		ok     bool
	}{
		{"th-2,1", "th", 2, 1, -1, true},
		{"ph-2,1", "pv", 2, 1, -1, true},
		{"th-3,1", "th", 3, 1, 1, true},
		{"pv-3,1", "pv", 3, 1, 1, true},
		{"bad", "th", 0, 0, 0, false},
	} {
		p := &packerV01{
			req: &request{
				msg:    osc.NewMessage("/test/inputGain", 1),
				layout: tt.layout,
				x:      tt.x,
				y:      tt.y,
			},
			pkt: &router.Packet{
				Signal: signals.Input,
			},
		}
		t.Logf("%s: request: %s", tt.desc, p.req)
		t.Logf("%s: packet before: %s", tt.desc, p.pkt)
		p.inputGain()
		if p.err != nil && tt.ok {
			t.Errorf("%s: unexpected error: %s", tt.desc, p.err)
		}
		if p.err == nil && !tt.ok {
			t.Errorf("%s: expected an error", tt.desc)
		}
		if !tt.ok {
			continue
		}
		if !p.done() {
			t.Errorf("%s: expected packing to be done", tt.desc)
		}
		t.Logf("%s: packet after: %s", tt.desc, p.pkt)

		if got, want := p.pkt.Action, actions.InputGain; got != want {
			t.Errorf("%s: packGain() y = %d: pkt.Action = %v, want = %v", tt.desc, tt.y, got, want)
		}
		if got, want := p.pkt.Value, tt.want; got != want {
			t.Errorf("%s: packGain() y = %d: pkt.Value = %d, want = %d", tt.desc, tt.y, got, want)
		}
	}
}

func TestVenueAuxGroup(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		y     int
		sig   signals.Signal
		sigNo int
	}{
		{"aux1", 1, signals.Aux, 1},
		{"aux9", 5, signals.Aux, 9},
		{"grp1", 9, signals.Group, 1},
	} {
		req := &request{y: tt.y}
		sig, sigNo := venueAuxGroup(req)
		if got, want := sig, tt.sig; got != want {
			t.Errorf("%s: sig: got %d, want %d", tt.desc, got, want)
		}
		if got, want := sigNo, tt.sigNo; got != want {
			t.Errorf("%s: sigNo: got %d, want %d", tt.desc, got, want)
		}
	}
}
