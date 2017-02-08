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
	// InputGain sets the input gain of a channel.
	InputGain

	// SelectOutput channel for adjustment.
	SelectOutput
	// OutputLevel sets the output level of a channel.
	OutputLevel
)
