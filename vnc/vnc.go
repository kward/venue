package vnc

import (
	"context"
	"fmt"
	"image"
	"net"
	"time"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/venue/venuelib"
)

const (
	refresh   = 1000 * time.Millisecond
	port      = 5900
	maxInputs = 96 // Maximum number of signal inputs the code can handle.
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

	glog.Infof("Connecting to VENUE VNC server...")
	addr := fmt.Sprintf("%s:%d", v.opts.host, v.opts.port)
	nc, err := net.DialTimeout("tcp", addr, deadline.Sub(time.Now()))
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

// ListenAndHandle VNC server messages.
func (v *VNC) ListenAndHandle() {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	go v.conn.ListenAndHandle()
	for {
		msg := <-v.cfg.ServerMessageCh
		switch msg.Type() {
		case vnclib.FramebufferUpdateMsg:
			if glog.V(5) {
				glog.Info("ListenAndHandle() FramebufferUpdateMessage")
			}
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
			glog.Errorf("ListenAndHandle() unknown message type:%v msg:%s\n", msg.Type(), msg)
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
			glog.Errorf("framebuffer refresh error: %s", err)
		}
		if p == 0 {
			break
		}
		time.Sleep(p)
	}
}

// Snapshot requests updated image info from the VNC server.
func (v *VNC) Snapshot(r image.Rectangle) error {
	if glog.V(5) {
		glog.Infof("Snapshot(%v)\n", r)
	}
	w, h := uint16(r.Max.X-r.Min.X), uint16(r.Max.Y-r.Min.Y)
	if err := v.conn.FramebufferUpdateRequest(
		vnclib.RFBTrue, uint16(r.Min.X), uint16(r.Min.Y), w, h); err != nil {
		glog.Errorf("Snapshot() error: %s", err)
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
func (v *VNC) MouseMove(p image.Point) error {
	return v.conn.PointerEvent(vnclib.ButtonNone, uint16(p.X), uint16(p.Y))
}

// MouseLeftClick moves the mouse to a position and left clicks.
func (v *VNC) MouseLeftClick(p image.Point) error {
	if err := v.MouseMove(p); err != nil {
		return err
	}
	if err := v.conn.PointerEvent(vnclib.ButtonLeft, uint16(p.X), uint16(p.Y)); err != nil {
		return err
	}
	return v.conn.PointerEvent(vnclib.ButtonNone, uint16(p.X), uint16(p.Y))
}

// MouseDrag moves the mouse, clicks, and drags to a new position.
func (v *VNC) MouseDrag(p, d image.Point) error {
	if err := v.MouseMove(p); err != nil {
		return err
	}
	if err := v.conn.PointerEvent(vnclib.ButtonLeft, uint16(p.X), uint16(p.Y)); err != nil {
		return err
	}
	p = p.Add(d) // Add delta.
	if err := v.conn.PointerEvent(vnclib.ButtonLeft, uint16(p.X), uint16(p.Y)); err != nil {
		return err
	}
	return v.conn.PointerEvent(vnclib.ButtonNone, uint16(p.X), uint16(p.Y))
}

// IntToKeys returns a slice of uint32 values that represent the key presses
// required to input an int.
// TODO(kward:20161126) This should move to upstream VNC library.
func IntToKeys(v int) []uint32 {
	keys := map[rune]uint32{
		'-': vnclib.KeyMinus,
		'0': vnclib.Key0,
		'1': vnclib.Key1,
		'2': vnclib.Key2,
		'3': vnclib.Key3,
		'4': vnclib.Key4,
		'5': vnclib.Key5,
		'6': vnclib.Key6,
		'7': vnclib.Key7,
		'8': vnclib.Key8,
		'9': vnclib.Key9,
	}
	k := []uint32{}
	s := fmt.Sprintf("%d", v)
	for _, c := range s {
		k = append(k, keys[c])
	}
	return k
}
