package switches

type SwitchSize int

// Switch size examples:
// - Tiny: Channel solo or mute on bank (13x13 px)
// - Small: Encoder ON (18x14 px)
// - Medium: 48V, Pad, Guess (26x16 px)
// - Large: Channel solo or mute (32x18)
const ( // Switch size.
	Tiny SwitchSize = iota
	Small
	Medium
	Large
)

type SwitchKind int

const (
	Toggle SwitchKind = iota
	PushButton
)
