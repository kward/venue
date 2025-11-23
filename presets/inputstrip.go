package presets

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	pb "github.com/kward/venue/presets/proto"
)

// Body > Input Strip

var _ Adjuster = new(InputStrip)

type InputStrip struct {
	pb.DShowInputChannel_InputStrip
}

const (
	EQShelf = false
	EQCurve = true
)

// inputStripSchema returns the field descriptors for InputStrip.
func inputStripSchema(a *InputStrip) []FieldDescriptor {
	return []FieldDescriptor{
		// Input controls
		{"patched", 0x00, FieldTypeBool, &a.Patched},
		{"phantom", 0x01, FieldTypeBool, &a.Phantom},
		{"pad", 0x02, FieldTypeBool, &a.Pad},
		{"gain_db", 0x03, FieldTypeFloat32x10, &a.Gain},
		{"input_direct", 0x0d, FieldTypeBool, &a.InputDirect},

		// EQ controls
		{"eq_in", 0x0e, FieldTypeBool, &a.EqIn},
		{"analog_eq", 0x0f, FieldTypeBool, &a.AnalogEq},

		// EQ High band
		{"eq_high_in", 0x10, FieldTypeBool, &a.EqHighIn},
		{"eq_high_type", 0x11, FieldTypeEqType, &a.EqHighType},
		{"eq_high_gain_db", 0x12, FieldTypeFloat32x10, &a.EqHighGain},
		{"eq_high_freq_hz", 0x16, FieldTypeInt32LE, &a.EqHighFreq},
		{"eq_high_q_bw", 0x1a, FieldTypeFloat32x100, &a.EqHighQBw},

		// EQ High-Mid band
		{"eq_high_mid_in", 0x1e, FieldTypeBool, &a.EqHighMidIn},
		{"eq_high_mid_type", 0x1f, FieldTypeEqType, &a.EqHighMidType},
		{"eq_high_mid_gain_db", 0x20, FieldTypeFloat32x10, &a.EqHighMidGain},
		{"eq_high_mid_freq_hz", 0x24, FieldTypeInt32LE, &a.EqHighMidFreq},
		{"eq_high_mid_q_bw", 0x28, FieldTypeFloat32x100, &a.EqHighMidQBw},

		// EQ Low-Mid band
		{"eq_low_mid_in", 0x2c, FieldTypeBool, &a.EqLowMidIn},
		{"eq_low_mid_type", 0x2d, FieldTypeEqType, &a.EqLowMidType},
		{"eq_low_mid_gain_db", 0x2e, FieldTypeFloat32x10, &a.EqLowMidGain},
		{"eq_low_mid_freq_hz", 0x32, FieldTypeInt32LE, &a.EqLowMidFreq},
		{"eq_low_mid_q_bw", 0x36, FieldTypeFloat32x100, &a.EqLowMidQBw},

		// EQ Low band
		{"eq_low_in", 0x3a, FieldTypeBool, &a.EqLowIn},
		{"eq_low_type", 0x3b, FieldTypeEqType, &a.EqLowType},
		{"eq_low_gain_db", 0x3c, FieldTypeFloat32x10, &a.EqLowGain},
		{"eq_low_freq_hz", 0x40, FieldTypeInt32LE, &a.EqLowFreq},
		{"eq_low_q_bw", 0x44, FieldTypeFloat32x100, &a.EqLowQBw},

		// Bus routing
		{"bus_1", 0x48, FieldTypeBool, &a.Bus1},
		{"bus_2", 0x4a, FieldTypeBool, &a.Bus2},
		{"bus_3", 0x4c, FieldTypeBool, &a.Bus3},
		{"bus_4", 0x4e, FieldTypeBool, &a.Bus4},
		{"bus_5", 0x50, FieldTypeBool, &a.Bus5},
		{"bus_6", 0x52, FieldTypeBool, &a.Bus6},
		{"bus_7", 0x54, FieldTypeBool, &a.Bus7},
		{"bus_8", 0x56, FieldTypeBool, &a.Bus8},

		// Aux channels 1-2
		{"aux_1_in", 0x5c, FieldTypeBool, &a.Aux1In},
		{"aux_1_pre", 0x5d, FieldTypeBool, &a.Aux1Pre},
		{"aux_1_level_db", 0x5e, FieldTypeFloat32x10, &a.Aux1Level},
		{"aux_2_in", 0x63, FieldTypeBool, &a.Aux2In},
		{"aux_2_pre", 0x64, FieldTypeBool, &a.Aux2Pre},
		{"aux_2_level_db", 0x65, FieldTypeFloat32x10, &a.Aux2Level},
		{"aux_1_pan", 0x69, FieldTypeInt32LE, &a.Aux1Pan},

		// Aux channels 3-4
		{"aux_3_in", 0x6d, FieldTypeBool, &a.Aux3In},
		{"aux_3_pre", 0x6e, FieldTypeBool, &a.Aux3Pre},
		{"aux_3_level_db", 0x6f, FieldTypeFloat32x10, &a.Aux3Level},
		{"aux_4_in", 0x74, FieldTypeBool, &a.Aux4In},
		{"aux_4_pre", 0x75, FieldTypeBool, &a.Aux4Pre},
		{"aux_4_level_db", 0x76, FieldTypeFloat32x10, &a.Aux4Level},
		{"aux_3_pan", 0x7a, FieldTypeInt32LE, &a.Aux3Pan},

		// Aux channels 5-6
		{"aux_5_in", 0x7e, FieldTypeBool, &a.Aux5In},
		{"aux_5_pre", 0x7f, FieldTypeBool, &a.Aux5Pre},
		{"aux_5_level_db", 0x80, FieldTypeFloat32x10, &a.Aux5Level},
		{"aux_6_in", 0x85, FieldTypeBool, &a.Aux6In},
		{"aux_6_pre", 0x86, FieldTypeBool, &a.Aux6Pre},
		{"aux_6_level_db", 0x87, FieldTypeFloat32x10, &a.Aux6Level},
		{"aux_5_pan", 0x8b, FieldTypeInt32LE, &a.Aux5Pan},

		// Aux channels 7-8
		{"aux_7_in", 0x8f, FieldTypeBool, &a.Aux7In},
		{"aux_7_pre", 0x90, FieldTypeBool, &a.Aux7Pre},
		{"aux_7_level_db", 0x91, FieldTypeFloat32x10, &a.Aux7Level},
		{"aux_8_in", 0x96, FieldTypeBool, &a.Aux8In},
		{"aux_8_pre", 0x97, FieldTypeBool, &a.Aux8Pre},
		{"aux_8_level_db", 0x98, FieldTypeFloat32x10, &a.Aux8Level},
		{"aux_7_pan", 0x9c, FieldTypeInt32LE, &a.Aux7Pan},

		// Aux channels 9-10
		{"aux_9_in", 0xa0, FieldTypeBool, &a.Aux9In},
		{"aux_9_pre", 0xa1, FieldTypeBool, &a.Aux9Pre},
		{"aux_9_level_db", 0xa2, FieldTypeFloat32x10, &a.Aux9Level},
		{"aux_10_in", 0xa7, FieldTypeBool, &a.Aux10In},
		{"aux_10_pre", 0xa8, FieldTypeBool, &a.Aux10Pre},
		{"aux_10_level_db", 0xa9, FieldTypeFloat32x10, &a.Aux10Level},
		{"aux_9_pan", 0xad, FieldTypeInt32LE, &a.Aux9Pan},

		// Aux channels 11-12
		{"aux_11_in", 0xb1, FieldTypeBool, &a.Aux11In},
		{"aux_11_pre", 0xb2, FieldTypeBool, &a.Aux11Pre},
		{"aux_11_level_db", 0xb3, FieldTypeFloat32x10, &a.Aux11Level},
		{"aux_12_in", 0xb8, FieldTypeBool, &a.Aux12In},
		{"aux_12_pre", 0xb9, FieldTypeBool, &a.Aux12Pre},
		{"aux_12_level_db", 0xba, FieldTypeFloat32x10, &a.Aux12Level},
		{"aux_11_pan", 0xbe, FieldTypeInt32LE, &a.Aux11Pan},
	}
}

