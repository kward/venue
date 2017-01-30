package router

import (
	"fmt"
	"reflect"

	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/router/controls"
	"github.com/kward/venue/router/signals"
)

// Packet represents a audio action to perform.
type Packet struct {
	Source   string           // Source device of the packet.
	Action   actions.Action   // Action to be, or that was, performed.
	Control  controls.Control // Control to be acted upon.
	Signal   signals.Signal   // Signal being acted upon.
	SignalNo int              // The signal number (e.g. input #1, or aux #3).
	Value    interface{}
}

// Equal returns true if the two packets are equal.
func (p *Packet) Equal(p2 *Packet) bool {
	return reflect.DeepEqual(p, p2)
}

// String returns a human readable representation of the packet.
func (p *Packet) String() string {
	return fmt.Sprintf("{ Source: %s Action: %s, Control: %s Signal: %s SignalNo: %d Value: %v }",
		p.Source, p.Action, p.Control, p.Signal, p.SignalNo, p.Value)
}
