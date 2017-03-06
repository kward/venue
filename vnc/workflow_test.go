package vnc

import (
	"image"
	"reflect"
	"testing"
	"time"

	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/go-vnc/rfbflags"
)

type Events []interface{}

func TestKeyPress(t *testing.T) {
	conn := NewMockConn()
	wf := NewWorkflow(conn)
	wf.KeyPress(keys.F1)
	wf.Execute()

	events := Events{
		keyEvent{keys.F1, true},
		keyEvent{keys.F1, false},
	}
	if got, want := conn.(*mockConn).events, events; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected events; < %v > != < %v >", got, want)
	}
}

func TestMouseMove(t *testing.T) {
	conn := NewMockConn()
	wf := NewWorkflow(conn)
	wf.MouseMove(image.Point{100, 200})
	wf.Execute()

	events := Events{
		pointerEvent{buttons.None, 100, 200},
	}
	if got, want := conn.(*mockConn).events, events; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected events; < %v > != < %v >", got, want)
	}
}

func TestMouseClick(t *testing.T) {
	conn := NewMockConn()
	wf := NewWorkflow(conn)
	wf.MouseClick(buttons.Left, image.Point{100, 200})
	wf.Execute()

	events := Events{
		pointerEvent{buttons.None, 100, 200},
		pointerEvent{buttons.Left, 100, 200},
		pointerEvent{buttons.None, 100, 200},
	}
	if got, want := conn.(*mockConn).events, events; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected events; < %v > != < %v >", got, want)
	}
}

func TestMouseDrag(t *testing.T) {
	conn := NewMockConn()
	wf := NewWorkflow(conn)
	wf.MouseDrag(image.Point{100, 200}, image.Point{10, -20})
	wf.Execute()

	events := Events{
		pointerEvent{buttons.None, 100, 200},
		pointerEvent{buttons.Left, 100, 200},
		pointerEvent{buttons.Left, 110, 180},
		pointerEvent{buttons.None, 110, 180},
	}
	if got, want := conn.(*mockConn).events, events; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected events; < %v > != < %v >", got, want)
	}
}

func TestSleep(t *testing.T) {
	conn := NewMockConn()
	wf := NewWorkflow(conn)
	wf.sleeper = newMockSleeper()
	wf.Sleep(1 * time.Second)
	wf.Execute()

	events := Events{
		sleepEvent{1 * time.Second},
	}
	if got, want := wf.sleeper.(*mockSleeper).events, events; !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected events; < %v > != < %v >", got, want)
	}
}

//-----------------------------------------------------------------------------
// mockConn implements the ClientConn interface.
type mockConn struct {
	events Events
}

func NewMockConn() ClientConn { return &mockConn{} }

func (c *mockConn) FramebufferHeight() uint16 { return 0 }
func (c *mockConn) FramebufferWidth() uint16  { return 0 }

func (c *mockConn) KeyEvent(key keys.Key, down bool) error {
	c.events = append(c.events, keyEvent{key, down})
	return nil
}

func (c *mockConn) PointerEvent(button buttons.Button, x, y uint16) error {
	c.events = append(c.events, pointerEvent{button, x, y})
	return nil
}

func (c *mockConn) Close() error                                                           { return nil }
func (c *mockConn) DebugMetrics()                                                          {}
func (c *mockConn) FramebufferUpdateRequest(inc rfbflags.RFBFlag, x, y, w, h uint16) error { return nil }
func (c *mockConn) ListenAndHandle() error                                                 { return nil }

//-----------------------------------------------------------------------------
// mockSleeper implements the Sleeper interface.
type mockSleeper struct {
	events Events
}

func newMockSleeper() Sleeper { return &mockSleeper{} }

func (s *mockSleeper) Sleep(d time.Duration) { s.events = append(s.events, sleepEvent{d}) }
