/*
The console represents the full known state of the console, including all
signals that are not currently exposed by the UI.
*/
package venue

import (
	"math"

	"github.com/kward/venue/router/signals"
)

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

// Input represents an input signal.
type Input struct {
	sig   signals.Signal
	sigNo signals.SignalNo
	prop  Signals
	sends Signals
}

func NewInput(sig signals.Signal, sigNo signals.SignalNo) *Input {
	i := &Input{
		sig:   sig,
		sigNo: sigNo,
		prop: Signals{
			"Fader": NewSignal(math.Inf(-1), math.Inf(-1), 15, 1, "dB", true),
			"Gain":  NewSignal(10, 10, 60, 1, "dB", true),
			"Delay": NewSignal(0, 0, 250, 0, "ms", false),
			"HPF":   NewSignal(100, 20, 500, 0, "Hz", true),
		},
		sends: Signals{
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

// Output represents an input signal.
type Output struct {
	sig   signals.Signal
	sigNo signals.SignalNo
	prop  Signals
	sends Signals
}

func NewOutput(sig signals.Signal, sigNo signals.SignalNo) *Output {
	o := &Output{
		sig:   sig,
		sigNo: sigNo,
		prop: Signals{
			"Fader": NewSignal(math.Inf(-1), math.Inf(-1), 15, 1, "dB", true),
		},
		sends: Signals{
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
