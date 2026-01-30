package router

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
