package datatypes

//go:generate stringer -type=DataType

// DataTypes
type DataType byte

// DataType constants.
const (
	Int32      DataType = 0x06
	Token      DataType = 0x0a // A null terminated string.
	TokenCount DataType = 0x0b // Indicator for how many tokens will follow.
	Bytes      DataType = 0x0d
)
