package oscparse

import (
	"testing"

	"github.com/kward/go-osc/osc"
	"github.com/kward/venue"
)

func TestV01Parse(t *testing.T) {
	for _, tt := range []parseTest{
		{"thGain",
			"/venue/0.1/th/soundcheck/input/gain/1/3",
			[]interface{}{},
			&venue.Packet{
				Ctrl:   venue.CtrlInput,
				Cmd:    venue.CmdGain,
				Val:    venue.ValInteger,
				IntVal: -1,
			},
			true},
	} {
		msg := osc.NewMessage(tt.addr)
		for _, v := range tt.val {
			msg.Append(v)
		}

		pkt, err := parse(msg)
		if err != nil && tt.ok {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
		}
		if err == nil && !tt.ok {
			t.Errorf("%s: expected error", tt.name)
		}
		if !tt.ok {
			continue
		}
		if !pkt.Equal(*tt.pkt) {
			t.Errorf("%s: packets not equal; got = %v, want = %v", tt.name, pkt, tt.pkt)
		}
	}
}

func TestV01PackGain(t *testing.T) {
	for _, tt := range []struct {
		y    int
		want int
	}{
		{0, 0},
		{1, 5},
		{2, 1},
		{3, -1},
		{4, -5},
	} {
		p := &packerV01{
			req: request{
				dev:    devTablet,
				orient: orientHoriz,
				y:      tt.y,
			},
			pkt: &venue.Packet{
				Ctrl: venue.CtrlInput,
			},
		}

		fn := p.packGain()
		if fn != nil {
			t.Errorf("expected nil fn, got = %v", fn)
		}
		if got, want := p.pkt.Cmd, venue.CmdGain; got != want {
			t.Errorf("packGain() y = %d: pkt.Cmd = %v, want = %v", tt.y, got, want)
		}
		if got, want := p.pkt.Val, venue.ValInteger; got != want {
			t.Errorf("packGain() y = %d: pkt.Val = %d, want = %d", tt.y, got, want)
		}
		if got, want := p.pkt.IntVal, tt.want; got != want {
			t.Errorf("packGain() y = %d: pkt.IntVal = %d, want = %d", tt.y, got, want)
		}
		t.Log(p, p.pkt)
	}
}
