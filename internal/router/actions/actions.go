// Package actions defines the supported endpoint actions.
package actions

// Action identifies the type of action.
type Action int

//go:generate stringer -type=Action

const (
	// Unknown indicates the action wasn't specified.
	Unknown Action = iota

	// Noop is a special request that indicates the packet should be ignored.
	// This action is given instead of nil to differentiate it from an unknown
	// error condition.
	Noop
	// Ping is a periodic request to indicate the client is still alive.
	Ping

	// SelectInput channel for adjustment.
	SelectInput
	// InputBank selects another input bank.
	InputBank
	// InputGain sets the gain of an input channel.
	InputGain
	// InputGuess guesses the gain level of an input channel.
	InputGuess
	// InputMute toggles the state of the input mute button.
	InputMute
	// InputSolo toggles the state of the input solo button.
	InputSolo
	// InputPad toggles the state of the 20 dB Pad.
	InputPad
	// InputPhantom toggles the state of the 48V phantom button.
	InputPhantom

	// SelectOutput channel for adjustment.
	SelectOutput
	// OutputLevel sets the level of an output channel.
	OutputLevel
)
