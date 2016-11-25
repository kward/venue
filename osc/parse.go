/*
parse takes an OSC string, tokenizes it using the lexer, and then provides
functionality to each of the tokens. As different versions of Venue support
different functionality, the parser enables different versions to provide their
own custom functionality.
*/
package oscparse

import (
	"errors"
	"fmt"

	"github.com/kward/go-osc/osc"
	"github.com/kward/venue"
	"github.com/kward/venue/venuelib"
)

type request struct {
	msg     *osc.Message
	version string // VENUE OSC protocol version
	dev     int    // device type
	orient  int    // orientation
	page    string
	control string
	command string
	bank    int
	x, y    int
	label   bool
}

type devType int
type orientType int

type PackerI interface {
	init(req request)
	done() bool
	error() error
	errorf(format string, args ...interface{}) packerM
	packer() packerM
	setPacker(pack packerM)
	pack()
	packet() *venue.Packet
}
type packerM func() packerM
type packerT struct {
	PackerI
	client string        // The name of client.
	err    error         // An error message, if present.
	fn     packerM       // The next packer state to enter.
	req    request       // The request to pack.
	pkt    *venue.Packet // The packet to pack.
}

var packers = map[string]PackerI{
	"0.1": &packerV01{},
}

func parse(msg *osc.Message) (*venue.Packet, error) {
	req := request{msg: msg, x: -1, y: -1}
	l := lex("OSC", msg.Address)
Parsing:
	for {
		item := l.nextItem()
		fmt.Printf("item: %v\n", item)
		switch item.typ {
		case itemCommand:
			req.command = item.val
		case itemControl:
			req.control = item.val
		case itemLabel:
			req.label = true
		case itemLayout:
			dev, orient, err := layout(item.val)
			if err != nil {
				return nil, err
			}
			req.dev, req.orient = dev, orient
		case itemPage:
			req.page = item.val
		case itemPositionX:
			req.x = venuelib.ToInt(item.val)
		case itemPositionY:
			req.y = venuelib.ToInt(item.val)
		case itemVersion:
			req.version = item.val
		case itemError:
			return nil, errors.New(fmt.Sprintf("unable to parse item %v", item))
		case itemEOF:
			break Parsing
		}
	}

	packer, ok := packers[req.version]
	if !ok {
		return nil, errors.New(fmt.Sprintf("unable to pack version %s", req.version))
	}
	packer.init(req)
	for packer.setPacker(packer.packer()); !packer.done(); {
		packer.pack()
	}
	if err := packer.error(); err != nil {
		return nil, err
	}
	return packer.packet(), nil
}

const (
	devInvalid int = 0x01 << iota
	devPhone
	devTablet
)
const (
	orientInvalid int = 0x10 << iota
	orientHoriz
	orientVert
)

func layout(s string) (int, int, error) {
	switch s {
	case "pv":
		return devPhone, orientVert, nil
	case "th":
		return devTablet, orientHoriz, nil
	}
	return devInvalid, orientInvalid, errors.New(fmt.Sprintf("unable to parse layout %s", s))
}
