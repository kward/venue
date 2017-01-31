package touchosc

import (
	"github.com/kward/go-osc/osc"
	"github.com/kward/venue/router"
)

type parseTest struct {
	name string
	msg  *osc.Message
	pkt  *router.Packet
	ok   bool
}
