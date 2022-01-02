package datatypes

//go:generate stringer -type=DataType

// DataTypes
type DataType byte

// DataType constants.
const (
	UInt32 DataType = 0x06
	Token  DataType = 0x0a // A null terminated string.
	Bytes  DataType = 0x0b
	String DataType = 0x0d // A null prefixed and terminated string.
)
