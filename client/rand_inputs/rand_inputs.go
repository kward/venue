package main

import (
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/kward/venue"
	"github.com/kward/venue/venuelib"
	"golang.org/x/net/context"
)

var (
	host            = flag.String("venue_host", "localhost", "Venue host.")
	port            = flag.Uint("venue_port", 5900, "Venue port.")
	passwd          string
	maxProtoVersion = flag.String("vnc_max_proto_version", "", "VNC max protocol version")
)

func flagInit() {
	flag.StringVar(&passwd, "venue_passwd", "", "Venue password.")
	flag.Parse()
}

func main() {
	flagInit()

	if passwd == "" {
		passwd = venuelib.GetPasswd()
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
		i := r.Intn(48)
		v.SetInput(i)
		time.Sleep(2 * time.Second)
	}
}
