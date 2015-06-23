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
		input: 1,
	}
}

func (s *state) handleBundle(b *osc.Bundle) {
	log.Print("OSC Bundle:")
	for i, msg := range b.Messages {
		log.Printf("OSC Message #%d: ", i+1, msg.Address)
	}
}

func (s *state) handleMessage(v *venue.Venue, msg *osc.Message) {
	var (
		dxInput, dyInput int
		dxOutput         int
	)

	addr := msg.Address
	log.Printf("OSC Message: %v", addr)

	// Strip version
	version, addr := car(addr), cdr(addr)
	if version == "/ping" {
		return
	}

	layout, addr := car(addr), cdr(addr)
	switch layout {
	case "pv":
		dxInput, dyInput = 8, 3
		dxOutput = 4
	}

	page, addr := car(addr), cdr(addr)
	log.Printf("Page: %v", page)

	control, addr := car(addr), cdr(addr)
	log.Printf("Control: %v", control)
	switch control {
	case "input":
		action := car(addr)
		switch action {
		case "bank":
			bank := car(cdr(addr))
			log.Printf("Input bank %v selected.", bank)
			switch bank {
			case "a":
				s.inputBank = 0
			case "b":
				s.inputBank = 1
			}

		default:
			val := msg.Arguments[0].(float32)
			if val == 0 { // Only handle presses, not releases.
				break
			}

			x, y := carInt(addr), carInt(cdr(addr))
			ch := x + (y-1)*dxInput
			input := ch + s.inputBank*dxInput*dyInput
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
		action, addr := car(addr), cdr(addr)
		switch action {
		case "bank":
			bank := car(addr)
			log.Printf("Output bank %v selected.", bank)
			switch bank {
			case "a":
				s.outputBank = 0
			case "b":
				s.outputBank = 1
			case "c":
				s.outputBank = 1
			}

		case "level":
			log.Println("Output level.")
			val := msg.Arguments[0].(float32)
			if val == 0 { // Only handle presses, not releases.
				break
			}

			x, y := carInt(addr), carInt(cdr(addr))
			ch := x
			output := (ch+s.outputBank*dxOutput)*2 - 1
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
				v.Page(venue.OutputsPage)
				vp := v.Pages[venue.OutputsPage]
				e := vp.Elements[name+"solo"]
				fmt.Printf("output:%v name:%v element:%v\n", output, name, e)
				e.(*venue.Switch).Select(v)
				v.Page(venue.InputsPage)
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

			x := carInt(addr)
			ch := x
			output := (ch+s.outputBank*dxOutput)*2 - 1
			fmt.Printf("x:%v output:%v\n", x, output)

			// Select output.
			var name string
			if output < 16 {
				name = fmt.Sprintf("aux%d", output) // TOOD(kward): replace aux with constant.
			} else {
				name = fmt.Sprintf("grp%d", output-16)
			}

			// Solo output if needed.
			if s.output != output {
				v.Page(venue.OutputsPage)
				vp := v.Pages[venue.OutputsPage]
				e := vp.Elements[name+"solo"]
				fmt.Printf("output:%v name:%v element:%v\n", output, name, e)
				e.(*venue.Switch).Select(v)
				v.Page(venue.InputsPage)
			}
		}
	}
}

func abs(x int) int {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0
	}
	return x
}

func car(s string) string {
	sp := strings.SplitN(s, "/", 3)
	if len(sp) > 1 {
		return sp[1]
	}
	return ""
}

func carInt(s string) int {
	s = car(s)
	i, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return -1
	}
	return int(i)
}

func cdr(s string) string {
	sp := strings.SplitN(s, "/", 3)
	if len(sp) > 2 {
		return "/" + sp[2]
	}
	return ""
}

func main() {
	flagInit()

	if venuePasswd == "" {
		fmt.Printf("Password: ")
		venuePasswd = string(gopass.GetPasswdMasked())
	}

	v := venue.NewVenue(venueHost, venuePort, venuePasswd)
	if err := v.Connect(); err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	log.Println("Venue connection established.")

	v.Initialize()
	go v.HandleServer()
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
