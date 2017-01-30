/*
parse takes an OSC string, tokenizes it using the lexer, and then provides
functionality to each of the tokens. As different versions of Venue support
different functionality, the parser enables different versions to provide their
own custom functionality.
*/
package touchosc

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/kward/go-osc/osc"
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
		case itemPingReq, itemVenueReq:
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
			return nil, fmt.Errorf("unable to parse item %v", item)
		case itemEOF:
			break Processing
		}
	}

	// Check for supported requests.
	switch req.request {
	case PingReq:
		return &router.Packet{Action: actions.Ping}, nil
	case VenueReq:
	default:
		return nil, fmt.Errorf("unrecognized request %q", req.request)
	}

	// Packetize the request.
	packer, ok := packers[req.version]
	if !ok {
		return nil, fmt.Errorf("unable to pack version %s", req.version)
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

// request holds the raw OSC message, and its lexed equivalent.
type request struct {
	msg *osc.Message
	// Lexed values.
	request string
	version string
	layout  string
	page    string
	control string
	command string
	bank    int
	x, y    int
	label   bool
}

// String returns a human readable representation of the request.
func (req request) String() string {
	return fmt.Sprintf("{ msg: %v request: %q version: %q layout: %q page: %q control: %q command: %q bank: %d x: %d y: %d label: %v",
		req.msg, req.request, req.version, req.layout, req.page, req.control, req.command, req.bank, req.x, req.y, req.label)
}

// isTablet returns true for tablet client requests.
func (req *request) isTablet() bool {
	switch req.layout {
	case "th", "tv":
		return true
	case "ph", "pv":
		return false
	}
	glog.Errorf("invalid request: %s", req)
	return false
}

// isHorizontal returns true for horizontally oriented client requests.
func (req *request) isHorizontal() bool {
	switch req.layout {
	case "th", "ph":
		return true
	case "tv", "pv":
		return false
	}
	glog.Errorf("invalid request: %s", req)
	return false
}

// The multi* UI controls report their x and y position as /x/y.
// In a vertical orientation, x and y correspond to the top-left of the control,
// with x increasing to the right and y increasing downwards.
// In a horizontal orientation, x and y correspond to the bottom-left of the
// control, with x increasing vertically and y increasing to the right.
//
// Vertical: 1, 1 is top-left, X inc right, Y inc down
// | 1 2 3 |
// | 2 2 3 |
// | 3 3 3 |dy=3
//       dx=3
//
// Horizontal: 1, 1 is bottom-left, X inc up, Y inc right
// | 3 3 3 |
// | 2 2 3 |
// | 1 2 3 |

// multiPosition returns the absolute position on a Multi-* control.
//
// Assuming a 3x2 Multi-Push/-Toggle control, with coordinates mapped to the
// vertical orientation (see multiRotate()), the "absolute position" within the
// control can be seen as this:
// | 1/1 2/1 3/1 | --> | 1 2 3 |
// | 1/2 2/2 3/2 |     | 4 5 6 |
//
// This is useful when one wants to turn a large XxY control into a single
// value.
//
// `x` and `y` correspond to the parsed OSC values for a control, and `dx`
// is the number of controls on the X axis.
func (req *request) multiPosition(dx, dy int) int {
	if req.isHorizontal() {
		x, y := req.multiRotate(dx)
		return x + dy*(y-1)
	}
	return req.x + dx*(req.y-1)
}

// multiRotate returns the rotated position for a Multi-* control.
//
// When a 3x2 Multi-Push/-Toggle control is drawn in vertical orientation, its
// x/y coordinates look like this:
// | 1/1 2/1 3/1 |
// | 1/2 2/2 3/2 |
//
// Drawing that same 3x2 control in the horizontal orientation, the x/y
// coordinates rotate with the control, which looks like this:
// | 3/1 3/2 |
// | 2/1 2,2 |
// | 1/1 1/2 |
//
// Although the orientation may change, OSC clients will transmit the same
// coordinates for either orientation. From a code perspective though, it can be
// nice to always refer the upper-left coordinate as 1/1, as though it were in
// the vertical rotation. This function maps a horizontal coordinate into its
// vertical equivalent. Therefore, the coordinates translate like this:
// | 3/1 3/2 |     | 1/1 2/1 |
// | 2/1 2,2 | --> | 1/2 2/2 |
// | 1/1 1/2 |     | 1/3 2/3 |
// horizontal      vertical-ref
//
// `x` and `y` correspond to the parsed OSC values for a control, and `dx`
// is the number of controls on the X axis.
func (req *request) multiRotate(dx int) (int, int) {
	if req.isHorizontal() {
		return req.y, dx - req.x + 1
	}
	return req.x, req.y
}
