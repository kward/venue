/*
This file contains UI elements representing the various signals.
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

func (s *Signal) Enabled() bool {
	return s.ena
}

func (s *Signal) Value() float64 {
	return s.val
}

func (s *Signal) Reset() {
	s.val = s.defVal
	s.ena = s.defEna
}

const (
	Ichannel = iota
	Ifx
)

// Input represents an input signal.
type Input struct {
	v     *Venue
	ch    int // Channel number
	kind  int //
	prop  Signals
	sends Signals
}

func NewInput(v *Venue, ch int, kind int) *Input {
	i := &Input{
		v:    v,
		ch:   ch,
		kind: kind,
		prop: map[string]*Signal{
			"signal": NewSignal(math.Inf(-1), math.Inf(-1), 15, 1, "dB", true),
			"gain":   NewSignal(10, 10, 60, 1, "dB", true),
			"delay":  NewSignal(0, 0, 250, 0, "ms", false),
			"hpf":    NewSignal(100, 20, 500, 0, "Hz", true),
		},
		sends: map[string]*Signal{
			"pan":        NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux1":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux2":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux12pan":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux3":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux4":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux34pan":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux5":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux6":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux56pan":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux7":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux8":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux78pan":   NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux9":       NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux10":      NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux910pan":  NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux11":      NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux12":      NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux1112pan": NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux13":      NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux14":      NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux1314pan": NewSignal(panDef, panMin, panMax, 0, "", false),
			"aux15":      NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux16":      NewSignal(auxDef, auxMin, auxMax, 1, "dB", true),
			"aux1516pan": NewSignal(panDef, panMin, panMax, 0, "", false),
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
			"signal": NewSignal(math.Inf(-1), math.Inf(-1), 15, 1, "dB", true),
		},
		sends: map[string]*Signal{
			"pan": NewSignal(panDef, panMin, panMax, 0, "", false),
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
