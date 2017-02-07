package vnc

import (
	"reflect"
	"testing"

	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
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

type mockConn struct {
	events Events
}

func NewMockConn() ClientConn {
	return &mockConn{}
}

func (c *mockConn) FramebufferHeight() uint16 { return 0 }
func (c *mockConn) FramebufferWidth() uint16  { return 0 }

type keyEvent struct {
	key  keys.Key
	down bool
}

func (c *mockConn) KeyEvent(key keys.Key, down bool) error {
	c.events = append(c.events, keyEvent{key, down})
	return nil
}

func (c *mockConn) PointerEvent(button buttons.Button, x, y uint16) error {
	return nil
}

func (c *mockConn) Close() error                                                { return nil }
func (c *mockConn) DebugMetrics()                                               {}
func (c *mockConn) FramebufferUpdateRequest(int uint8, x, y, w, h uint16) error { return nil }
func (c *mockConn) ListenAndHandle() error                                      { return nil }
