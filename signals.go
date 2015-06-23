/*
This file contains UI elements representing the various signals.
*/
package venue

import ()

// Input represents an input signal.
type Input struct {
	v    *Venue
	ch   int // Channel number
	prop map[string]*InputProperty
}

func NewInput(v *Venue, ch int) (input *Input) {
	const (
		auxMin = -999.0
		auxMax = 12.0
		auxDef = -20.0
		panMin = -100.0
		panMax = 100.0
		panDef = 0.0
	)

	input = &Input{
		v:  v,
		ch: ch,
		prop: map[string]*InputProperty{
			"gain":       NewInputProperty(10, 10, 60, 1, "dB", true),
			"delay":      NewInputProperty(0, 0, 250, 0, "ms", false),
			"hpf":        NewInputProperty(100, 20, 500, 0, "Hz", true),
			"pan":        NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux1":       NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux12pan":   NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux3":       NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux34pan":   NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux5":       NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux56pan":   NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux7":       NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux78pan":   NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux9":       NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux910pan":  NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux11":      NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux1112pan": NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux13":      NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux1314pan": NewInputProperty(panDef, panMin, panMax, 0, "", false),
			"aux15":      NewInputProperty(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux1516pan": NewInputProperty(panDef, panMin, panMax, 0, "", false),
		},
	}
	input.Reset()
	return
}

func (i *Input) Reset() {
	for _, prop := range i.prop {
		prop.Reset()
	}
}

type InputProperty struct {
	Val, defVal, min, max float32 // Value
	prec                  int     // Precision
	unit                  string  // Measurement unit
	Ena, defEna           bool    // Enabled
}

func NewInputProperty(val, min, max float32, prec int, unit string, ena bool) *InputProperty {
	return &InputProperty{val, val, min, max, prec, unit, ena, ena}
}

func (p *InputProperty) Reset() {
	p.Val = p.defVal
}

// Output represents an input signal.
type Output struct {
	v    *Venue
	ch   int // Channel number
	prop map[string]*InputProperty
}
