package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/howeyc/gopass"
	"github.com/kward/venue"
	"golang.org/x/net/context"
)

var (
	host            string
	port            uint
	passwd          string
	timeout         time.Duration
	havePasswd      bool
	maxProtoVersion string
)

func flagInit() {
	flag.StringVar(&host, "venue_host", "localhost", "Venue host.")
	flag.UintVar(&port, "venue_port", 5900, "Venue port.")
	flag.StringVar(&passwd, "venue_passwd", "", "Venue password.")
	flag.StringVar(&maxProtoVersion, "vnc_max_proto_version", "", "VNC max protocol version")

	flag.Parse()
}

func main() {
	flagInit()

	if passwd == "" {
		fmt.Printf("Password: ")
		passwd = string(gopass.GetPasswdMasked())
	}

	ctx := context.Background()
	if maxProtoVersion != "" {
		ctx = context.WithValue(ctx, "vnc_max_proto_version", maxProtoVersion)
	}

	v := venue.NewVenue(host, port, passwd)
	if err := v.Connect(ctx); err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	log.Println("Venue connection established.")

	v.Initialize()
	go v.ListenAndHandle()
	go v.FramebufferRefresh()

	// Randomly adjust an input.
	for {
		i := rand.Intn(48)
		v.SetInput(i)
		time.Sleep(1 * time.Second)
	}
}
