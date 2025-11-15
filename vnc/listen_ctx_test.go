package vnc

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	vnclib "github.com/kward/go-vnc"
	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/go-vnc/rfbflags"
)

type mockConnListen struct {
	startedCh chan struct{}
	stopCh    chan struct{}
	closed    atomic.Bool
}

func newMockConnListen() *mockConnListen {
	return &mockConnListen{startedCh: make(chan struct{}), stopCh: make(chan struct{})}
}

// Implement the subset of ClientConn used by ListenAndHandleCtx.
func (m *mockConnListen) FramebufferHeight() uint16                          { return 0 }
func (m *mockConnListen) FramebufferWidth() uint16                           { return 0 }
func (m *mockConnListen) KeyEvent(_ keys.Key, _ bool) error                  { return nil }
func (m *mockConnListen) PointerEvent(_ buttons.Button, _x, _y uint16) error { return nil }
func (m *mockConnListen) DebugMetrics()                                      {}
func (m *mockConnListen) FramebufferUpdateRequest(_ rfbflags.RFBFlag, _, _, _, _ uint16) error {
	return nil
}

func (m *mockConnListen) ListenAndHandle() error {
	close(m.startedCh)
	<-m.stopCh
	return nil
}

func (m *mockConnListen) Close() error {
	if !m.closed.Load() {
		m.closed.Store(true)
		close(m.stopCh)
	}
	return nil
}

func TestListenAndHandleCtx_CancelClosesAndReturns(t *testing.T) {
	mc := newMockConnListen()
	v := &VNC{conn: mc, cfg: &vnclib.ClientConfig{}}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan struct{})
	go func() { v.ListenAndHandleCtx(ctx); close(done) }()

	// Ensure the internal listener started.
	select {
	case <-mc.startedCh:
	case <-time.After(1 * time.Second):
		t.Fatal("listener did not start")
	}

	// Cancel and ensure shutdown completes and connection was closed.
	cancel()

	select {
	case <-done:
		if !mc.closed.Load() {
			t.Fatal("expected connection to be closed on cancel")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("ListenAndHandleCtx did not return after cancel")
	}
}

func TestListenAndHandleCtx_ChannelClosedReturns(t *testing.T) {
	mc := newMockConnListen()
	v := &VNC{conn: mc, cfg: &vnclib.ClientConfig{ServerMessageCh: make(chan vnclib.ServerMessage)}}

	// Close the server message channel to simulate remote shutdown.
	close(v.cfg.ServerMessageCh)

	done := make(chan struct{})
	go func() { v.ListenAndHandleCtx(context.Background()); close(done) }()

	select {
	case <-done:
		// ok
	case <-time.After(1 * time.Second):
		t.Fatal("ListenAndHandleCtx did not return after channel close")
	}
}
