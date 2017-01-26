package controls

// Control identifies the type of control.
type Control int

const (
	// Unknown indicates a control wasn't specified.
	Unknown Control = iota

	// Input indicates an input was selected.
	Input

	// Aux indicates an aux was selected.
	Aux

	// Group indicates a group was selected.
	Group
)
