package router

import (
	"github.com/golang/glog"
	"github.com/kward/venue/venuelib"
)

// type Router interface {
// 	// Subscribe to router messages.
// 	Subscribe() <-chan *Packet
// 	// Unsubscribe from router messages.
// 	Unsubscribe() error
// 	// Process a packet.
// 	Process(in <-chan *Packet)
// }

// An Endpoint can handle routed packets.
type Endpoint interface {
	// EndpointName returns the name of the endpoint.
	EndpointName() string

	// Handle a packet.
	Handle(pkt *Packet)
}

type Router struct {
	endpoints []Endpoint
}

func (r *Router) RegisterEndpoint(e Endpoint) {
	r.endpoints = append(r.endpoints, e)
}

// Dispatch the packet to the endpoints.
func Dispatch(r *Router, pkt *Packet) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	// Ask every endpoint to handle the packet.
	for _, e := range r.endpoints {
		if glog.V(2) {
			glog.Infof("Dispatching packet to %q endpoint.", e.EndpointName())
		}
		e.Handle(pkt)
		if glog.V(2) {
			glog.Info("Dispatching complete.")
		}
	}
}
