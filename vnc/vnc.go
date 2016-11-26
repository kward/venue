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
	ui   *UI
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

// SetPage changes the VENUE page.
func (v *VNC) SetPage(p int) error { return v.Press(v.ui.inputs) }

// Widget returns a pointer to a widget with name `n` on a given page `p`.
func (v *VNC) Widget(p int, n string) Widget {
	switch p {
	case InputsPage:
		return v.ui.inputs.Widget(n)
	case OutputsPage:
		return v.ui.outputs.Widget(n)
	}
	return nil
}

// Outputs returns a pointer to the OUTPUTS page.
func (v *VNC) Outputs() *Page { return v.ui.outputs }

// SelectInput for interaction.
func (v *VNC) SelectInput(input uint16) error {
	if input > maxInputs {
		return fmt.Errorf("input number %d exceeds maximum number of inputs %d", input, maxInputs)
	}
	log.Printf("Selecting input #%v.", input)

	if err := v.SetPage(InputsPage); err != nil {
		return err
	}
	return v.selectInput(1)
}

// selectInput directly.
func (v *VNC) selectInput(input uint16) error {
	if input < 10 {
		if err := v.KeyPress(vnclib.Key0); err != nil {
			return err
		}
	}
	for _, key := range intToKeys(int(input)) {
		if err := v.KeyPress(key); err != nil {
			return err
		}
	}
	// TODO(kward:20161126): Start a timer that expires after 1750ms. Additional
	// key presses aren't allowed until the time expires, but mouse input is.
	time.Sleep(1750 * time.Millisecond)

	return nil
}

// SelectOutput for interaction.
func (v *VNC) SelectOutput(output string) error {
	v.SetPage(OutputsPage)

	// Clear solo.
	log.Printf("Clearing output solo.")
	widget := v.ui.outputs.Widget("solo_clear")
	if err := widget.Press(v); err != nil {
		return err
	}

	// Solo output.
	log.Printf("Soloing %v output.", output)
	widget = v.ui.outputs.Widget(output + "solo")
	if err := widget.Press(v); err != nil {
		return err
	}

	return nil
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
// TODO(kward): Don't ignore the errors!
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

// TODO(kward:20161126) This should move to upstream VNC library.
func intToKeys(v int) []uint32 {
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
