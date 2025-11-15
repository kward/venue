package vnc

import (
	"context"
	"image"
	"net"
	"strconv"
	"time"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/go-vnc/buttons"
	"github.com/kward/go-vnc/encodings"
	"github.com/kward/go-vnc/keys"
	"github.com/kward/go-vnc/messages"
	"github.com/kward/go-vnc/rfbflags"
	"github.com/kward/venue/codes"
	"github.com/kward/venue/venuelib"
)

const (
	refresh   = 1000 * time.Millisecond
	port      = 5900
	maxInputs = 96 // Maximum number of signal inputs the code can handle.
)

// ClientConn is a local interface to enable mocking of the go-vnc ClientConn.
type ClientConn interface {
	FramebufferHeight() uint16
	FramebufferWidth() uint16
	KeyEvent(key keys.Key, down bool) error
	PointerEvent(button buttons.Button, x, y uint16) error

	Close() error
	DebugMetrics()
	FramebufferUpdateRequest(inc rfbflags.RFBFlag, x, y, w, h uint16) error
	ListenAndHandle() error
}

// The VNC type contains various handles relating to a VNC connection.
type VNC struct {
	opts *options
	cfg  *vnclib.ClientConfig
	conn ClientConn
	fb   *Framebuffer
}

// New returns a populated VNC structure.
func New(opts ...func(*options) error) (*VNC, error) {
	o := &options{}
	o.setPort(port)
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	return &VNC{opts: o}, nil
}

// Close a VNC connection.
func (v *VNC) Close() error {
	return v.conn.Close()
}

// Connect to a Venue console.
func (v *VNC) Connect(ctx context.Context) error {
	if v.conn != nil {
		return venuelib.Errorf(codes.FailedPrecondition, "already connected")
	}
	// TODO(kward:20161122) Add check for a reasonably sufficient deadline.
	deadline, ok := ctx.Deadline()
	if !ok {
		return venuelib.Errorf(codes.FailedPrecondition, "context missing deadline")
	}

	glog.Infof("Connecting to VENUE VNC server...")
	addr := net.JoinHostPort(v.opts.host, strconv.Itoa(int(v.opts.port)))
	nc, err := net.DialTimeout("tcp", addr, time.Until(deadline))
	if err != nil {
		return err
	}

	glog.Infof("Establishing session...")
	v.cfg = vnclib.NewClientConfig(v.opts.passwd)
	conn, err := vnclib.Connect(ctx, nc, v.cfg)
	if err != nil {
		return err
	}
	v.conn = conn

	// Initialize a framebuffer for updates.
	v.fb = NewFramebuffer(int(v.conn.FramebufferWidth()), int(v.conn.FramebufferHeight()))
	// Setup channel to listen to server messages.
	v.cfg.ServerMessageCh = make(chan vnclib.ServerMessage)

	return nil
}

func (v *VNC) ClientConn() ClientConn {
	return v.conn
}

// ListenAndHandle VNC server messages.
// ListenAndHandle maintains backward compatibility by using a background context.
// Deprecated: prefer ListenAndHandleCtx.
func (v *VNC) ListenAndHandle() { v.ListenAndHandleCtx(context.Background()) }

// ListenAndHandleCtx listens for server messages until the context is cancelled
// or the ServerMessage channel is closed.
func (v *VNC) ListenAndHandleCtx(ctx context.Context) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = v.conn.ListenAndHandle()
	}()
	for {
		select {
		case <-ctx.Done():
			if glog.V(2) {
				glog.Infof("ListenAndHandleCtx stopping: %v", ctx.Err())
			}
			// Close the connection to unblock the listener and wait briefly
			// for the goroutine to exit to avoid leaks.
			_ = v.conn.Close()
			select {
			case <-done:
			case <-time.After(2 * time.Second):
			}
			return
		case msg, ok := <-v.cfg.ServerMessageCh:
			if !ok {
				if glog.V(2) {
					glog.Info("ListenAndHandleCtx channel closed")
				}
				return
			}
			switch msg.Type() {
			case messages.FramebufferUpdate:
				if glog.V(5) {
					glog.Info("ListenAndHandleCtx FramebufferUpdateMessage")
				}
				fu := msg.(*vnclib.FramebufferUpdate)
				for i := uint16(0); i < fu.NumRect; i++ {
					var colors []vnclib.Color
					rect := fu.Rects[i]
					switch rect.Enc.Type() {
					case encodings.Raw:
						colors = rect.Enc.(*vnclib.RawEncoding).Colors
					}
					v.fb.Paint(rect, colors)
				}
			default:
				glog.Errorf("ListenAndHandleCtx unknown message type:%d msg:%s", msg.Type(), msg)
			}
		}
	}
}

// FramebufferRefresh refreshes the local framebuffer image of the VNC server
// every period `p`.
// FramebufferRefresh maintains backward compatibility by using a background context.
// Deprecated: prefer FramebufferRefreshCtx.
func (v *VNC) FramebufferRefresh(p time.Duration) { v.FramebufferRefreshCtx(context.Background(), p) }

// FramebufferRefreshCtx periodically requests framebuffer updates until the
// context is cancelled. If p == 0 a single refresh is performed.
func (v *VNC) FramebufferRefreshCtx(ctx context.Context, p time.Duration) {
	r := image.Rectangle{image.Point{0, 0}, image.Point{v.fb.Width(), v.fb.Height()}}
	if p == 0 {
		_ = v.Snapshot(r)
		return
	}
	ticker := time.NewTicker(p)
	defer ticker.Stop()
	for {
		if err := v.Snapshot(r); err != nil {
			glog.Errorf("framebuffer refresh error: %s", err)
		}
		select {
		case <-ctx.Done():
			if glog.V(2) {
				glog.Infof("FramebufferRefreshCtx stopping: %v", ctx.Err())
			}
			return
		case <-ticker.C:
			continue
		}
	}
}

// Snapshot requests updated image info from the VNC server.
func (v *VNC) Snapshot(r image.Rectangle) error {
	if glog.V(5) {
		glog.Infof("Snapshot(%v)\n", r)
	}
	w, h := uint16(r.Max.X-r.Min.X), uint16(r.Max.Y-r.Min.Y)
	if err := v.conn.FramebufferUpdateRequest(
		rfbflags.RFBTrue, uint16(r.Min.X), uint16(r.Min.Y), w, h); err != nil {
		glog.Errorf("Snapshot() error: %s", err)
		return err
	}
	return nil
}

// DebugMetrics passes the call through to the connection.
func (v *VNC) DebugMetrics() {
	v.conn.DebugMetrics()
}
