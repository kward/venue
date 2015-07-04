package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/howeyc/gopass"
	osc "github.com/kward/go-osc"
	vnc "github.com/kward/go-vnc"
	"github.com/kward/venue"
	"golang.org/x/net/context"
)

var (
	oscClientHost string
	oscClientPort uint
	oscServerHost string
	oscServerPort uint
	venueHost     string
	venuePort     uint
	venuePasswd   string
)

const (
	maxArrowKeys = 8
)

func flagInit() {
	flag.StringVar(&oscClientHost, "osc_client_host", "127.0.0.1", "OSC client host/IP.")
	flag.UintVar(&oscClientPort, "osc_client_port", 9000, "OSC client port.")

	flag.StringVar(&oscServerHost, "osc_server_host", "0.0.0.0", "OSC client host/IP.")
	flag.UintVar(&oscServerPort, "osc_server_port", 8000, "OSC client port.")

	flag.StringVar(&venueHost, "venue_host", "localhost", "Venue VNC host/IP.")
	flag.UintVar(&venuePort, "venue_port", 5900, "Venue VNC port.")
	flag.StringVar(&venuePasswd, "venue_passwd", "", "Venue VNC password.")

	flag.Parse()
}

type state struct {
	input      int
	inputBank  int
	output     int
	outputBank int
}

func NewState() *state {
	return &state{
		input:      1,
		inputBank:  1,
		outputBank: 1,
	}
}

func (s *state) handleBundle(b *osc.Bundle) {
	log.Print("OSC Bundle:")
	for i, msg := range b.Messages {
		log.Printf("OSC Message #%d: ", i+1, msg.Address)
	}
}

func (s *state) handleMessage(v *venue.Venue, msg *osc.Message) {
	const (
		vertical = iota
		horizontal
	)
	var (
		// The dx and dy vars are always based on a vertical orientation.
		dxInput, dyInput int
		dxOutput         int
		orientation      int
	)

	// The address is expected to be in this format:
	// /version/layout/page/control[/mod][/num][/label]
	addr := msg.Address
	log.Printf("OSC Message: %v", addr)

	version, addr := car(addr), cdr(addr)
	if version != "0.0" {
		log.Printf("Unsupported message.")
		return
	}
	log.Printf("Version: %v", version)

	layout, addr := car(addr), cdr(addr)
	switch layout {
	case "pv":
		dxInput, dyInput = 8, 4
		dxOutput = 6
	case "th":
		dxInput, dyInput = 12, 4
		dxOutput = 12
		orientation = horizontal
	}
	log.Printf("Layout: %v", layout)

	page, addr := car(addr), cdr(addr)
	log.Printf("Page: %v", page)

	control, addr := car(addr), cdr(addr)
	log.Printf("Control: %v", control)
	switch control {
	case "input":
		mod := car(addr)
		switch mod {
		case "bank": // Only present on the phone layout.
			bank := car(cdr(addr))
			log.Printf("Input bank %v selected.", bank)
			switch bank {
			case "a":
				s.inputBank = 1
			case "b":
				s.inputBank = 2
			}

		default:
			val := msg.Arguments[0].(float32)
			if val == 0 { // Only handle presses, not releases.
				fmt.Println("Ignoring release.")
				break
			}

			x, y := toInt(car(addr)), toInt(cadr(addr))
			if orientation == horizontal {
				x, y = multiRotate(x, y, dyInput)
			}
			input := multiPosition(x, y, dxInput, dyInput, s.inputBank)
			fmt.Printf("x:%v y:%v input:%v\n", x, y, input)

			const (
				left  = false
				right = true
			)

			kp := abs(input - s.input)
			if kp <= maxArrowKeys {
				var dir bool
				if input-s.input > 0 {
					dir = right
				}
				for i := 0; i < kp; i++ {
					if dir == left {
						v.KeyPress(vnc.KeyLeft)
					} else {
						v.KeyPress(vnc.KeyRight)
					}
				}
			} else {
				v.Input(input)
			}
			s.input = input
		}

	case "output":
		mod, addr := car(addr), cdr(addr)
		switch mod {
		case "bank": // Only present on the phone layout.
			bank := car(addr)
			log.Printf("Output bank %v selected.", bank)
			switch bank {
			case "a":
				s.outputBank = 1
			case "b":
				s.outputBank = 2
			case "c":
				s.outputBank = 3
			}

		case "level":
			log.Println("Output level.")
			val := msg.Arguments[0].(float32)
			if val == 0 { // Only handle presses, not releases.
				fmt.Println("Ignoring release.")
				break
			}

			x, y := toInt(car(addr)), toInt(cadr(addr))
			if orientation == horizontal {
				x, y = multiRotate(x, y, 4) // TODO(kward): 4 should be a constant.
			}
			output := x*2 - 1
			fmt.Printf("x:%v y:%v output:%v\n", x, y, output)

			var (
				clicks int
			)
			switch y {
			case 1:
				clicks = 6 // +4.1 dB
			case 2:
				clicks = 1 // ~+0.7 dB
			case 3:
				clicks = -1 // ~-0.7 dB
			case 4:
				clicks = -6 // ~-4.1 dB
			}

			// Select output.
			var name string
			if output < 16 {
				name = fmt.Sprintf("aux%d", output) // TOOD(kward): replace aux with constant.
			} else {
				name = fmt.Sprintf("grp%d", output-16)
			}

			// Solo output if needed.
			if s.output != output {
				prevPage := v.Page()
				v.SetPage(venue.OutputsPage)
				vp := v.Pages[venue.OutputsPage]

				// Clear solo.
				fmt.Println("Clearing solo.")
				e := vp.Elements["solo_clear"]
				fmt.Printf("output:%v name:%v element:%v\n", output, "solo_clear", e)
				e.(*venue.Switch).Update(v)

				// Solo output.
				solo := name + "solo"
				e = vp.Elements[solo]
				fmt.Printf("output:%v name:%v element:%v\n", output, solo, e)
				e.(*venue.Switch).Update(v)

				v.SetPage(prevPage)
			}

			// Adjust output.
			vp := v.Pages[venue.InputsPage]
			e := vp.Elements[name]
			fmt.Printf("output:%v name:%v element:%v\n", output, name, e)
			e.(*venue.Encoder).Adjust(v, clicks)

			s.output = output

		case "select":
			log.Println("Output select.")
			val := msg.Arguments[0].(float32)
			if val == 0 { // Only handle presses, not releases.
				break
			}

			x, y := toInt(car(addr)), toInt(cadr(addr))
			fmt.Printf("x:%v y:%v\n", x, y)
			if orientation == horizontal {
				x, y = multiRotate(x, y, 1) // TODO(kward): 1 should be a constant.
			}
			fmt.Printf("x:%v y:%v\n", x, y)
			output := multiPosition(x, y, dxOutput, 1, s.outputBank)*2 - 1
			fmt.Printf("x:%v y:%v output:%v\n", x, y, output)

			// Select output.
			var name string
			if output < 16 {
				name = fmt.Sprintf("aux%d", output) // TOOD(kward): replace aux with constant.
			} else {
				name = fmt.Sprintf("grp%d", output-16)
			}

			// Solo output if needed.
			if s.output != output {
				prevPage := v.Page()
				v.SetPage(venue.OutputsPage)
				vp := v.Pages[venue.OutputsPage]

				// Clear solo.
				fmt.Println("Clearing solo.")
				e := vp.Elements["solo_clear"]
				fmt.Printf("output:%v name:%v element:%v\n", output, "solo_clear", e)
				e.(*venue.Switch).Update(v)

				// Solo output.
				fmt.Println("Soloing output.")
				solo := name + "solo"
				e = vp.Elements[solo]
				fmt.Printf("output:%v name:%v element:%v\n", output, solo, e)
				e.(*venue.Switch).Update(v)

				v.SetPage(prevPage)
			}
		}
	}
}

