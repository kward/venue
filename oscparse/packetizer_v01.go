// TODO(kward:20170122) The packetizer should be separate from the lexer/parser.
package oscparse

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kward/venue/oscparse/commands"
	"github.com/kward/venue/oscparse/controls"
	"github.com/kward/venue/venuelib"
)

type packerV01 packerT

// Verify that the PackerI interface is honored.
var _ PackerI = new(packerV01)

func (p *packerV01) init(req request) {
	p.setPacker(p.packByControl)
	p.pkt = &Packet{}
	p.req = req
}
func (p *packerV01) done() bool { return p.fn == nil }

func (p *packerV01) error() error { return p.err }
func (p *packerV01) errorf(format string, args ...interface{}) packerFn {
	p.err = fmt.Errorf(format, args...)
	glog.Errorf("packer error: %s", p.err)
	return nil
}

func (p *packerV01) packer() packerFn      { return p.fn }
func (p *packerV01) setPacker(fn packerFn) { p.fn = fn }
func (p *packerV01) pack()                 { p.fn = p.fn() }

func (p *packerV01) packet() *Packet { return p.pkt }

func (p *packerV01) packByControl() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Packing control %q.", p.req.control)
	}
	switch p.req.control {
	case "input":
		return p.input
	case "output":
		return p.output
	default:
		return p.errorf("invalid control %q", p.req.control)
	}
}

//-----------------------------------------------------------------------------
// Input control.

func (p *packerV01) input() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Packing input command %q.", p.req.command)
	}

	p.pkt.Control = controls.Input
	switch p.req.command {
	case "bank":
		return p.inputBank
	case "gain":
		return p.inputGain
	case "guess":
		return p.inputGuess
	case "mute":
		return p.inputMute
	case "pad":
		return p.inputPad
	case "select":
		return p.inputSelect
	case "solo":
		return p.inputSolo
	default:
		return p.errorf("invalid input %q", p.req.command)
	}
}

func (p *packerV01) inputBank() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	return p.errorf("%s unimplemented", venuelib.FnName())
}

// inputGain translates the gain MultiPush HID element position into a VENUE
// gain UI up/down click count.
//
// The gain control is a Multi-XY widget. On a horizontal layout, X/Y is the
// bottom-left, with X increasing vertically and Y increasing horizontally.
func (p *packerV01) inputGain() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	clicks := clicks(p.req.x)
	if clicks == 0 {
		return p.errorf("invalid gain control x/y: %d/%d", p.req.x, p.req.y)
	}

	p.pkt.Command = commands.InputGain
	p.pkt.Val = clicks
	return nil
}

func (p *packerV01) inputGuess() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	return p.errorf("%s unimplemented", venuelib.FnName())
}

func (p *packerV01) inputMute() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	return p.errorf("%s unimplemented", venuelib.FnName())
}

func (p *packerV01) inputPad() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	return p.errorf("%s unimplemented", venuelib.FnName())
}

const (
	dxInputSelect = 4  // Multi-Push/-Toggle Y value.
	dyInputSelect = 12 // Multi-Push/-Toggle Y value.
)

func (p *packerV01) inputSelect() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	pos := p.req.multiPosition(dxInputSelect, dyInputSelect)
	p.pkt = &Packet{Control: controls.Input, Command: commands.SelectInput, Pos: pos}
	return nil
}

func (p *packerV01) inputSolo() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	return p.errorf("%s unimplemented", venuelib.FnName())
}

//-----------------------------------------------------------------------------
// Output control.

func (p *packerV01) output() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Packing output command %q.", p.req.command)
	}

	switch p.req.command {
	case "level":
		return p.outputLevel
	case "pan":
		return p.outputPan
	case "select":
		return p.outputSelect
	default:
		return p.errorf("invalid input %q", p.req.command)
	}
}

func (p *packerV01) outputLevel() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	ctrl, pos := venueAuxGroup(p.req)
	clicks := clicks(p.req.x)
	if clicks == 0 {
		return p.errorf("invalid level control x/y: %d/%d", p.req.x, p.req.y)
	}

	p.pkt = &Packet{Control: ctrl, Command: commands.OutputLevel, Pos: pos, Val: clicks}
	return nil
}

func (p *packerV01) outputPan() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	return p.errorf("%s unimplemented", venuelib.FnName())
}

func (p *packerV01) outputSelect() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	ctrl, pos := venueAuxGroup(p.req)
	p.pkt = &Packet{Control: ctrl, Command: commands.SelectOutput, Pos: pos}
	return nil
}

//-----------------------------------------------------------------------------
// Miscellaneous.

// clicks converts the X value of 4x1 (XxY) multi UI control into a count of
// mouse clicks, which in Venue equates to a dB value in-/decrease. A value of
// 0 is an error.
func clicks(x int) int {
	switch x {
	case 4:
		return 5
	case 3:
		return 1
	case 2:
		return -1
	case 1:
		return -5
	default:
		return 0
	}
}

// venueAugGroup converts request into a Control and position.
// Note: a Bus Configuration of "16 Auxes + 8 Variable Groups (24 bus)" is
// assumed.
func venueAuxGroup(req request) (controls.Control, int) {
	ctrl := controls.Aux
	pos := req.y
	if req.y > 8 {
		ctrl = controls.Group
		pos = req.y - 8
	}
	pos = pos*2 - 1 // Convert position into stereo channel number.
	return ctrl, pos
}
