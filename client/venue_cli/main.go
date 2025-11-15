package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kward/venue/router"
	"github.com/kward/venue/router/actions"
	"github.com/kward/venue/router/signals"
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

	// App context cancelled on SIGINT/SIGTERM.
	ctxApp, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := v.Connect(ctxApp, *host, *port, passwd); err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	log.Println("Venue connection established.")

	v.Initialize()
	go v.ListenAndHandleCtx(ctxApp)
	//go v.FramebufferRefresh()

	// Randomly adjust an input.
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		v.Handle(&router.Packet{
			SourceName: "venue_cli",
			Action:     actions.SelectInput,
			Signal:     signals.Input,
			SignalNo:   (signals.SignalNo)(r.Intn(int(*numInputs))),
		})
		if *period == 0 {
			break
		}
		select {
		case <-ctxApp.Done():
			return
		case <-time.After(*period):
		}
	}
}