// The multi* UI controls report their x and y position as /X/Y, with x and y
// corresponding to the top-left of the control, with x increasing to the right
// and y increasing downwards, on a vertical orientation. When the layout
// orientation is changed to horizontal, the x and y correspond to the
// bottom-left corner, with x increasing vertically, and y increasing to the
// right.
//
// Vertical: 1, 1 is top-left, X inc right, Y inc down
// | 1 2 3 |
// | 2 2 3 |
// | 3 3 3 |
//
// Horizontal: 1, 1 is bottom-left, X inc up, Y inc right
// | 3 3 3 |
// | 2 2 3 |
// | 1 2 3 |

// multiPosition returns the absolute position on a multi UI control.
func multiPosition(x, y, dx, dy, bank int) int {
	return x + (y-1)*dx + dx*dy*(bank-1)
}

// multiRotate returns rotated x and y values for a multi UI control.
func multiRotate(x, y, dy int) (int, int) {
	return y, dy - x + 1
}

// abs returns the absolute value of an int.
func abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0
	}
	return x
}

// car returns the first element of an OSC address.
func car(s string) string {
	sp := strings.SplitN(s, "/", 3)
	if len(sp) > 1 {
		return sp[1]
	}
	return ""
}

// cadr is equivalent to car(cdr(s)).
func cadr(s string) string {
	return car(cdr(s))
}

// cdr returns an OSC address sans the first element.
func cdr(s string) string {
	sp := strings.SplitN(s, "/", 3)
	if len(sp) > 2 {
		return "/" + sp[2]
	}
	return ""
}

// toInt converts a string to an int.
func toInt(s string) int {
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return -1
	}
	return int(i)
}

func main() {
	flagInit()

	if venuePasswd == "" {
		fmt.Printf("Password: ")
		venuePasswd = string(gopass.GetPasswdMasked())
	}

	v := venue.NewVenue(venueHost, venuePort, venuePasswd)
	if err := v.Connect(context.Background()); err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	log.Println("Venue connection established.")

	v.Initialize()
	go v.ListenAndHandle()
	go v.FramebufferRefresh()

	server := &osc.Server{}
	conn, err := net.ListenPacket("udp", fmt.Sprintf("%v:%v", oscServerHost, oscServerPort))
	if err != nil {
		log.Fatal("Error starting OSC server:", err)
	}
	defer conn.Close()
	log.Println("OSC server started.")

	go func() {
		s := NewState()

		for {
			p, err := server.ReceivePacket(conn)
			if err != nil {
				log.Fatalf("OSC error: %v", err)
			}
			if p == nil {
				continue
			}

			switch p.(type) {
			case *osc.Bundle:
				s.handleBundle(p.(*osc.Bundle))
			case *osc.Message:
				s.handleMessage(v, p.(*osc.Message))
			default:
				log.Println("Error: Unrecognized packet type.")
			}
		}
	}()

	for {
		log.Println("--- checkpoint ---")
		time.Sleep(1 * time.Minute)
	}
}
