// Package controls defines the supported console controls.
package controls

// Control identifies the type of control.
type Control int

//go:generate stringer -type=Control

const (
	// Unknown indicates a control wasn't specified.
	Unknown Control = iota

	// Select the signal to operate on. Frequently a 'Select' button.
	Select

	Aux
	AuxPan
	Delay
	Fader
	Gain
	Group
	GroupPan
	HPF
	Pan
	Solo
	SoloClear
	VarGroups
)
