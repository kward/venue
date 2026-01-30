// Package controls defines the supported console controls.
package controls

// Control identifies the type of control.
type Control int

//go:generate stringer -type=Control

const (
	// Unknown indicates a control wasn't specified.
	Unknown Control = iota

	// -- Common --

	// Mute en-/disables the mute for a channel.
	Mute
	// Select the signal to operate on. Frequently a 'Select' button.
	Select
	// Solo en-/disables the solo for a channel.
	Solo
	// SoloClear clears the current channel solo.
	SoloClear

	// -- Inputs --

	Delay
	Fader
	Gain
	Guess
	HPF
	Pad
	Pan
	Phantom
	Phase

	// -- Outputs --

	Aux
	AuxPan
	Group
	GroupPan
	VarGroups
)
