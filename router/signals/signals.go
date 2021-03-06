// Package signals defines the supported types of console signals.
package signals

// Signal identifies the type of signal.
type Signal int

//go:generate stringer -type=Signal

const (
	// Unknown indicates the action wasn't specified.
	Unknown Signal = iota

	// Input signals.
	Input
	FXReturn

	// Output signals.
	Direct
	Insert
	Aux
	Group
	VarGroup
	Matrix

	// Zero-based channel index: FL, FR, C LFE, SL, SR, SBL, SBR
	// https://en.wikipedia.org/wiki/Surround_sound#Channel_identification
	Mains
)

// SignalNo is the signal number.
type SignalNo int
