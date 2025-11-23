package presets

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"

	pb "github.com/kward/venue/presets/proto"
)

// FieldType represents the data type of a field in binary data.
type FieldType int

const (
	FieldTypeBool FieldType = iota
	FieldTypeInt32LE
	FieldTypeFloat32x10  // Float stored as int32 * 10
	FieldTypeFloat32x96  // Float stored as int32 * 96 (special delay case)
	FieldTypeFloat32x100 // Float stored as int32 * 100
	FieldTypeEqType      // EQ Type enum (bool mapped to enum)
	FieldTypeEqDyn       // EQ/Dyn order enum (bool mapped to enum)
)

// FieldDescriptor contains all information needed to read/write a field.
type FieldDescriptor struct {
	Name      string
	Offset    int
	FieldType FieldType
	Ptr       interface{} // Pointer to the field in the struct
}

// BinaryReader provides error-accumulating reader for binary data.
type BinaryReader struct {
	data []byte
	err  error
}

func NewBinaryReader(data []byte) *BinaryReader {
	return &BinaryReader{data: data}
}

func (r *BinaryReader) ReadBool(offset int) bool {
	if r.err != nil || offset >= len(r.data) {
		r.err = fmt.Errorf("at offset 0x%x: read bool: %w", offset, io.ErrUnexpectedEOF)
		return false
	}
	return r.data[offset] != 0
}

func (r *BinaryReader) ReadInt32LE(offset int) int32 {
	if r.err != nil || offset+4 > len(r.data) {
		r.err = fmt.Errorf("at offset 0x%x: read int32: %w", offset, io.ErrUnexpectedEOF)
		return 0
	}
	return int32(binary.LittleEndian.Uint32(r.data[offset : offset+4]))
}

func (r *BinaryReader) ReadFloat32Scaled(offset int, scale float32) float32 {
	i := r.ReadInt32LE(offset)
	return float32(i) / scale
}

func (r *BinaryReader) ReadFloat32Delay(offset int) float32 {
	i := r.ReadInt32LE(offset)
	return float32(math.Trunc(float64(i) / 96))
}

func (r *BinaryReader) ReadEqType(offset int) pb.DShowInputChannelEqType {
	b := r.ReadBool(offset)
	if b {
		return pb.DShowInputChannel_EQ_CURVE
	}
	return pb.DShowInputChannel_EQ_SHELF
}

func (r *BinaryReader) ReadEqDyn(offset int) pb.DShowInputChannelEqDyn {
	b := r.ReadBool(offset)
	if b {
		return pb.DShowInputChannel_EQ_POST_DYN
	}
	return pb.DShowInputChannel_EQ_PRE_DYN
}

func (r *BinaryReader) Err() error {
	return r.err
}

// ReadField reads a single field based on its descriptor.
func ReadField(r *BinaryReader, fd FieldDescriptor) {
	switch fd.FieldType {
	case FieldTypeBool:
		*fd.Ptr.(*bool) = r.ReadBool(fd.Offset)
	case FieldTypeInt32LE:
		*fd.Ptr.(*int32) = r.ReadInt32LE(fd.Offset)
	case FieldTypeFloat32x10:
		*fd.Ptr.(*float32) = r.ReadFloat32Scaled(fd.Offset, 10)
	case FieldTypeFloat32x96:
		*fd.Ptr.(*float32) = r.ReadFloat32Delay(fd.Offset)
	case FieldTypeFloat32x100:
		*fd.Ptr.(*float32) = r.ReadFloat32Scaled(fd.Offset, 100)
	case FieldTypeEqType:
		*fd.Ptr.(*pb.DShowInputChannelEqType) = r.ReadEqType(fd.Offset)
	case FieldTypeEqDyn:
		*fd.Ptr.(*pb.DShowInputChannelEqDyn) = r.ReadEqDyn(fd.Offset)
	}
}

// BinaryWriter provides stateful writer for binary data.
type BinaryWriter struct {
	data []byte
}

func NewBinaryWriter(size int) *BinaryWriter {
	return &BinaryWriter{data: make([]byte, size)}
}

func (w *BinaryWriter) WriteBool(offset int, v bool) {
	if offset >= len(w.data) {
		return
	}
	if v {
		w.data[offset] = 1
	} else {
		w.data[offset] = 0
	}
}

func (w *BinaryWriter) WriteInt32LE(offset int, v int32) {
	if offset+4 > len(w.data) {
		return
	}
	binary.LittleEndian.PutUint32(w.data[offset:offset+4], uint32(v))
}

func (w *BinaryWriter) WriteFloat32Scaled(offset int, v float32, scale float32) {
	w.WriteInt32LE(offset, int32(v*scale))
}

func (w *BinaryWriter) WriteFloat32Delay(offset int, v float32) {
	w.WriteInt32LE(offset, int32(v)*96)
}

func (w *BinaryWriter) WriteEqType(offset int, v pb.DShowInputChannelEqType) {
	w.WriteBool(offset, v.Number() == 1)
}

func (w *BinaryWriter) WriteEqDyn(offset int, v pb.DShowInputChannelEqDyn) {
	w.WriteBool(offset, v.Number() == 1)
}

func (w *BinaryWriter) Bytes() []byte {
	return w.data
}

// WriteField writes a single field based on its descriptor.
func WriteField(w *BinaryWriter, fd FieldDescriptor) {
	switch fd.FieldType {
	case FieldTypeBool:
		w.WriteBool(fd.Offset, *fd.Ptr.(*bool))
	case FieldTypeInt32LE:
		w.WriteInt32LE(fd.Offset, *fd.Ptr.(*int32))
	case FieldTypeFloat32x10:
		w.WriteFloat32Scaled(fd.Offset, *fd.Ptr.(*float32), 10)
	case FieldTypeFloat32x96:
		w.WriteFloat32Delay(fd.Offset, *fd.Ptr.(*float32))
	case FieldTypeFloat32x100:
		w.WriteFloat32Scaled(fd.Offset, *fd.Ptr.(*float32), 100)
	case FieldTypeEqType:
		w.WriteEqType(fd.Offset, *fd.Ptr.(*pb.DShowInputChannelEqType))
	case FieldTypeEqDyn:
		w.WriteEqDyn(fd.Offset, *fd.Ptr.(*pb.DShowInputChannelEqDyn))
	}
}
