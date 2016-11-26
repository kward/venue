package venue

import "reflect"

// ctrlType identifies the type of control
type ctrlType int

const (
	CtrlUnknown ctrlType = iota
	CtrlInput
	CtrlOutput
)

type cmdType int

const (
	CmdUnknown cmdType = iota
	CmdBank
	CmdGain
	CmdSelect
)

// ValueType identifies what type of value a packet holds.
type ValueType int

const (
	ValUnknown ValueType = iota // The value is undefined.
	ValFloat
	ValInteger
	ValString
)

// Packet represents a Venue action to perform.
type Packet struct {
	Ctrl      ctrlType  // The control.
	Cmd       cmdType   // The command.
	Pos       int       // The position or channel number.
	Val       ValueType // The type of the value.
	FloatVal  float64
	IntVal    int
	StringVal string
}

func (p Packet) Equal(p2 Packet) bool {
	return reflect.DeepEqual(p, p2)
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
