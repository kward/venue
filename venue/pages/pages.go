package pages

// Page indicates the type of UI page.
type Page int

const (
	// Inputs refers to the Venue INPUTS page. This page manages input signals.
	Inputs Page = 0
	// Outputs refers to the Venue OUTPUTS page. This page manages output signals.
	Outputs Page = 1
	// Filing refers to the Venue FILING page. This page supports loading,
	// saving, and transferring of data.
	Filing Page = 2
	// Snapshots refers to the Venue SNAPSHOTS page. This page manages the scene
	// snapshots.
	Snapshots Page = 3
	// Patchbay refers to the Venue PATCHBAY page. This page controls all signal
	// routing from inputs to outputs.
	Patchbay Page = 4
	// Plugins refers to the Venue PLUG-INS page. This page manages the various
	// signal processing plugins.
	Plugins Page = 5
	// Options refers to the Venue OPTIONS page. This page manages the various
	// user definable system options.
	Options Page = 6
)
