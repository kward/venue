package ping

import (
	"github.com/golang/glog"
	"github.com/kward/venue/internal/router"
	"github.com/kward/venue/internal/router/actions"
	"github.com/kward/venue/internal/venuelib"
	"github.com/kward/venue/touchosc"
)

// Ping is a Venue endpoint.
type Ping struct{}

// Verify that the expected interface is implemented properly.
var _ router.Endpoint = new(Ping)

// EndpointName implements router.Endpoint.
func (e *Ping) EndpointName() string { return "Ping" }

// Handle implements router.Endpoint.
func (e *Ping) Handle(pkt *router.Packet) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if pkt.SourceName != touchosc.TouchOSC && pkt.Action != actions.Ping {
		return
	}
	if glog.V(2) {
		glog.Infof("Received ping request from %s/%s.", pkt.SourceName, pkt.SourceAddr)
	}
}
