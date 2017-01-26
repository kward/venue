package oscparse

import (
	"fmt"
	"reflect"

	"github.com/kward/venue/oscparse/commands"
)

// Ctrl identifies the type of control.
type Ctrl int

const (
	UnknownCtrl Ctrl = iota
	InputCtrl
	AuxCtrl
	GroupCtrl
)

var Ctrls = map[Ctrl]string{
	UnknownCtrl: "unknown",
	InputCtrl:   "input",
	AuxCtrl:     "aux",
	GroupCtrl:   "grp",
}

// Packet represents a Venue action to perform.
type Packet struct {
	Ctrl    Ctrl             // The control.
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
	return fmt.Sprintf("{ Ctrl: %s Command: %s Pos: %d Value: %v }",
		Ctrls[p.Ctrl], p.Command, p.Pos, p.Val)
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
