package encoders

// Encoder represents the encoder position.
type Encoder int

const (
	TopLeft Encoder = iota
	TopRight
	MiddleLeft
	MiddleRight
	BottomLeft
	BottomCenter
	BottomRight
)
