/*
A workflow wraps a set of VNC requests together into a single object, allowing
that object to be acted upon as a whole, rather than as individual requests.

The purpose of a workflow is to enable the VENUE VNC server to focus on a
single client at a time, as simultaneous requests from multiple clients would
otherwise conflict with one another.
*/
package vnc

import (
	"fmt"
	"image"
	"time"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/venue/codes"
	"github.com/kward/venue/venuelib"
	"github.com/kward/venue/vnc/messages"
)

// Event describes a single workflow event.
type Event struct {
	desc string           // Description of the event.
	msg  messages.Message // Type of event.
	data interface{}      // Event data.
}

type keyEvent struct {
	key  keys.Key
	down bool
}

type pointerEvent struct {
	button buttons.Button
	x, y   uint16
}

type sleepEvent struct {
	d time.Duration
}

// Workflow holds a client connection to the VNC server, and a list of events.
type Workflow struct {
	conn    ClientConn
	events  []*Event
	sleeper Sleeper
}

// NewWorkflow returns a new workflow object.
func NewWorkflow(conn ClientConn) *Workflow {
	return &Workflow{
		conn:    conn,
		sleeper: newWorkflowSleeper(),
	}
}

func (wf *Workflow) enqueue(e *Event) {
	wf.events = append(wf.events, e)
}

// Execute the workflow against the VNC server.
func (wf *Workflow) Execute() error {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	if wf.conn == nil {
		return venuelib.Errorf(codes.Internal, "invalid VNC connection")
	}

	for eventNo, event := range wf.events {
		// TODO(kward:20170207) Send an 'ESC' key for certain workflow errors.
		if glog.V(4) {
			glog.Infof("Handling event #%d: %s", eventNo+1, event.desc)
		}
		switch event.msg {
		case messages.KeyEvent:
			e := event.data.(keyEvent)
			if err := wf.conn.KeyEvent(e.key, e.down); err != nil {
				return err
			}
		case messages.PointerEvent:
			e := event.data.(pointerEvent)
			if err := wf.conn.PointerEvent(e.button, e.x, e.y); err != nil {
				return err
			}
		case messages.Sleep:
			e := event.data.(sleepEvent)
			wf.sleeper.Sleep(e.d)
			time.Sleep(e.d)
		}
		time.Sleep(10 * time.Millisecond)
	}

	return nil
}

// KeyPress presses a key on the VENUE console.
func (wf *Workflow) KeyPress(key keys.Key) {
	wf.enqueue(&Event{
		fmt.Sprintf("press %s", key),
		messages.KeyEvent,
		keyEvent{key, vnclib.PressKey}})
	wf.enqueue(&Event{
		fmt.Sprintf("release %s", key),
		messages.KeyEvent,
		keyEvent{key, vnclib.ReleaseKey}})
}

// MouseMove moves the mouse.
func (wf *Workflow) MouseMove(p image.Point) {
	wf.enqueue(&Event{
		fmt.Sprintf("move mouse to %s", p),
		messages.PointerEvent,
		pointerEvent{buttons.None, uint16(p.X), uint16(p.Y)}})
}

// MouseClick moves the mouse to a position and left clicks.
func (wf *Workflow) MouseClick(b buttons.Button, p image.Point) {
	wf.MouseMove(p)
	wf.enqueue(&Event{
		fmt.Sprintf("%s button click at %s", b, p),
		messages.PointerEvent,
		pointerEvent{b, uint16(p.X), uint16(p.Y)}})
	wf.enqueue(&Event{
		fmt.Sprintf("%s button release", b),
		messages.PointerEvent,
		pointerEvent{buttons.None, uint16(p.X), uint16(p.Y)}})
}

// MouseDrag moves the mouse, clicks, and drags to a new position.
func (wf *Workflow) MouseDrag(p, d image.Point) {
	wf.MouseMove(p)
	wf.enqueue(&Event{
		fmt.Sprintf("%s button click at %s", buttons.Left, p),
		messages.PointerEvent,
		pointerEvent{buttons.Left, uint16(p.X), uint16(p.Y)}})
	p = p.Add(d) // Add delta.
	wf.enqueue(&Event{
		fmt.Sprintf("mouse drag to %s", p),
		messages.PointerEvent,
		pointerEvent{buttons.Left, uint16(p.X), uint16(p.Y)}})
	wf.enqueue(&Event{
		fmt.Sprintf("%s button release", buttons.Left),
		messages.PointerEvent,
		pointerEvent{buttons.None, uint16(p.X), uint16(p.Y)}})
}

// Sleep the workflow for at least the duration d.
func (wf *Workflow) Sleep(d time.Duration) {
	wf.enqueue(&Event{
		fmt.Sprintf("sleep for %s", d),
		messages.Sleep,
		sleepEvent{d}})
}

type Sleeper interface {
	Sleep(d time.Duration)
}

type workflowSleeper struct{}

func newWorkflowSleeper() Sleeper { return &workflowSleeper{} }

func (s *workflowSleeper) Sleep(d time.Duration) { time.Sleep(d) }
