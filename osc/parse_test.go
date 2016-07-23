package oscparse

import "github.com/kward/venue"

type parseTest struct {
	name string
	addr string
	val  []interface{}
	pkt  *venue.Packet
	ok   bool
}
