package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/kward/venue/router"
	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/router/signals"
	"github.com/kward/venue/venue"
	"github.com/kward/venue/venuelib"
)

var (
	host            = flag.String("venue_host", "localhost", "Venue host.")
	port            = flag.Uint("venue_port", 5900, "Venue port.")
	passwd          string
	maxProtoVersion = flag.String("vnc_max_proto_version", "", "VNC max protocol version")

	numInputs = flag.Uint("num_inputs", 48, "number of inputs")
	period    = flag.Duration("period", 100*time.Millisecond, "period for random adjustment")
)

func flagInit() {
	flag.StringVar(&passwd, "venue_passwd", "", "Venue password.")
	flag.Parse()
}

func main() {
	flagInit()

	if passwd == "" {
		var err error
		passwd, err = venuelib.GetPasswd()
		if err != nil {
			log.Fatalf("failed to get password; %s", err)
		}
	}

	v, err := venue.New()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	// TODO(kward:20161124) Fix how the context value is handled.
	if *maxProtoVersion != "" {
		ctx = context.WithValue(ctx, "vnc_max_proto_version", *maxProtoVersion)
	}

	// Establish connection with the VENUE VNC server.
	if err := v.Connect(ctx, *host, *port, passwd); err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	log.Println("Venue connection established.")

	v.Initialize()
	go v.ListenAndHandle()

	// Randomly adjust an input.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		v.Handle(&router.Packet{
			SourceName: "rand_inputs",
			Action:     actions.SelectInput,
			Signal:     signals.Input,
			SignalNo:   (signals.SignalNo)(r.Intn(int(*numInputs))),
		})
		if *period == 0 {
			break
		}
		time.Sleep(*period)
	}
}
