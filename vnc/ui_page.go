package vnc

import (
	"fmt"
	"image"
	"log"

	vnclib "github.com/kward/go-vnc"
)

type Page struct {
	widgets map[string]Widget
}

type Pages map[string]*Page

// Verify that the Widget interface is honored.
var _ Widget = new(Page)

const (
	InputsPage = iota
	OutputsPage
	FilingPage
	SnapshotsPage
	PatchbayPage
	PluginsPage
	OptionsPage
)

var (
	pageNames = map[int]string{
		InputsPage:    "INPUTS",
		OutputsPage:   "OUTPUTS",
		FilingPage:    "FILING",
		SnapshotsPage: "SNAPSHOTS",
		PatchbayPage:  "PATCHBAY",
		PluginsPage:   "PLUG-INS",
		OptionsPage:   "OPTIONS",
	}
)

func (w *Page) Read(v *VNC) (interface{}, error) { return nil, nil }

func (w *Page) Update(v *VNC, val interface{}) error {
	var page int = val.(int)
	if err := v.KeyPress(vnclib.KeyF1 + uint32(page)); err != nil {
		return err
	}
	return nil
}

func (w *Page) Press(v *VNC) error { return nil }

const (
	bankX  = 8   // X position of 1st bank.
	bankDX = 131 // dX between banks.
	chanDX = 15  // dX between channels in a bank.

	// Inputs
	auxOddX  = 316
	auxPanX  = 473
	aux12Y   = 95
	aux34Y   = 146
	aux56Y   = 197
	aux78Y   = 248
	aux910Y  = 299
	aux1112Y = 350
	aux1314Y = 401
	aux1516Y = 452

	// Outputs
	meterY = 512
	muteY  = 588
	soloY  = 573
)

// NewInputsPage returns a populated Inputs page.
func NewInputsPage() *Page {
	return &Page{Widgets{
		"gain":       &Encoder{image.Point{167, 279}, encoderBL, true},
		"delay":      &Encoder{image.Point{168, 387}, encoderBL, false},
		"hpf":        &Encoder{image.Point{168, 454}, encoderBL, true},
		"pan":        &Encoder{image.Point{239, 443}, encoderBC, false},
		"var_groups": NewPushButton(226, 299, mediumSwitch),
		"aux1":       &Encoder{image.Point{auxOddX, aux12Y}, encoderTR, true},
		"aux1pan":    &Encoder{image.Point{auxPanX, aux12Y}, encoderTL, false},
		"aux3":       &Encoder{image.Point{auxOddX, aux34Y}, encoderTR, true},
		"aux3pan":    &Encoder{image.Point{auxPanX, aux34Y}, encoderTL, false},
		"aux5":       &Encoder{image.Point{auxOddX, aux56Y}, encoderTR, true},
		"aux5pan":    &Encoder{image.Point{auxPanX, aux56Y}, encoderTL, false},
		"aux7":       &Encoder{image.Point{auxOddX, aux78Y}, encoderTR, true},
		"aux7pan":    &Encoder{image.Point{auxPanX, aux78Y}, encoderTL, false},
		"aux9":       &Encoder{image.Point{auxOddX, aux910Y}, encoderTR, true},
		"aux9pan":    &Encoder{image.Point{auxPanX, aux910Y}, encoderTL, false},
		"aux11":      &Encoder{image.Point{auxOddX, aux1112Y}, encoderTR, true},
		"aux11pan":   &Encoder{image.Point{auxPanX, aux1112Y}, encoderTL, false},
		"aux13":      &Encoder{image.Point{auxOddX, aux1314Y}, encoderTR, true},
		"aux13pan":   &Encoder{image.Point{auxPanX, aux1314Y}, encoderTL, false},
		"aux15":      &Encoder{image.Point{auxOddX, aux1516Y}, encoderTR, true},
		"aux15pan":   &Encoder{image.Point{auxPanX, aux1516Y}, encoderTL, false},
		"grp1":       &Encoder{image.Point{auxOddX, aux12Y}, encoderTR, true},
		"grp1pan":    &Encoder{image.Point{auxPanX, aux12Y}, encoderTL, false},
		"grp3":       &Encoder{image.Point{auxOddX, aux34Y}, encoderTR, true},
		"grp3pan":    &Encoder{image.Point{auxPanX, aux34Y}, encoderTL, false},
		"grp5":       &Encoder{image.Point{auxOddX, aux56Y}, encoderTR, true},
		"grp5pan":    &Encoder{image.Point{auxPanX, aux56Y}, encoderTL, false},
		"grp7":       &Encoder{image.Point{auxOddX, aux78Y}, encoderTR, true},
		"grp7pan":    &Encoder{image.Point{auxPanX, aux78Y}, encoderTL, false},
		"solo_clear": NewPushButton(979, 493, mediumSwitch),
	}}
}

// NewOutputsPage returns a populated Outputs page.
func NewOutputsPage() *Page {
	widgets := Widgets{
		"solo_clear": NewPushButton(980, 490, mediumSwitch),
	}

	// Auxes
	for _, b := range []int{1, 2} { // bank
		pre := "aux"
		for c := 1; c <= 8; c++ { // bank channel
			ch, x := (b-1)*8+c, bankX+(b-1)*bankDX+(c-1)*chanDX

			n := fmt.Sprintf("%v%vsolo", pre, ch)
			log.Println("NewOutput() element[%v]:", n)
			widgets[n] = NewToggle(x, soloY, tinySwitch, false)

			n = fmt.Sprintf("%v%vmeter", pre, ch)
			log.Println("NewOutput() element[%v]:", n)
			widgets[n] = &Meter{
				pos:  image.Point{x, meterY},
				size: smallVMeter,
			}
		}
	}

	// Groups
	b := 5 // bank
	pre := "grp"
	for c := 1; c <= 8; c++ { // bank channel
		ch, x := c, bankX+(b-1)*bankDX+(c-1)*chanDX

		n := fmt.Sprintf("%v%vsolo", pre, ch)
		log.Println("NewOutput() element[%v]:", n)
		widgets[n] = NewToggle(x, soloY, tinySwitch, false)

		n = fmt.Sprintf("%v%vmeter", pre, ch)
		log.Println("NewOutput() element[%v]:", n)
		widgets[n] = &Meter{
			pos:  image.Point{x, meterY},
			size: smallVMeter,
		}
	}

	return &Page{widgets}
}

// Widget returns the named widget.
func (w *Page) Widget(n string) Widget {
	return w.widgets[n]
}
