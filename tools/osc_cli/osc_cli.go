package main

import (
	"flag"

	"github.com/kward/go-osc/osc"
)

const base = "/0.1/th/soundcheck"

var (
	host string
	port int
)

func flagInit() {
	flag.StringVar(&host, "host", "localhost", "host")
	flag.IntVar(&port, "port", 8000, "port")

	flag.Parse()
}

func main() {
	flagInit()

	var msg *osc.Message

	client := osc.NewClient(host, port)

	msg = osc.NewMessage(base + "/output/select/1/1")
	msg.Append(float32(1))
	msg.Append(float32(0))
	client.Send(msg)
	msg = osc.NewMessage(base + "/output/pan")
	msg.Append(float32(10))
	client.Send(msg)

	msg = osc.NewMessage(base + "/output/select/1/2")
	msg.Append(float32(1))
	msg.Append(float32(0))
	client.Send(msg)
	msg = osc.NewMessage(base + "/output/pan")
	msg.Append(float32(-20))
	client.Send(msg)
}
