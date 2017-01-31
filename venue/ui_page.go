package venue

import (
	"fmt"
	"image"

	"github.com/golang/glog"
	vnclib "github.com/kward/go-vnc"
	"github.com/kward/venue/router/controls"
	"github.com/kward/venue/venue/pages"
	"github.com/kward/venue/vnc"
)

type Page struct {
	page    pages.Page
	widgets map[string]Widget
}
type Pages map[pages.Page]*Page

// Verify that the expected interface is implemented properly.
var _ Widget = new(Page)

// Read implements the Widget interface.
func (w *Page) Read(v *vnc.VNC) (interface{}, error) {
	return nil, fmt.Errorf("page.Read() is unsupported")
}

// Update implements the Widget interface.
func (w *Page) Update(v *vnc.VNC, val interface{}) error {
	// var p int = val.(int)
	// if err := v.KeyPress(vnclib.KeyF1 + uint32(p)); err != nil {
	// 	return err
	// }
	// return nil
	return fmt.Errorf("page.Update() is unsupported")
}

// Press implements the Widget interface.
func (w *Page) Press(v *vnc.VNC) error {
	if err := v.KeyPress(vnclib.KeyF1 + uint32(w.page)); err != nil {
		return err
	}
	return nil
}

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
	return &Page{
		pages.Inputs,
		Widgets{
			"Gain":         &Encoder{image.Point{167, 279}, encoderBL, true},
			"Delay":        &Encoder{image.Point{168, 387}, encoderBL, false},
			"HPF":          &Encoder{image.Point{168, 454}, encoderBL, true},
			"Pan":          &Encoder{image.Point{239, 443}, encoderBC, false},
			"VarGroups":    NewPushButton(226, 299, mediumSwitch),
			"Aux 1":        &Encoder{image.Point{auxOddX, aux12Y}, encoderTR, true},
			"AuxPan 1/2":   &Encoder{image.Point{auxPanX, aux12Y}, encoderTL, false},
			"Aux 3":        &Encoder{image.Point{auxOddX, aux34Y}, encoderTR, true},
			"AuxPan 3/4":   &Encoder{image.Point{auxPanX, aux34Y}, encoderTL, false},
			"Aux 5":        &Encoder{image.Point{auxOddX, aux56Y}, encoderTR, true},
			"AuxPan 5/6":   &Encoder{image.Point{auxPanX, aux56Y}, encoderTL, false},
			"Aux 7":        &Encoder{image.Point{auxOddX, aux78Y}, encoderTR, true},
			"AuxPan 7/8":   &Encoder{image.Point{auxPanX, aux78Y}, encoderTL, false},
			"Aux 9":        &Encoder{image.Point{auxOddX, aux910Y}, encoderTR, true},
			"AuxPan 9/10":  &Encoder{image.Point{auxPanX, aux910Y}, encoderTL, false},
			"Aux 11":       &Encoder{image.Point{auxOddX, aux1112Y}, encoderTR, true},
			"AuxPan 11/12": &Encoder{image.Point{auxPanX, aux1112Y}, encoderTL, false},
			"Aux 13":       &Encoder{image.Point{auxOddX, aux1314Y}, encoderTR, true},
			"AuxPan 13/14": &Encoder{image.Point{auxPanX, aux1314Y}, encoderTL, false},
			"Aux 15":       &Encoder{image.Point{auxOddX, aux1516Y}, encoderTR, true},
			"AuxPan 15/16": &Encoder{image.Point{auxPanX, aux1516Y}, encoderTL, false},
			"Group 1":      &Encoder{image.Point{auxOddX, aux12Y}, encoderTR, true},
			"GroupPan 1/2": &Encoder{image.Point{auxPanX, aux12Y}, encoderTL, false},
			"Group 3":      &Encoder{image.Point{auxOddX, aux34Y}, encoderTR, true},
			"GroupPan 3/4": &Encoder{image.Point{auxPanX, aux34Y}, encoderTL, false},
			"Group 5":      &Encoder{image.Point{auxOddX, aux56Y}, encoderTR, true},
			"GroupPan 5/6": &Encoder{image.Point{auxPanX, aux56Y}, encoderTL, false},
			"Group 7":      &Encoder{image.Point{auxOddX, aux78Y}, encoderTR, true},
			"GroupPan 7/8": &Encoder{image.Point{auxPanX, aux78Y}, encoderTL, false},
			"SoloClear":    NewPushButton(979, 493, mediumSwitch),
		}}
}

// NewOutputsPage returns a populated Outputs page.
func NewOutputsPage() *Page {
	widgets := Widgets{
		"SoloClear": NewPushButton(980, 490, mediumSwitch),
	}

	// Auxes
	for _, b := range []int{1, 2} { // Bank.
		pre := controls.Aux.String()
		for c := 1; c <= 8; c++ { // Bank channel.
			ch, x := (b-1)*8+c, bankX+(b-1)*bankDX+(c-1)*chanDX

			n := fmt.Sprintf("%s %d Solo", pre, ch)
			if glog.V(4) {
				glog.Infof("NewOutput() element[%v]:", n)
			}
			widgets[n] = NewToggle(x, soloY, tinySwitch, false)

			n = fmt.Sprintf("%s %d Value", pre, ch)
			if glog.V(4) {
				glog.Infof("NewOutput() element[%v]:", n)
			}
			widgets[n] = &Meter{
				pos:  image.Point{x, meterY},
				size: smallVMeter,
			}
		}
	}

	// Groups
	b := 5 // bank
	pre := controls.Group.String()
	for c := 1; c <= 8; c++ { // bank channel
		ch, x := c, bankX+(b-1)*bankDX+(c-1)*chanDX

		n := fmt.Sprintf("%s %d Solo", pre, ch)
		if glog.V(4) {
			glog.Infof("NewOutput() element[%v]:", n)
		}
		widgets[n] = NewToggle(x, soloY, tinySwitch, false)

		n = fmt.Sprintf("%s %d Meter", pre, ch)
		if glog.V(4) {
			glog.Infof("NewOutput() element[%v]:", n)
		}
		widgets[n] = &Meter{
			pos:  image.Point{x, meterY},
			size: smallVMeter,
		}
	}

	return &Page{pages.Outputs, widgets}
}

// Widget returns the named widget.
func (w *Page) Widget(n string) (Widget, error) {
	v, ok := w.widgets[n]
	if !ok {
		return nil, fmt.Errorf("invalid page widget %q", n)
	}
	return v, nil
}
