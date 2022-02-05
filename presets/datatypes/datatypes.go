package datatypes

import (
	"bytes"
	"encoding/binary"
	log "github.com/sirupsen/logrus"
)

//go:generate stringer -type=DataType

// DataTypes
type DataType byte

// DataType constants.
const (
	Int32      DataType = 0x06
	String     DataType = 0x0a // A null terminated string.
	TokenCount DataType = 0x0b // Indicator for how many tokens will follow.
	Bytes      DataType = 0x0d
	Invalid    DataType = 0xff // Chosen out of thin air. Might break in future.
)

func ReadDataType(bs []byte, offset int) (DataType, int) {
	log.Debugf("ReadDataType()")
	if len(bs) < offset+1 {
		log.Tracef(" datatype = %s (0x%02x)", Invalid, byte(Invalid))
		return Invalid, 1
	}
	dt := DataType(bs[offset])
	log.Tracef(" datatype = %s (0x%02x)", dt, byte(dt))
	return dt, 1
}

func WriteBytes(v []byte) []byte {
	log.Debugf("WriteBytes(v)")
	// log.Tracef(" bytes: %v", v)
	bs := []byte{byte(Bytes)}
	bs = append(bs, writeInt32(int32(len(v)))...)
	return append(bs, v...)
}

func WriteInt32(v int32) []byte {
	log.Debugf("WriteInt32(%d)", v)
	return append([]byte{byte(Int32)}, writeInt32(v)...)
}

func writeInt32(v int32) []byte {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
		return []byte{0xff, 0xff, 0xff, 0xff}
	}
	return buf.Bytes()
}

func WriteString(v string) []byte {
	log.Debugf("WriteString(%s)", v)
	bs := append([]byte{0x0a}, v...)
	return append(bs, []byte{0x00}...)
}

func WriteTokenBytes(t string, v []byte) []byte {
	log.Debugf("WriteTokenBytes(%s, v)", t)
	// log.Tracef(" bytes: %v", v)
	return append(WriteString(t), WriteBytes(v)...)
}

func WriteTokenCount(v int32) []byte {
	log.Debugf("WriteTokenCount(%d)", v)
	bs := WriteInt32(v)
	bs[0] = byte(TokenCount)
	return bs
}

func WriteTokenInt32(t string, v int32) []byte {
	log.Debugf("WriteTokenInt32(%s, %d)", t, v)
	return append(WriteString(t), WriteInt32(v)...)
}

func WriteTokenString(t string, v string) []byte {
	log.Debugf("WriteTokenString(%s, %s)", t, v)
	return append(WriteString(t), WriteString(v)...)
}
