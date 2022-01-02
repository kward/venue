// Package presets provides support for VENUE presets.
package main

import (
	"io/ioutil"
	"os"

	pb "github.com/kward/venue/presets/proto"
)

func main() {
	b, err := ioutil.ReadFile("testdata/D-Show Input Channel/211231.00 Ch 1 Clear Console.ich")
	if err != nil {
		os.Exit(1)
	}

	preset := pb.NewDShowInputChannel()
	preset.Read(b)
}
