package meters

type Meter int

const (
	SmallVertical    Meter = iota // Channel (13x50 px)
	MediumHorizontal              // Comp/Lim or Exp/Gate ()
	LargeVertical                 // Input ()
)
