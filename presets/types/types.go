package types

//go:generate stringer -type=Typ

// Typ represents the type of preset data token.
type Typ int

// Typ constants.
const (
	Unknown   Typ = iota
	NLCString     // Newline C String
	Bytes         // Bytes
)
