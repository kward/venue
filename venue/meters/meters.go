// Package meters defines the supported types of meters.
package meters

type Meter int

//go:generate stringer -type=Meter

const (
	SmallVertical    Meter = iota // Channel (13x50 px)
	MediumHorizontal              // Comp/Lim or Exp/Gate ()
	LargeVertical                 // Input ()
)
