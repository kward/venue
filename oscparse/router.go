package oscparse

import (
	"fmt"
	"reflect"

	"github.com/kward/venue/oscparse/commands"
	"github.com/kward/venue/oscparse/controls"
)

// Packet represents a Venue action to perform.
type Packet struct {
	Control controls.Control // The control.
	Command commands.Command // The command.
	Pos     int              // The position or channel number.
	Val     interface{}
}

// Equal returns true if the two packets are equal.
func (p *Packet) Equal(p2 *Packet) bool {
	return reflect.DeepEqual(p, p2)
}

// String returns a human readable representation of the packet.
func (p *Packet) String() string {
	return fmt.Sprintf("{ Control: %s Command: %s Pos: %d Value: %v }",
		p.Control, p.Command, p.Pos, p.Val)
}

type packetBus struct {
	in  chan Packet   // Incoming packet.
	out []chan Packet // Slice of packet listeners.
}

type Router interface {
	// Subscribe to router messages.
	Subscribe() <-chan *Packet
	// Unsubscribe from router messages.
	Unsubscribe() error
	// Process a packet.
	Process(in <-chan *Packet)
}
