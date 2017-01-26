package commands

// Command identifies the type of command.
type Command int

const (
	// Unknown indicates a command wasn't specified.
	Unknown Command = iota

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
