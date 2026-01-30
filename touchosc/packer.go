package touchosc

import "github.com/kward/venue/internal/router"

const TouchOSC = "TouchOSC"

type Packer interface {
	// init prepares the packer to parse the request `req`.
	init(req *request)
	// done returns true when packing is complete.
	done() bool

	// error returns the packet error.
	error() error
	// errorf stores a formatted error for later recovery.
	errorf(format string, args ...interface{}) packerFn

	// packer returns the current packer function.
	packer() packerFn
	// setPacker sets the next packer function.
	setPacker(pack packerFn)
	// pack calls the current packer function.
	pack()

	// packet returns the constructed packet.
	packet() *router.Packet
}

type packerFn func() packerFn

type packerT struct {
	err error          // An error message, if present.
	fn  packerFn       // The next packer state to enter.
	req *request       // The request to pack.
	pkt *router.Packet // The packet to pack.
}

var (
	packers = map[string]Packer{
		"0.1": &packerV01{},
	}
)
