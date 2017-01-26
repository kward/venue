package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/golang/glog"
	"github.com/kward/go-osc/osc"
	"github.com/kward/venue/oscparse"
	"github.com/kward/venue/venue"
	"github.com/kward/venue/venuelib"
)

var (
	oscClientHost = flag.String("osc_client_host", "127.0.0.1", "OSC client hostname/IP.")
	oscClientPort = flag.Uint("osc_client_port", 9000, "OSC client port.")
	oscServerHost = flag.String("osc_server_host", "0.0.0.0", "OSC server hostname/IP.")
	oscServerPort = flag.Uint("osc_server_port", 8000, "OSC server port.")

	venueHost    = flag.String("venue_host", "", "Venue VNC host/IP.")
	venuePort    = flag.Uint("venue_port", 5900, "Venue VNC port.")
	venuePasswd  string
	venueTimeout = flag.Duration("venue_timeout", 15*time.Second, "Venue VNC timeout.")

	venueFbRefresh = flag.Bool("enable_venue_fb_refresh", false, "Enable Venue framebuffer refresh.")
)

func flagInit() {
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

func (s *state) handleBundle(b *osc.Bundle, remote net.Addr) {
	if glog.V(2) {
		glog.Infof("Received OSC bundle from %v:", remote)
	}
	for i, msg := range b.Messages {
		if glog.V(4) {
			glog.Infof("OSC message #%d: ", i+1, msg.Address)
		}
	}
}

func (s *state) handleMessage(v *venue.Venue, msg *osc.Message, remote net.Addr) {
	// The address is expected to be in this format:
	// /version/layout/page/control/command[/num1][/num2][/label]
	addr := msg.Address
	if glog.V(2) {
		glog.Infof("Received OSC message from %v: %v", remote, addr)
	}

	pkt, err := oscparse.Parse(msg)
	if err != nil {
		glog.Errorf("Failed to parse OSC message %q; %s", msg.Address, err)
		return
	}
	if glog.V(4) {
		glog.Infof("Venue packet: %s", pkt)
	}
}

func main() {
	flagInit()

	log.SetFlags(log.Flags() | log.Lmicroseconds | log.Lshortfile)

	if *venueHost == "" {
		glog.Fatal("missing --venue_host flag")
	}
	if venuePasswd == "" {
		venuePasswd = venuelib.GetPasswd()
	}

	v, err := venue.New()
	if err != nil {
		glog.Fatal(err)
	}

	// Establish connection with the VENUE VNC server.
	ctx, cancel := context.WithTimeout(context.Background(), *venueTimeout)
	defer cancel()
	if err := v.Connect(ctx, *venueHost, *venuePort, venuePasswd); err != nil {
		glog.Fatal(err)
	}
	defer v.Close()
	glog.Info("Venue connection established.")

	v.Initialize()
	//time.Sleep(1 * time.Second)
	go v.ListenAndHandle()

	o := &osc.Server{}
	conn, err := net.ListenPacket("udp", fmt.Sprintf("%v:%v", *oscServerHost, *oscServerPort))
	if err != nil {
		glog.Fatalf("Error starting OSC server:", err)
	}
	defer conn.Close()
	glog.Info("OSC server started.")

	go func() {
		s := NewState()

		for {
			p, remote, err := o.ReceivePacket(context.Background(), conn)
			if err != nil {
				glog.Fatalf("OSC error: %v", err)
			}
			if p == nil {
				continue
			}

			switch t := p.(type) {
			case *osc.Bundle:
				s.handleBundle(p.(*osc.Bundle), remote)
			case *osc.Message:
				s.handleMessage(v, p.(*osc.Message), remote)
			default:
				glog.Errorf("unrecognized packet type %v", t)
			}
		}
	}()

	for {
		if glog.V(5) {
			glog.Infof("--- checkpoint ---")
		}
		time.Sleep(1 * time.Minute)
	}
}
