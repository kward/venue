package types

// Typ represents the type of preset data token.
type Typ int

//go:generate stringer -type=Typ

// Typ constants.
const (
	Unknown   Typ = iota
	NLCString     // Newline C String
	CString       // C String
	Bytes         // Bytes
)
