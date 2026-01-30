// Package router provides endpoint packet routing functionality.
package router

import (
	"github.com/golang/glog"
	"github.com/kward/venue/codes"
	"github.com/kward/venue/internal/router/actions"
	"github.com/kward/venue/venuelib"
)

// Handler describes an Endpoint handler.
type Handler func(ep Endpoint, pkt *Packet) error

// Handlers is a map of HandlerSpec keyed on actions.Action.
type Handlers map[actions.Action]HandlerSpec

// HandlerSpec represents an Endpoint handler specification.
type HandlerSpec struct {
	Action  actions.Action
	Handler Handler
}

// Handle requests that a packet be handled by the appropriate handler from
// Handlers `hs` for the Endpoint `ep`.
func Handle(ep Endpoint, pkt *Packet, hs Handlers) error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	if glog.V(2) {
		glog.Infof("Handling %s packet.", pkt.Action)
	}

	h, ok := hs[pkt.Action]
	if !ok {
		return venuelib.Errorf(codes.Unimplemented, "%s action unimplemented for %s.", pkt.Action, ep.EndpointName())
	}
	return h.Handler(ep, pkt)
}

// Router holds the representation of a Router.
type Router struct {
	endpoints []Endpoint
}

// RegisterEndpoint to make the endpoint available for accepting requests.
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
