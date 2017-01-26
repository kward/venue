package oscparse

import (
	"fmt"
	"reflect"
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

// Cmd identifies the type of command.
type Cmd int

const (
	UnknownCmd Cmd = iota
	BankCmd
	GainCmd
	PingCmd
	SetOutputLevelCmd
	SelectOutputCmd
	SelectInputCmd
)

var Cmds = map[Cmd]string{
	UnknownCmd:        "unknown",
	BankCmd:           "bank",
	GainCmd:           "gain",
	PingCmd:           "ping",
	SetOutputLevelCmd: "set_output_level",
	SelectOutputCmd:   "select_output",
	SelectInputCmd:    "select_input",
}

// Packet represents a Venue action to perform.
type Packet struct {
	Ctrl Ctrl // The control.
	Cmd  Cmd  // The command.
	Pos  int  // The position or channel number.
	Val  interface{}
}

// Equal returns true if the two packets are equal.
func (p *Packet) Equal(p2 *Packet) bool {
	return reflect.DeepEqual(p, p2)
}

// String returns a human readable representation of the packet.
func (p *Packet) String() string {
	return fmt.Sprintf("{ Ctrl: %s Cmd: %s Pos: %d Value: %v }",
		Ctrls[p.Ctrl], Cmds[p.Cmd], p.Pos, p.Val)
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
