// Packages switches defines different sizes, kinds, and states of switches.
package switches

// Size defines size of a switch.
type Size int

//go:generate stringer -type=Size

const (
	// Tiny: Channel solo or mute on bank (13x13 px)
	Tiny Size = iota
	// Small: Encoder ON (18x14 px)
	Small
	// Medium: 48V, Pad, Guess (26x16 px)
	Medium
	// Large: Channel solo or mute (32x18)
	Large
)

// Kind defines the kind of a switch.
type Kind int

//go:generate stringer -type=Kind

const (
	Toggle Kind = iota
	PushButton
)

// Toggle switch states.
const (
	// Enabled indicates the switch was pressed. On the console, this is
	// indicated by an LED on the switch, and on the UI, an equivalent LED-like
	// indicator.
	Enabled = true

	// Disabled indicates the switch is not active.
	Disabled = false
)
