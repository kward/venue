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
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.Info("VENUE presets")

	fn := testdata + "/D-Show Input Channel/211231.00 Ch 1 Clear Console.ich"
	if len(os.Args) > 1 {
		fn = os.Args[1]
	}

	b, err := ioutil.ReadFile(fn)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	preset := presets.NewDShowInputChannel()
	preset.Read(b)

	fmt.Println(preset)
}
