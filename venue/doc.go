/*
Package venue the Avid™ VENUE console as a Venue endpoint.

This package enables programmatic access to an Avid™ VENUE console via the VNC
interface.

The code is broken into multiple parts:
- Endpoint (venue.go) -- standard API with the rest of the Venue program. This
  module contains all programming logic for driving the UI.
- UI (ui*.go) -- interaction with the VENUE interface via VNC. This module has
  the ability to interact with the console via VNC, but has no knowledge or
  understanding past enabling standardized interactions (e.g. pressing a
  button).
- Console (console.go) -- in-memory representation of the console. This module
  maintains the state of the system, but has no ability to influence the state
  of the system.
*/
package venue // import "github.com/kward/venue/venue"
