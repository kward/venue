package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"time"

	"github.com/kward/venue/venue"
	"github.com/kward/venue/venuelib"
)

var (
	host   = flag.String("host", "localhost", "Venue host.")
	port   = flag.Uint("port", 5900, "Venue port.")
	passwd string

	numInputs = flag.Uint("num_inputs", 48, "number of inputs")
	period    = flag.Duration("period", 100*time.Millisecond, "period for random adjustment")
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
		v.SelectInput(uint16(rand.Intn(int(*numInputs))))

		if *period == 0 {
			break
		}
		time.Sleep(*period)
	}
}
