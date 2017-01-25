package oscparse

import (
	"flag"
	"testing"

	"github.com/kward/go-osc/osc"
)

func init() {
	if testing.Verbose() {
		flag.Set("alsologtostderr", "true")
		flag.Set("v", "5")
	}
}

func TestV01Parse(t *testing.T) {
	for _, tt := range []parseTest{
		{"thGain",
			"/venue/0.1/th/soundcheck/input/gain/4/1",
			[]interface{}{},
			&Packet{
				Ctrl: CtrlInput,
				Cmd:  CmdGain,
				Val:  5,
			},
			true},
	} {
		msg := osc.NewMessage(tt.addr)
		for _, v := range tt.val {
			msg.Append(v)
		}

		pkt, err := Parse(msg)
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
			req: request{
				layout: tt.layout,
				x:      tt.x,
				y:      tt.y,
			},
			pkt: &Packet{
				Ctrl: CtrlInput,
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
			t.Error("%s: expected packing to be done", tt.desc)
		}
		t.Logf("%s: packet after: %s", tt.desc, p.pkt)

		if got, want := p.pkt.Cmd, CmdGain; got != want {
			t.Errorf("%s: packGain() y = %d: pkt.Cmd = %v, want = %v", tt.desc, tt.y, got, want)
		}
		if got, want := p.pkt.Val, tt.want; got != want {
			t.Errorf("%s: packGain() y = %d: pkt.Val = %d, want = %d", tt.desc, tt.y, got, want)
		}
	}
}

func TestVenueAuxGroup(t *testing.T) {
	for _, tt := range []struct {
		desc string
		y    int
		ctrl Ctrl
		pos  int
	}{
		{"aux1", 1, CtrlAux, 1},
		{"aux9", 5, CtrlAux, 9},
		{"grp1", 9, CtrlGroup, 1},
	} {
		req := request{y: tt.y}
		ctrl, pos := venueAuxGroup(req)
		if got, want := ctrl, tt.ctrl; got != want {
			t.Errorf("%s: control: got %s, want %s", tt.desc, got, want)
		}
		if got, want := pos, tt.pos; got != want {
			t.Errorf("%s: position: got %d, want %d", tt.desc, got, want)
		}
	}
}
