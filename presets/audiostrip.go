package presets

import (
	"fmt"

	pb "github.com/kward/venue/presets/proto"
	log "github.com/sirupsen/logrus"
)

//-----------------------------------------------------------------------------
// Body > AudioStrip

var _ Adjuster = new(AudioStrip)

type AudioStrip struct {
	pb.DShowInputChannel_AudioStrip
}

// audioStripSchema defines the complete schema for AudioStrip.
func audioStripSchema(a *AudioStrip) []FieldDescriptor {
	return []FieldDescriptor{
		{"phaseIn", 0x00, FieldTypeBool, &a.PhaseIn},
		{"delayIn", 0x01, FieldTypeBool, &a.DelayIn},
		{"delay", 0x02, FieldTypeFloat32x96, &a.Delay},
		{"directOutIn", 0x07, FieldTypeBool, &a.DirectOutIn},
		{"directOut", 0x0b, FieldTypeFloat32x10, &a.DirectOut},
		{"pan", 0x11, FieldTypeInt32LE, &a.Pan},
		{"left_right", 0x19, FieldTypeBool, &a.LeftRight},
		{"center_mono", 0x1a, FieldTypeBool, &a.CenterMono},
	}
}

func NewAudioStrip() *AudioStrip {
	return &AudioStrip{}
}

// Read Adjuster values from a slice of bytes.
func (p *AudioStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())

	r := NewBinaryReader(bs)
	for _, fd := range audioStripSchema(p) {
		ReadField(r, fd)
	}
	if r.Err() != nil {
		return 0, r.Err()
	}
	return len(bs), nil
}

const audioStripSize = 0x49

// Marshal the Adjuster into a slice of bytes.
func (p *AudioStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())

	w := NewBinaryWriter(audioStripSize)
	for _, fd := range audioStripSchema(p) {
		WriteField(w, fd)
	}
	return w.Bytes(), nil
}

func (p *AudioStrip) Name() string { return "AudioStrip" }
func (p *AudioStrip) String() string {
	var s string
	s += fmt.Sprintf("%s\n", p.Name())
	for _, fd := range audioStripSchema(p) {
		s += fmt.Sprintf(" %s: ", fd.Name)
		switch fd.FieldType {
		case FieldTypeBool:
			s += fmt.Sprintf("%v\n", *fd.Ptr.(*bool))
		case FieldTypeInt32LE:
			s += fmt.Sprintf("%d\n", *fd.Ptr.(*int32))
		case FieldTypeFloat32x10, FieldTypeFloat32x96:
			s += fmt.Sprintf("%.1f\n", *fd.Ptr.(*float32))
		}
	}
	return s
}
