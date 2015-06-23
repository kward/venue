package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/howeyc/gopass"
	"github.com/kward/venue"
)

var (
	host       string
	port       uint
	passwd     string
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

	v := venue.NewVenue(host, port, passwd)
	if err := v.Connect(); err != nil {
		log.Fatal(err)
	}
	defer v.Close()
	log.Println("Venue connection established.")

	v.Initialize()
	go v.HandleServer()
	go v.FramebufferRefresh()

	// Randomly adjust an input.
	for {
		i := rand.Intn(48)
		v.Input(i)
		time.Sleep(1 * time.Second)
	}
}
