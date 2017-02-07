// Package encoders defines the supported types of encoders.
package encoders

// Encoder represents the encoder position.
type Encoder int

//go:generate stringer -type=Encoder

const (
	TopLeft Encoder = iota
	TopRight
	MiddleLeft
	MiddleRight
	BottomLeft
	BottomCenter
	BottomRight
)
