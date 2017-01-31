/*
Everything representing the UI state of the console.
*/
package venue

import "math"

const (
	auxDef = -20.0
	auxMax = 12.0

	panDef = 0.0
	panMin = -100.0
	panMax = 100.0
)

var (
	auxMin = math.Inf(-1)
)

type Signal struct {
	val, defVal, min, max float64 // Value
	prec                  int     // Precision
	unit                  string  // Measurement unit
	ena, defEna           bool    // Enabled
}
type Signals map[string]*Signal

func NewSignal(defVal, min, max float64, prec int, unit string, defEna bool) *Signal {
	return &Signal{
		defVal: defVal,
		min:    min,
		max:    max,
		prec:   prec,
		unit:   unit,
		defEna: defEna,
	}
}

func (sig *Signal) Enabled() bool {
	return sig.ena
}

func (sig *Signal) Value() float64 {
	return sig.val
}

func (sig *Signal) Reset() {
	sig.val = sig.defVal
	sig.ena = sig.defEna
}

type signalEnum int

const (
	signalChannel signalEnum = iota
	signalFx
)

// Input represents an input signal.
type Input struct {
	v     *Venue
	ch    int        // Channel number
	sig   signalEnum //
	prop  Signals
	sends Signals
}

func NewInput(v *Venue, ch int, sig signalEnum) *Input {
	i := &Input{
		v:   v,
		ch:  ch,
		sig: sig,
		prop: map[string]*Signal{
			"Fader": NewSignal(math.Inf(-1), math.Inf(-1), 15, 1, "dB", true),
			"Gain":  NewSignal(10, 10, 60, 1, "dB", true),
			"Delay": NewSignal(0, 0, 250, 0, "ms", false),
			"HPF":   NewSignal(100, 20, 500, 0, "Hz", true),
		},
		sends: map[string]*Signal{
			"Pan":          NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 1":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 2":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 1/2":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 3":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 4":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 3/4":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 5":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 6":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 5/6":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 7":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 8":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 7/8":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 9":        NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 10":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 9/10":  NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 11":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 12":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 11/12": NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 13":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 14":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 13/14": NewSignal(panDef, panMin, panMax, 0, "", false),
			"Aux 15":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"Aux 16":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"AuxPan 15/16": NewSignal(panDef, panMin, panMax, 0, "", false),
		},
	}
	i.Reset()
	return i
}

func (i *Input) Reset() {
	for _, p := range i.prop {
		p.Reset()
	}
	for _, s := range i.sends {
		s.Reset()
	}
}

const (
	Oaux = iota
	Opq
	Omatrix
	Ogroup
	Ovca
	Oleft
	Omiddle
	Oright
	Ogeq
)

// Output represents an input signal.
type Output struct {
	v     *Venue
	name  string
	kind  int
	prop  Signals
	sends Signals
}

func NewOutput(v *Venue, name string, kind int) *Output {
	o := &Output{
		v:    v,
		name: name,
		kind: kind,
		prop: map[string]*Signal{
			"Fader": NewSignal(math.Inf(-1), math.Inf(-1), 15, 1, "dB", true),
		},
		sends: map[string]*Signal{
			"Pan": NewSignal(panDef, panMin, panMax, 0, "", false),
		},
	}
	o.Reset()

	return o
}

func (o *Output) Reset() {
	for _, p := range o.prop {
		p.Reset()
	}
}
