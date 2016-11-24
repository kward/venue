package vnc

import (
	"context"
	"fmt"
	"image"
	"log"
	"net"
	"time"

	vnclib "github.com/kward/go-vnc"
)

const (
	refresh = 1000 * time.Millisecond
	port    = 5900
)

// The VNC type contains various handles relating to a VNC connection.
type VNC struct {
	opts *options
	cfg  *vnclib.ClientConfig
	conn *vnclib.ClientConn
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
		return fmt.Errorf("already connected")
	}
	// TODO(kward:20161122) Add check for a reasonably sufficient deadline.
	deadline, ok := ctx.Deadline()
	if !ok {
		return fmt.Errorf("context missing deadline")
	}

	log.Println("Connecting to VENUE VNC server...")
	addr := fmt.Sprintf("%s:%d", v.opts.host, v.opts.port)
	nc, err := net.DialTimeout("tcp", addr, time.Since(deadline))
	if err != nil {
		return err
	}

	log.Println("Establishing session...")
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

// ListenAndHandle VNC server messages.
func (v *VNC) ListenAndHandle() {
	log.Println("ListenAndHandle()")
	go v.conn.ListenAndHandle()
	for {
		msg := <-v.cfg.ServerMessageCh
		switch msg.Type() {
		case vnclib.FramebufferUpdateMsg:
			log.Println("ListenAndHandle() FramebufferUpdateMessage")
			for i := uint16(0); i < msg.(*vnclib.FramebufferUpdate).NumRect; i++ {
				var colors []vnclib.Color
				rect := msg.(*vnclib.FramebufferUpdate).Rects[i]
				switch rect.Enc.Type() {
				case vnclib.Raw:
					colors = rect.Enc.(*vnclib.RawEncoding).Colors
				}
				v.fb.Paint(rect, colors)
			}

		default:
			log.Printf("ListenAndHandle() unknown message type:%s msg:%s\n", msg.Type(), msg)
		}
	}
}

// FramebufferRefresh refreshes the local framebuffer image of the VNC server
// every period `p`.
func (v *VNC) FramebufferRefresh(p time.Duration) {
	r := image.Rectangle{image.Point{0, 0}, image.Point{v.fb.Width, v.fb.Height}}
	for {
		if err := v.Snapshot(r); err != nil {
			// TODO(kward:20161124) Return errors on a channel.
			log.Printf("framebuffer refresh error: %s", err)
		}
		if p == 0 {
			break
		}
		time.Sleep(p)
	}
}

// Snapshot requests updated image info from the VNC server.
func (v *VNC) Snapshot(r image.Rectangle) error {
	log.Printf("Snapshot(%v)\n", r)
	w, h := uint16(r.Max.X-r.Min.X), uint16(r.Max.Y-r.Min.Y)
	if err := v.conn.FramebufferUpdateRequest(
		vnclib.RFBTrue, uint16(r.Min.X), uint16(r.Min.Y), w, h); err != nil {
		log.Println("Snapshot() error:", err)
		return err
	}
	return nil
}

// DebugMetrics passes the call through to the connection.
func (v *VNC) DebugMetrics() {
	v.conn.DebugMetrics()
}

// KeyPress presses a key on the VENUE console.
func (v *VNC) KeyPress(key uint32) error {
	if err := v.conn.KeyEvent(key, vnclib.PressKey); err != nil {
		return err
	}
	if err := v.conn.KeyEvent(key, vnclib.ReleaseKey); err != nil {
		return err
	}
	return nil
}

// MouseMove moves the mouse.
func (v *VNC) MouseMove(p image.Point) {
	v.conn.PointerEvent(vnclib.ButtonNone, uint16(p.X), uint16(p.Y))
}

// MouseLeftClick moves the mouse to a position and left clicks.
func (v *VNC) MouseLeftClick(p image.Point) {
	v.MouseMove(p)
	v.conn.PointerEvent(vnclib.ButtonLeft, uint16(p.X), uint16(p.Y))
	v.conn.PointerEvent(vnclib.ButtonNone, uint16(p.X), uint16(p.Y))
}

// MouseDrag moves the mouse, clicks, and drags to a new position.
func (v *VNC) MouseDrag(p, d image.Point) {
	v.MouseMove(p)
	v.conn.PointerEvent(vnclib.ButtonLeft, uint16(p.X), uint16(p.Y))
	p = p.Add(d) // Add delta.
	v.conn.PointerEvent(vnclib.ButtonLeft, uint16(p.X), uint16(p.Y))
	v.conn.PointerEvent(vnclib.ButtonNone, uint16(p.X), uint16(p.Y))
}
