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
	host   = flag.String("host", "localhost", "Venue host.")
	port   = flag.Uint("port", 5900, "Venue port.")
	passwd string
	period = flag.Duration("period", 100*time.Millisecond, "period for random adjustment")
)

func flagInit() {
	flag.StringVar(&passwd, "passwd", "", "Venue password.")
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
	if err := v.Connect(ctx, *host, *port, passwd); err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	log.Println("Venue connection established.")

	v.Initialize()
	go v.ListenAndHandle()
	//go v.FramebufferRefresh()

	// Randomly adjust an input.
	for {
		i := rand.Intn(48)
		v.SetInput(i)

		if *period == 0 {
			break
		}
		time.Sleep(*period)
	}
}
