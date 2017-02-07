// Package messages defines the VNC Workflow messages used by Venue.
package messages // import "github.com/kward/venue/vnc/messages"

type Message int

//go:generate stringer -type=Message

const (
	Unknown Message = iota

	// Client-to-Server
	SetPixelFormat
	SetEncodings
	FramebufferUpdateRequest
	KeyEvent
	PointerEvent
	ClientCutText

	// Server-to-Client
	FramebufferUpdate
	SetColorMapEntries
	Bell
	ServerCutText

	// Non-VNC

	// Sleep enables a workflow to pause for a given duration.
	Sleep
)
