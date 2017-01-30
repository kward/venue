package router

type Router interface {
	// Subscribe to router messages.
	Subscribe() <-chan *Packet
	// Unsubscribe from router messages.
	Unsubscribe() error
	// Process a packet.
	Process(in <-chan *Packet)
}
