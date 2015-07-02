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
	host       string
	port       uint
	passwd     string
	timeout    time.Duration
	havePasswd bool
)

func flagInit() {
	flag.StringVar(&host, "host", "localhost", "Venue host.")
	flag.UintVar(&port, "port", 5900, "Venue port.")
	flag.StringVar(&passwd, "passwd", "", "Venue password.")

	flag.Parse()
}

func main() {
	flagInit()

	if passwd == "" {
		fmt.Printf("Password: ")
		passwd = string(gopass.GetPasswdMasked())
	}

	ctx := context.Background()
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
		v.Input(i)
		time.Sleep(1 * time.Second)
	}
}
