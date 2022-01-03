// Package presets provides support for VENUE presets.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	pb "github.com/kward/venue/presets/proto"
)

func main() {
	fn := "testdata/D-Show Input Channel/211231.00 Ch 1 Clear Console.ich"
	if len(os.Args) > 1 {
		fn = os.Args[1]
	}

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	preset := pb.NewDShowInputChannel()
	preset.Read(b)

	fmt.Printf("strip_name: %s\n", preset.StripName())
}
