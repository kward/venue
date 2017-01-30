package signals

// Signal identifies the type of signal.
type Signal int

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
