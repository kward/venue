// TODO(kward:20170122) The packetizer should be separate from the lexer/parser.
package touchosc

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kward/venue/internal/router"
	"github.com/kward/venue/internal/router/actions"
	"github.com/kward/venue/internal/router/controls"
	"github.com/kward/venue/internal/router/signals"
	"github.com/kward/venue/internal/venuelib"
	"github.com/kward/venue/touchosc/multistates"
)

type packerV01 packerT

// Verify that the expected interface is implemented properly.
var _ Packer = new(packerV01)

func (p *packerV01) init(req *request) {
	p.err = nil
	p.fn = p.packByControl
	p.req = req
	p.pkt = &router.Packet{}
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

func (p *packerV01) packet() *router.Packet { return p.pkt }

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

	p.pkt.Signal = signals.Input
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
	case "phantom":
		return p.inputPhantom
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

	args := p.req.msg.Arguments
	if len(args) == 0 {
		return p.errorf("missing OCS arguments")
	}
	switch multistates.State(args[0]) {
	case multistates.Released: // Do nothing.
		p.setPacket(router.NewNoopPacket())
		return nil
	case multistates.Unknown:
		return p.errorf("received invalid argument %v", args[0])
	}

	clicks := clicks(p.req.x)
	if clicks == 0 {
		return p.errorf("invalid gain control x/y: %d/%d", p.req.x, p.req.y)
	}

	p.setPacket(&router.Packet{
		Action:  actions.InputGain,
		Control: controls.Gain,
		Signal:  signals.Input,
		// No SignalNo as we expect to work on the currently selected channel.
		Value: clicks,
	})
	return nil
}

func (p *packerV01) inputGuess() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	args := p.req.msg.Arguments
	if len(args) != 1 {
		return p.errorf("expected 1 argument, got %d; %v", len(args), args)
	}
	switch multistates.State(args[0]) {
	case multistates.Released: // Do nothing.
		p.setPacket(router.NewNoopPacket())
		return nil
	case multistates.Unknown:
		return p.errorf("received invalid argument %v", args[0])
	}

	p.setPacket(&router.Packet{
		Action:  actions.InputGuess,
		Control: controls.Guess,
		Signal:  signals.Input,
	})
	return nil
}

func (p *packerV01) inputMute() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	return p.multiToggle(&router.Packet{
		Action:  actions.InputMute,
		Control: controls.Mute,
		Signal:  signals.Input,
	})
}

func (p *packerV01) inputPad() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	return p.multiToggle(&router.Packet{
		Action:  actions.InputPad,
		Control: controls.Pad,
		Signal:  signals.Input,
	})
}

const (
	dxInputSelect = 4  // Multi-Push/-Toggle Y value.
	dyInputSelect = 12 // Multi-Push/-Toggle Y value.
)

func (p *packerV01) inputSelect() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	args := p.req.msg.Arguments
	if len(args) == 0 {
		return p.errorf("missing OCS arguments")
	}
	switch multistates.State(args[0]) {
	case multistates.Released: // Do nothing.
		p.setPacket(router.NewNoopPacket())
		return nil
	case multistates.Unknown:
		return p.errorf("invalid OSC argument %v", args[0])
	}

	pos := p.req.multiPosition(dxInputSelect, dyInputSelect)
	p.setPacket(&router.Packet{
		Action:   actions.SelectInput,
		Control:  controls.Select,
		Signal:   signals.Input,
		SignalNo: (signals.SignalNo)(pos),
	})
	return nil
}

func (p *packerV01) inputPhantom() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	return p.multiToggle(&router.Packet{
		Action:  actions.InputPhantom,
		Control: controls.Phantom,
		Signal:  signals.Input,
	})
}

func (p *packerV01) inputSolo() packerFn {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	return p.multiToggle(&router.Packet{
		Action:  actions.InputSolo,
		Control: controls.Solo,
		Signal:  signals.Input,
	})
}

func (p *packerV01) multiToggle(pkt *router.Packet) packerFn {
	args := p.req.msg.Arguments
	if len(args) != 1 {
		return p.errorf("expected 1 argument, got %d; %v", len(args), args)
	}
	state := multistates.State(args[0])
	if state == multistates.Unknown {
		return p.errorf("received invalid argument %v", args[0])
	}
	pkt.Value = state
	p.setPacket(pkt)
	return nil
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

	args := p.req.msg.Arguments
	if len(args) == 0 {
		return p.errorf("missing OCS arguments")
	}
	switch multistates.State(args[0]) {
	case multistates.Released: // Do nothing.
		p.setPacket(router.NewNoopPacket())
		return nil
	case multistates.Unknown:
		return p.errorf("received invalid argument %v", args[0])
	}

	sig, sigNo := venueAuxGroup(p.req)
	clicks := clicks(p.req.x)
	if clicks == 0 {
		return p.errorf("invalid level control x/y: %d/%d", p.req.x, p.req.y)
	}
	p.setPacket(&router.Packet{
		Action:   actions.OutputLevel,
		Signal:   sig,
		SignalNo: sigNo,
		Value:    clicks,
	})
	if glog.V(4) {
		glog.Infof("packet: %s", p.pkt)
	}
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

	args := p.req.msg.Arguments
	if len(args) == 0 {
		return p.errorf("missing OCS arguments")
	}
	switch multistates.State(args[0]) {
	case multistates.Released: // Do nothing.
		p.setPacket(router.NewNoopPacket())
		return nil
	case multistates.Unknown:
		return p.errorf("received invalid argument %v", args[0])
	}

	sig, sigNo := venueAuxGroup(p.req)
	p.setPacket(&router.Packet{
		Action:   actions.SelectOutput,
		Signal:   sig,
		SignalNo: sigNo,
	})
	return nil
}

func (p *packerV01) setPacket(pkt *router.Packet) {
	p.pkt = pkt
	if pkt == nil {
		return
	}
	if p.pkt.SourceName == "" {
		p.pkt.SourceName = TouchOSC
	}
	if p.pkt.SourceAddr == "" {
		p.pkt.SourceAddr = p.req.msg.Addr()
	}
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
func venueAuxGroup(req *request) (signals.Signal, signals.SignalNo) {
	sig := signals.Aux
	sigNo := req.y
	if req.y > 8 {
		sig = signals.Group
		sigNo = req.y - 8
	}
	sigNo = sigNo*2 - 1 // Convert position into stereo channel number.
	return sig, (signals.SignalNo)(sigNo)
}
