/*
parse takes an OSC string, tokenizes it using the lexer, and then provides
functionality to each of the tokens. As different versions of Venue support
different functionality, the parser enables different versions to provide their
own custom functionality.
*/
package touchosc

import (
	"github.com/golang/glog"
	"github.com/kward/go-osc/osc"
	"github.com/kward/venue/codes"
	"github.com/kward/venue/router"
	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/venuelib"
)

// Parse the OSC message `msg` and transform it into a Packet.
func Parse(msg *osc.Message) (*router.Packet, error) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	req := &request{msg: msg, x: -1, y: -1}

	// Lex the request into tokens.
	l := lex("OSC", msg.Address)

	// Process the tokens.
Processing:
	for {
		item := l.nextItem()
		if glog.V(5) {
			glog.Infof("parsing item: %v", item)
		}
		switch item.typ {
		case itemRequest:
			req.request = item.val
		case itemCommand:
			req.command = item.val
		case itemControl:
			req.control = item.val
		case itemLabel:
			req.label = true
		case itemLayout:
			req.layout = item.val
		case itemPage:
			req.page = item.val
		case itemPositionX:
			req.x = venuelib.ToInt(item.val)
		case itemPositionY:
			req.y = venuelib.ToInt(item.val)
		case itemVersion:
			req.version = item.val
		case itemError:
			return nil, venuelib.Errorf(codes.InvalidArgument, "unable to parse item %v", item)
		case itemEOF:
			break Processing
		}
	}

	// Check for supported requests.
	switch req.request {
	case "ping":
		return &router.Packet{
			SourceName: TouchOSC,
			SourceAddr: req.msg.Addr(),
			Action:     actions.Ping,
		}, nil
	case "venue":
	default:
		return nil, venuelib.Errorf(codes.InvalidArgument, "unrecognized request %q", req.request)
	}

	// Packetize the request.
	packer, ok := packers[req.version]
	if !ok {
		return nil, venuelib.Errorf(codes.NotFound, "unable to pack version %s", req.version)
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
