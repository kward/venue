package actions

// Action identifies the type of action.
type Action int

const (
	// Unknown indicates the action wasn't specified.
	Unknown Action = iota

	// Ping is a periodic request to indicate the client is still alive.
	Ping
	// DropPacket is a special request that indicates the packet should be
	// dropped. This action is given instead of nil to differentiate it from an
	// unknown error condition.
	DropPacket

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
