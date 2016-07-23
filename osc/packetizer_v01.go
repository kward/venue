package oscparse

import (
	"errors"
	"fmt"

	"github.com/kward/venue"
)

type packerV01 struct {
	PackerI
	client string        // The name of client.
	err    error         // An error message, if present.
	fn     packerM       // The next packer state to enter.
	req    request       // The request to pack.
	pkt    *venue.Packet // The packet to pack.
	dev    packerV01Device
}

type packerV01Device struct {
	dxInput, dyInput int
	dxOutput         int
}

type packerV01ItemType int
type packerV01Item struct {
	typ packerV01ItemType
}

const (
	packerV01ItemError packerV01ItemType = iota // An error occurred, and is held in err.
	packerV01ItemInput
)

func (p *packerV01) packControl() packerM {
	switch p.req.control {
	case "input":
		return p.packInput
	}
	return p.errorf("invalid control")
}

func (p *packerV01) packInput() packerM {
	p.pkt.Ctrl = venue.CtrlInput
	switch p.req.command {
	case "bank":
		return p.packBank
	case "gain":
		return p.packGain
	case "select":
		return p.packSelect
	}
	return p.errorf("invalid input")
}

func (p *packerV01) packBank() packerM {
	return p.errorf("invalid bank")
}

// packGain translates the gain MultiPush HID element position into a VENUE gain
// UI up/down click count.
func (p *packerV01) packGain() packerM {
	p.pkt.Cmd = venue.CmdGain

	pos := 0
	switch p.req.dev | p.req.orient {
	case devTablet | orientHoriz:
		pos = p.req.y
	}

	clicks := 0
	switch pos {
	case 1:
		clicks = 5 // +5 dB
	case 2:
		clicks = 1 // +1 dB
	case 3:
		clicks = -1 // -1 dB
	case 4:
		clicks = -5 // -5 dB
	}
	//v := p.req.msg.Arguments[0]
	p.setValue(venue.ValInteger, clicks)

	return nil
}

func (p *packerV01) packChannel() packerM {
	fmt.Println(p.req.dev, p.req.orient)
	switch p.req.dev | p.req.orient {
	case devPhone | orientVert:
		p.dev.dxInput, p.dev.dyInput, p.dev.dxOutput = 8, 4, 6
	case devTablet | orientHoriz:
		p.dev.dxInput, p.dev.dyInput, p.dev.dxOutput = 12, 4, 12
	}
	// TODO: this is broken
	return nil
}

func (p *packerV01) packSelect() packerM {
	return p.errorf("invalid select")
}

func (p *packerV01) errorf(format string, args ...interface{}) packerM {
	fmt.Println("errorf()")
	p.err = errors.New(fmt.Sprintf(format, args...))
	return nil
}

func (p *packerV01) init(req request) {
	p.fn = p.packControl
	p.pkt = &venue.Packet{}
	p.req = req
}
func (p *packerV01) done() bool            { return p.fn == nil }
func (p *packerV01) error() error          { return p.err }
func (p *packerV01) packer() packerM       { return p.packControl }
func (p *packerV01) setPacker(fn packerM)  { p.fn = fn }
func (p *packerV01) pack()                 { p.fn = p.fn() }
func (p *packerV01) packet() *venue.Packet { return p.pkt }

func (p *packerV01) setValue(vt venue.ValueType, v interface{}) {
	p.pkt.Val = vt
	switch vt {
	case venue.ValFloat:
		p.pkt.FloatVal = v.(float64)
	case venue.ValInteger:
		p.pkt.IntVal = v.(int)
	case venue.ValString:
		p.pkt.StringVal = v.(string)
	}
}

// func NewPacker(req request) *Packer {
// 	p := &packerV01{pkt: &venue.Packet{}, req: req}
// 	p.fn = p.packControl
// 	return Packer(*p)
// }
