// Package pages defines the supported types of pages.
package pages

// Page indicates the type of UI page.
type Page int

//go:generate stringer -type=Page

const (
	// Inputs refers to the Venue INPUTS page. This page manages input signals.
	Inputs Page = iota
	// Outputs refers to the Venue OUTPUTS page. This page manages output signals.
	Outputs
	// Filing refers to the Venue FILING page. This page supports loading,
	// saving, and transferring of data.
	Filing
	// Snapshots refers to the Venue SNAPSHOTS page. This page manages the scene
	// snapshots.
	Snapshots
	// Patchbay refers to the Venue PATCHBAY page. This page controls all signal
	// routing from inputs to outputs.
	Patchbay
	// Plugins refers to the Venue PLUG-INS page. This page manages the various
	// signal processing plugins.
	Plugins
	// Options refers to the Venue OPTIONS page. This page manages the various
	// user definable system options.
	Options
)
