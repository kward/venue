package presets

import (
	"fmt"

	pb "github.com/kward/venue/presets/proto"
	log "github.com/sirupsen/logrus"
)

//-----------------------------------------------------------------------------
// Body > MicLineStrips

var _ Adjuster = new(MicLineStrips)

type MicLineStrips struct {
	pb.DShowInputChannel_MicLineStrips
	size int
}

// micLineStripsFields defines the complete schema for MicLineStrips.
func micLineStripsFields(m *MicLineStrips) []FieldDescriptor {
	return []FieldDescriptor{
		{"hpfIn", 0x00, FieldTypeBool, &m.HpfIn},
		{"hpf", 0x01, FieldTypeInt32LE, &m.Hpf},
		{"eq_dyn", 0x05, FieldTypeEqDyn, &m.EqDyn},
		{"comp_lim_in", 0x06, FieldTypeBool, &m.CompLimIn},
		{"comp_lim_threshold_db", 0x07, FieldTypeFloat32x10, &m.CompLimThreshold},
		{"comp_lim_ratio", 0x0b, FieldTypeFloat32x100, &m.CompLimRatio},
		{"comp_lim_attack_us", 0x0f, FieldTypeInt32LE, &m.CompLimAttack},
		{"comp_lim_release_ms", 0x13, FieldTypeInt32LE, &m.CompLimRelease},
		{"comp_lim_knee", 0x17, FieldTypeInt32LE, &m.CompLimKnee},
		{"comp_lim_gain_db", 0x1a, FieldTypeFloat32x10, &m.CompLimGain},
		{"exp_gate_in", 0x2f, FieldTypeBool, &m.ExpGateIn},
		{"exp_gate_threshold_db", 0x30, FieldTypeFloat32x10, &m.ExpGateThreshold},
		{"exp_gate_attack", 0x34, FieldTypeFloat32x10, &m.ExpGateAttack},
		{"exp_gate_ratio", 0x38, FieldTypeFloat32x10, &m.ExpGateRatio},
		{"exp_gate_release_ms", 0x3c, FieldTypeInt32LE, &m.ExpGateRelease},
		{"exp_gate_hold_ms", 0x40, FieldTypeInt32LE, &m.ExpGateHold},
		{"exp_gate_range_db", 0x44, FieldTypeFloat32x10, &m.ExpGateRange},
		{"exp_gate_sidechain_in", 0x52, FieldTypeBool, &m.ExpGateSidechainIn},
		{"lpfIn", 0x58, FieldTypeBool, &m.LpfIn},
		{"lpf", 0x59, FieldTypeInt32LE, &m.Lpf},
	}
}

func NewMicLineStrips() *MicLineStrips {
	return &MicLineStrips{}
}

// Read Adjuster values from a slice of bytes.
func (p *MicLineStrips) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	p.size = len(bs)

	r := NewBinaryReader(bs)
	for _, fd := range micLineStripsFields(p) {
		ReadField(r, fd)
	}
	if r.Err() != nil {
		return 0, r.Err()
	}
	return len(bs), nil
}

const micLineStripsSize = 0x65

// Marshal the Adjuster into a slice of bytes.
func (p *MicLineStrips) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	if p.size == 0 {
		p.size = micLineStripsSize
	}

	w := NewBinaryWriter(p.size)
	for _, fd := range micLineStripsFields(p) {
		WriteField(w, fd)
	}
	return w.Bytes(), nil
}

func (p *MicLineStrips) Name() string { return "MicLineStrips" }
func (p *MicLineStrips) String() string {
	var s string
	s += fmt.Sprintf("%s\n", p.Name())
	for _, fd := range micLineStripsFields(p) {
		s += fmt.Sprintf(" %s: ", fd.Name)
		switch fd.FieldType {
		case FieldTypeBool:
			s += fmt.Sprintf("%v\n", *fd.Ptr.(*bool))
		case FieldTypeInt32LE:
			s += fmt.Sprintf("%d\n", *fd.Ptr.(*int32))
		case FieldTypeFloat32x10, FieldTypeFloat32x100:
			s += fmt.Sprintf("%.2f\n", *fd.Ptr.(*float32))
		case FieldTypeEqDyn:
			s += fmt.Sprintf("%v\n", *fd.Ptr.(*pb.DShowInputChannelEqDyn))
		}
	}
	return s
}
