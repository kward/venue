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
	"github.com/kward/venue/ping"
	"github.com/kward/venue/router"
	"github.com/kward/venue/touchosc"
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

	venueFbRefresh   = flag.Bool("enable_venue_fb_refresh", false, "Enable Venue framebuffer refresh.")
	checkpointPeriod = flag.Duration("checkpoint_period", 1*time.Minute, "Checkpoint period")
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
	router     *router.Router
}

func NewState(router *router.Router) *state {
	return &state{
		input:      1,
		inputBank:  1,
		outputBank: 1,
		router:     router,
	}
}

func (s *state) handleBundle(b *osc.Bundle) {
	if glog.V(2) {
		glog.Infof("Received OSC bundle from %v:", b.Addr())
	}
	for i, msg := range b.Messages {
		if glog.V(4) {
			glog.Infof("OSC message #%d: ", i+1, msg.Address)
		}
	}
	glog.Errorf("%s unimplemented", venuelib.FnName())
}

func (s *state) handleMessage(v *venue.Venue, msg *osc.Message) {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}

	// The address is expected to be in this format:
	// /version/layout/page/control/command[/num1][/num2][/label]
	if glog.V(2) {
		glog.Infof("Received OSC message from %s: %q", msg.Addr(), msg)
	}

	pkt, err := touchosc.Parse(msg)
	if err != nil {
		glog.Errorf("Failed to parse OSC message %s; %s", msg, err)
		return
	}
	if glog.V(4) {
		glog.Infof("Parsed packet: %s", pkt)
	}

	router.Dispatch(s.router, pkt)
}

func main() {
	flagInit()

	log.SetFlags(log.Flags() | log.Lmicroseconds | log.Lshortfile)

	if *venueHost == "" {
		glog.Exitln("missing --venue_host flag")
	}
	passwd := venuePasswd
	if passwd == "" {
		pw, err := venuelib.GetPasswd()
		if err != nil {
			glog.Exitf("Failed to get password; %s\n", err)
		}
		passwd = pw
	}

	// Instantiate Venue client.
	v, err := venue.New()
	if err != nil {
		glog.Exitf("Failure instantiating Venue client; %s\n", err)
	}

	// Establish connection with the Venue VNC server.
	ctx, cancel := context.WithTimeout(context.Background(), *venueTimeout)
	defer cancel()
	if err := v.Connect(ctx, *venueHost, *venuePort, passwd); err != nil {
		glog.Exitf("Failed to connect to Venue VNC server; %s\n", err)
	}
	defer v.Close()
	glog.Info("Venue connection established.")
	if err := v.Initialize(); err != nil {
		glog.Exitf("Unable to initialize Venue properly; %s\n", err)
	}

	router := &router.Router{}
	router.RegisterEndpoint(v)
	router.RegisterEndpoint(&ping.Ping{})

	go v.ListenAndHandle()

	o := &osc.Server{}
	conn, err := net.ListenPacket("udp", fmt.Sprintf("%v:%v", *oscServerHost, *oscServerPort))
	if err != nil {
		glog.Exitf("Error starting OSC server; %s\n", err)
	}
	defer conn.Close()
	glog.Info("OSC server started.")

	go func() {
		s := NewState(router)

		for {
			p, err := o.ReceivePacket(context.Background(), conn)
			if err != nil {
				glog.Exitf("OSC error; %s\n", err)
			}
			if p == nil {
				continue
			}

			switch t := p.(type) {
			case *osc.Bundle:
				s.handleBundle(p.(*osc.Bundle))
			case *osc.Message:
				s.handleMessage(v, p.(*osc.Message))
			default:
				glog.Errorf("unrecognized packet type %v", t)
			}
		}
	}()

	for {
		if glog.V(5) {
			glog.Infof("--- checkpoint ---")
		}
		time.Sleep(*checkpointPeriod)
	}
}
