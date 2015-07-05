/*
The Venue package exposes the Avid™ VENUE VNC interface as a programmatic API.
*/
package venue

import (
	"flag"
	"fmt"
	"image"
	"log"
	"net"
	"strconv"
	"time"

	vnc "github.com/kward/go-vnc"
	"golang.org/x/net/context"
)

const (
	errPrefix = "Venue error."
	numInputs = 48
)

var (
	refresh = flag.Duration("venue_refresh", 1000*time.Millisecond, "framebuffer refresh period.")
	timeout = flag.Duration("venue_timeout", 10*time.Second, "timeout for Venue connection.")
)

// Venue holds information representing the state of the VENUE backend.
type Venue struct {
	host      string
	port      uint
	cfg       *vnc.ClientConfig
	conn      *vnc.ClientConn
	fb        *Framebuffer
	inputs    [numInputs]*Input
	currInput Input
	currPage  int

	Pages VenuePages
}

type VenuePages map[int]*Page

// NewVenue returns a populated Venue struct.
func NewVenue(host string, port uint, passwd string) *Venue {
	cfg := vnc.NewClientConfig(passwd)
	return &Venue{host: host, port: port, cfg: cfg}
}

// Connect to a VENUE console.
func (v *Venue) Connect(ctx context.Context) error {
	if v.conn != nil {
		return fmt.Errorf("%v Already connected.", errPrefix)
	}

	addr := v.host + ":" + strconv.FormatUint(uint64(v.port), 10)
	netConn, err := net.DialTimeout("tcp", addr, *timeout)
	if err != nil {
		return fmt.Errorf("%v Error connecting to host. %v", errPrefix, err)
	}

	var cancel context.CancelFunc
	if _, ok := ctx.Deadline(); ok {
		ctx, cancel = context.WithCancel(ctx)
	} else {
		ctx, cancel = context.WithTimeout(ctx, *timeout)
	}
	defer cancel()

	vncConn, err := vnc.Connect(ctx, netConn, v.cfg)
	if err != nil {
		return fmt.Errorf("%v Could not establish session. %v", errPrefix, err)
	}
	v.conn = vncConn
	return nil
}

// Close a connection to a VENUE console.
func (v *Venue) Close() error {
	return v.conn.Close()
}

// Initialize the in-memory state representation of a VENUE console.
func (v *Venue) Initialize() {
	// Create image to apply framebuffer updates to.
	v.fb = NewFramebuffer(int(v.conn.FramebufferWidth), int(v.conn.FramebufferHeight))

	// Setup channel to listen to server messages.
	v.cfg.ServerMessageCh = make(chan vnc.ServerMessage)

	// Initialize pages.
	v.Pages = VenuePages{}
	v.Pages[InputsPage] = NewInputsPage()
	v.Pages[OutputsPage] = NewOutputsPage()
	// Initialize inputs.
	for ch := 0; ch < numInputs; ch++ {
		input := NewInput(v, ch+1)
		v.inputs[ch] = input
	}

	v.SetPage(OptionsPage) // Ensure Inputs page shows first bank when selected.
	v.SetPage(InputsPage)
	v.SetInput(1)
	v.MouseMove(image.Point{0, 0})
}

// ListenAndHandle VNC server messages.
func (v *Venue) ListenAndHandle() {
	go v.conn.ListenAndHandle()
	for {
		msg := <-v.cfg.ServerMessageCh
		switch msg.Type() {
		case vnc.FramebufferUpdate:
			//log.Println("ListenAndHandle() FramebufferUpdateMessage")
			for i := uint16(0); i < msg.(*vnc.FramebufferUpdateMessage).NumRect; i++ {
				var colors []vnc.Color
				rect := msg.(*vnc.FramebufferUpdateMessage).Rects[i]
				switch rect.Enc.Type() {
				case vnc.Raw:
					colors = rect.Enc.(*vnc.RawEncoding).Colors
				}
				v.fb.Paint(v, rect, colors)
			}

		default:
			log.Printf("ListenAndHandle() unknown message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}

// FramebufferRefresh refreshes the local framebuffer image of the VNC server.
func (v *Venue) FramebufferRefresh() {
	screen := image.Rectangle{image.Point{0, 0}, image.Point{v.fb.Width, v.fb.Height}}
	for {
		v.Snapshot(screen)
		time.Sleep(*refresh)
	}
}

// Snapshot requests updated image info from the VNC server.
func (v *Venue) Snapshot(r image.Rectangle) error {
	//log.Printf("Snapshot(%v)\n", r)
	w, h := uint16(r.Max.X-r.Min.X), uint16(r.Max.Y-r.Min.Y)
	if err := v.conn.FramebufferUpdateRequest(
		vnc.RFBTrue, uint16(r.Min.X), uint16(r.Min.Y), w, h); err != nil {
		log.Println("Snapshot() error; ", err)
		return err
	}
	return nil
}
