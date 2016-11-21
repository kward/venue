package workflow

import (
	"fmt"
	"strings"
)

const (
	vertical = iota
	horizontal
)

// A Control alters the value of a control? by a given amount.
type Control func(x, y, dx, dy, b int)

var Workflows map[string]Workflow

// A Workflow maps a given Venue workflow into a specific function.
// TODO(kward:20161118) Why did I do all this again?
type Workflow struct {
	version     string  // Version number (0.1, 0.2)
	layout      string  // TouchOSC Layout (th == tablet horizontal)
	page        string  // OSC page visible (soundcheck)
	control     string  // Surface control (input, output)
	verb        string  // Action verb (select)
	ref1, ref2  int     // Control references (optional)
	label       bool    // Message refers to a text label
	orientation int     // Device orientation (horizontal, vertical)
	fx          Control // Function implementing the control
}

func (w Workflow) Addr() string {
	elem := []string{"", w.version, w.layout, w.page, w.control, w.verb}
	if w.ref1 > 0 {
		elem = append(elem, fmt.Sprintf("%d", w.ref1))
	}
	if w.ref2 > 0 {
		elem = append(elem, fmt.Sprintf("%d", w.ref2))
	}
	if w.label {
		elem = append(elem, "label")
	}
	return strings.Join(elem, "/")
}

// Register a workflow.
func Register(wf Workflow) error {
	addr := wf.Addr()
	if _, ok := Workflows[addr]; ok {
		return fmt.Errorf("workflow already exists: %v", addr)
	}
	Workflows[addr] = wf
	return nil
}

// The multi* UI controls report their x and y position as /X/Y, with x and y
// corresponding to the top-left of the control, with x increasing to the right
// and y increasing downwards, on a vertical orientation. When the layout
// orientation is changed to horizontal, the x and y correspond to the
// bottom-left corner, with x increasing vertically, and y increasing to the
// right.
//
// Vertical: 1, 1 is top-left, X inc right, Y inc down
// | 1 2 3 |
// | 2 2 3 |
// | 3 3 3 |
//
// Horizontal: 1, 1 is bottom-left, X inc up, Y inc right
// | 3 3 3 |
// | 2 2 3 |
// | 1 2 3 |

// multiPosition returns the absolute position on a multi UI control.
func multiPosition(x, y, dx, dy, bank int) int {
	return x + (y-1)*dx + dx*dy*(bank-1)
}

// multiRotate returns rotated x and y values for a dy sized multi UI control.
func multiRotate(x, y, dy int) (int, int) {
	return y, dy - x + 1
}

func init() {
	Workflows = make(map[string]Workflow)
}
