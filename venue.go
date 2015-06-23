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
)

const (
	errPrefix = "Venue error."
	numInputs = 48
	uiSettle  = 10 * time.Millisecond // Time to allow UI to settle.
)

var (
	refreshFlag = flag.Duration("refresh", 1000*time.Millisecond, "framebuffer refresh period.")
	timeoutFlag = flag.Duration("timeout", 10*time.Second, "timeout for Venue connection.")
)

type Venue struct {
	host      string
	port      uint
	cfg       *vnc.ClientConfig
	conn      *vnc.ClientConn
	fb        *Framebuffer
	inputs    [numInputs]*Input
	currInput Input
	currPage  Page

	Pages map[int]*Page
}

func NewVenue(host string, port uint, passwd string) *Venue {
	cfg := vnc.NewClientConfig(passwd)
	return &Venue{host: host, port: port, cfg: cfg}
}

func (v *Venue) Connect() error {
	if v.conn != nil {
		return fmt.Errorf("%v Already connected.", errPrefix)
	}

	addr := v.host + ":" + strconv.FormatUint(uint64(v.port), 10)
	netConn, err := net.DialTimeout("tcp", addr, *timeoutFlag)
	if err != nil {
		return fmt.Errorf("%v Error connecting to host. %v", errPrefix, err)
	}

	vncConn, err := vnc.Client(netConn, v.cfg)
	if err != nil {
		return fmt.Errorf("%v Could not establish session. %v", errPrefix, err)
	}
	v.conn = vncConn
	return nil
}

func (v *Venue) Close() error {
	return v.conn.Close()
}

func (v *Venue) Initialize() {
	// Create image to apply framebuffer updates to.
	v.fb = NewFramebuffer(int(v.conn.FramebufferWidth), int(v.conn.FramebufferHeight))

	// Setup channel to listen to server messages.
	v.cfg.ServerMessageCh = make(chan vnc.ServerMessage)

	// Initialize pages.
	v.Pages = map[int]*Page{}
	v.Pages[InputsPage] = NewInputsPage()
	v.Pages[OutputsPage] = NewOutputsPage()
	// Initialize inputs.
	for ch := 0; ch < numInputs; ch++ {
		input := NewInput(v, ch+1)
		v.inputs[ch] = input
	}

	v.Page(OptionsPage) // Ensure Inputs page shows first bank when selected.
	v.Page(InputsPage)
	v.Input(1)
	v.MouseMove(image.Point{0, 0})
}

func (v *Venue) HandleServer() {
	for {
		msg := <-v.cfg.ServerMessageCh
		switch msg.Type() {
		case vnc.FramebufferUpdate:
			log.Println("HandleServer() FramebufferUpdateMsg")
			for i := uint16(0); i < msg.(*vnc.FramebufferUpdateMsg).NumRect; i++ {
				var colors []vnc.Color
				rect := msg.(*vnc.FramebufferUpdateMsg).Rects[i]
				switch rect.Enc.Type() {
				case vnc.RawEnc:
					colors = rect.Enc.(*vnc.RawEncoding).Colors
				}
				v.fb.Paint(v, rect, colors)
			}

		default:
			log.Printf("HandleServer() unknown message type:%v msg:%v\n", msg.Type(), msg)
		}
	}
}

func (v *Venue) FramebufferRefresh() {
	//screen := image.Rectangle{image.Point{0, 0}, image.Point{v.fb.Width, v.fb.Height}}
	for {
		//v.Snapshot(screen)
		time.Sleep(*refreshFlag)
	}
}

func (v *Venue) Snapshot(r image.Rectangle) error {
	log.Printf("Snapshot(%v)\n", r)
	w, h := uint16(r.Max.X-r.Min.X), uint16(r.Max.Y-r.Min.Y)
	if err := v.conn.FramebufferUpdateRequest(
		vnc.RFBTrue, uint16(r.Min.X), uint16(r.Min.Y), w, h); err != nil {
		log.Println("Snapshot() error; ", err)
		return err
	}
	return nil
}