// NewInputStrip creates a new InputStrip with default values.
func NewInputStrip() *InputStrip {
	a := &InputStrip{}
	// Set default values.
	a.EqHighMidType = pb.DShowInputChannel_EQ_CURVE
	a.EqLowMidType = pb.DShowInputChannel_EQ_CURVE
	return a
}

// Read Adjuster values from a slice of bytes.
func (p *InputStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	r := NewBinaryReader(bs)

	for _, fd := range inputStripSchema(p) {
		ReadField(r, fd)
	}
	if r.Err() != nil {
		return 0, r.Err()
	}
	return len(bs), nil
}

const inputStripSize = 0x0128 // Old style preset.

// Marshal the Adjuster into a slice of bytes.
func (p *InputStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	w := NewBinaryWriter(inputStripSize)

	for _, fd := range inputStripSchema(p) {
		WriteField(w, fd)
	}
	return w.Bytes(), nil
}

func (p *InputStrip) Name() string { return "InputStrip" }

func (p *InputStrip) String() string {
	var s string
	s += fmt.Sprintf("%s\n", p.Name())
	for _, fd := range inputStripSchema(p) {
		s += fmt.Sprintf(" %s: ", fd.Name)
		switch fd.FieldType {
		case FieldTypeBool:
			s += fmt.Sprintf("%v\n", *fd.Ptr.(*bool))
		case FieldTypeInt32LE:
			s += fmt.Sprintf("%d\n", *fd.Ptr.(*int32))
		case FieldTypeFloat32x10:
			s += fmt.Sprintf("%.1f\n", *fd.Ptr.(*float32))
		case FieldTypeFloat32x100:
			s += fmt.Sprintf("%.2f\n", *fd.Ptr.(*float32))
		case FieldTypeEqType:
			s += fmt.Sprintf("%v\n", *fd.Ptr.(*pb.DShowInputChannelEqType))
		}
	}
	return s
}
