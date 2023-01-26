// Package presets provides support for VENUE presets.
package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kward/venue/presets"
	log "github.com/sirupsen/logrus"
)

const testdata = "presets/testdata"

func init() {
	// log.SetLevel(log.TraceLevel)
	// log.SetLevel(log.DebugLevel)
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Info("VENUE presets")

	fn := testdata + "/D-Show Input Channel/211231.00 Ch 1 Clear Console.ich"
	if len(os.Args) > 1 {
		fn = os.Args[1]
	}

	bs, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	p := presets.NewDShowInputChannel()
	if p == nil {
		fmt.Println("error: failed to initialize DShowInputChannel")
		os.Exit(1)
	}
	p.Read(bs)
	fmt.Println(p)

	bsNew, err := p.Marshal()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if old, new := len(bs), len(bsNew); old != new {
		fmt.Printf("lengths don't match: %d != %d\n", old, new)
	}

	if err := ioutil.WriteFile("/Users/kward/tmp/output", bsNew, 0644); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
