package presets

import (
	"bytes"
	"testing"

	pb "github.com/kward/venue/presets/proto"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.DebugLevel)
	// log.SetLevel(log.TraceLevel)
}

func TestDShowInputChannel(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		setFn func(*DShowInputChannel)
		getFn func(*DShowInputChannel) interface{}
	}{
		{"phantom_true",
			func(dsic *DShowInputChannel) { dsic.Body().InputStrip().Phantom = true },
			func(dsic *DShowInputChannel) interface{} { return dsic.Body().InputStrip().GetPhantom() }},
	} {
		// Marshal the proto to bytes.
		p1 := NewDShowInputChannel()
		tt.setFn(p1)
		value := tt.getFn(p1)
		m1, err := p1.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected Marshal() error %s", tt.desc, err)
			continue
		}

		// Read bytes back.
		p2 := NewDShowInputChannel()
		c, err := p2.Read(m1)
		if err != nil {
			t.Errorf("%s: unexpected Read() error %s", tt.desc, err)
			continue
		}
		if c == 0 {
			t.Errorf("%s: expected count > 0, got %d", tt.desc, c)
			continue
		}
		m2, err := p2.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected Marshal() error %s", tt.desc, err)
			continue
		}

		// Checks.
		if got, want := tt.getFn(p2), value; got != want {
			t.Errorf("%s: got = %v, want %v", tt.desc, got, want)
		}
		if m1len, m2len := len(m1), len(m2); m1len != m2len {
			t.Errorf("%s: len(m1) = %d != len(m2) = %d", tt.desc, m1len, m2len)
		}
	}
}

func TestHeader(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		setFn func(*Header)
		getFn func(*Header) interface{}
	}{
		{"version",
			func(h *Header) { h.Version = int32(1) },
			func(h *Header) interface{} { return h.GetVersion() }},
		{"file_type",
			func(h *Header) { h.FileType = "Super file type" },
			func(h *Header) interface{} { return h.GetVersion() }},
		{"user_comment",
			func(h *Header) { h.UserComment = "Super comment" },
			func(h *Header) interface{} { return h.GetVersion() }},
		{"user_comment_empty",
			func(h *Header) { h.UserComment = "" },
			func(h *Header) interface{} { return h.GetVersion() }},
	} {
		// Marshal the proto to bytes.
		p := NewHeader("Digidesign Storage - 1.0")
		tt.setFn(p)
		value := tt.getFn(p)
		m, err := p.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
			continue
		}

		// Read bytes back.
		p = NewHeader("Digidesign Storage - 1.0")
		c, err := p.Read(m)
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
		}
		if c == 0 {
			t.Errorf("%s: expected count > 0, got %d", tt.desc, c)
		}
		if err != nil {
			continue
		}

		// Verify the value.
		if got, want := tt.getFn(p), value; got != want {
			t.Errorf("%s: got = %v, want %v", tt.desc, got, want)
		}
	}
}

// func TestBody(t *testing.T) {
// 	p = NewBody()
// }

func TestAudioStrip(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		setFn func(*AudioStrip)
		getFn func(*AudioStrip) interface{}
	}{
		{"phase_in_true",
			func(as *AudioStrip) { as.PhaseIn = true },
			func(as *AudioStrip) interface{} { return as.GetPhaseIn() }},
		{"phase_in_false",
			func(as *AudioStrip) { as.PhaseIn = false },
			func(as *AudioStrip) interface{} { return as.GetPhaseIn() }},

		{"delay_in_true",
			func(as *AudioStrip) { as.DelayIn = true },
			func(as *AudioStrip) interface{} { return as.GetDelayIn() }},
		{"phase_false",
			func(as *AudioStrip) { as.DelayIn = false },
			func(as *AudioStrip) interface{} { return as.GetDelayIn() }},

		{"delay_0.0", // Minimum.
			func(as *AudioStrip) { as.Delay = 0.0 },
			func(as *AudioStrip) interface{} { return as.GetDelay() }},
		{"delay_250.0", // Maximum.
			func(as *AudioStrip) { as.Delay = 250.0 },
			func(as *AudioStrip) interface{} { return as.GetDelay() }},

		{"direct_out_in_true",
			func(as *AudioStrip) { as.DirectOutIn = true },
			func(as *AudioStrip) interface{} { return as.GetDirectOutIn() }},
		{"direct_out_in_false",
			func(as *AudioStrip) { as.DirectOutIn = false },
			func(as *AudioStrip) interface{} { return as.GetDirectOutIn() }},

		{"direct_out_-INF", // Minimum.
			func(as *AudioStrip) { as.DirectOut = -103.0 },
			func(as *AudioStrip) interface{} { return as.GetDirectOut() }},
		{"direct_out_+12.0_dB", // Maximum.
			func(as *AudioStrip) { as.DirectOut = 12.0 },
			func(as *AudioStrip) interface{} { return as.GetDirectOut() }},

		{"pan_left",
			func(as *AudioStrip) { as.Pan = -100 },
			func(as *AudioStrip) interface{} { return as.GetPan() }},
		{"pan_center",
			func(as *AudioStrip) { as.Pan = 0 },
			func(as *AudioStrip) interface{} { return as.GetPan() }},
		{"pan_right",
			func(as *AudioStrip) { as.Pan = 100 },
			func(as *AudioStrip) interface{} { return as.GetPan() }},

		{"left_right_true",
			func(as *AudioStrip) { as.LeftRight = true },
			func(as *AudioStrip) interface{} { return as.GetLeftRight() }},
		{"left_right_false",
			func(as *AudioStrip) { as.LeftRight = false },
			func(as *AudioStrip) interface{} { return as.GetLeftRight() }},

		{"center_mono_true",
			func(as *AudioStrip) { as.CenterMono = true },
			func(as *AudioStrip) interface{} { return as.GetCenterMono() }},
		{"center_mono_false",
			func(as *AudioStrip) { as.CenterMono = false },
			func(as *AudioStrip) interface{} { return as.GetCenterMono() }},
	} {
		// Marshal the proto to bytes.
		p := NewAudioStrip()
		tt.setFn(p)
		value := tt.getFn(p)
		m, err := p.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
			continue
		}

		// Read bytes back.
		p = NewAudioStrip()
		c, err := p.Read(m)
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
		}
		if c == 0 {
			t.Errorf("%s: expected count > 0, got %d", tt.desc, c)
		}
		if err != nil {
			continue
		}

		// Verify the value.
		if got, want := tt.getFn(p), value; got != want {
			t.Errorf("%s: got = %v, want %v", tt.desc, got, want)
		}
	}
}

// TestAudioStrip_Delay verifies the marshaled value against a known byte slice.
func TestAudioStrip_Delay(t *testing.T) {
	as := NewAudioStrip()
	as.Delay = 250.0
	bs, err := as.Marshal()
	if err != nil {
		t.Errorf("unexpected error; %s", err)
	}
	o := as.params["delay"].offset
	if got, want := bs[o:o+4], []byte{0xc0, 0x5d, 0x00, 0x00}; !bytes.Equal(got, want) {
		t.Errorf("got = %v, want %v", got, want)
	}
}

func TestInputStrip(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		setFn func(*InputStrip)
		getFn func(*InputStrip) interface{}
	}{
		{"patched_true",
			func(is *InputStrip) { is.Patched = true },
			func(is *InputStrip) interface{} { return is.GetPatched() }},
		{"patched_false",
			func(is *InputStrip) { is.Patched = false },
			func(is *InputStrip) interface{} { return is.GetPatched() }},

		{"phantom_true",
			func(is *InputStrip) { is.Phantom = true },
			func(is *InputStrip) interface{} { return is.GetPhantom() }},
		{"phantom_false",
			func(is *InputStrip) { is.Phantom = false },
			func(is *InputStrip) interface{} { return is.GetPhantom() }},

		{"pad_true",
			func(is *InputStrip) { is.Pad = true },
			func(is *InputStrip) interface{} { return is.GetPad() }},
		{"pad_false",
			func(is *InputStrip) { is.Pad = false },
			func(is *InputStrip) interface{} { return is.GetPad() }},

		{"gain_+10.0_dB", // Minimum.
			func(is *InputStrip) { is.Gain = 10.0 },
			func(is *InputStrip) interface{} { return is.GetGain() }},
		{"gain_+10.1_dB",
			func(is *InputStrip) { is.Gain = 10.1 },
			func(is *InputStrip) interface{} { return is.GetGain() }},
		{"gain_+59.9_dB",
			func(is *InputStrip) { is.Gain = 59.9 },
			func(is *InputStrip) interface{} { return is.GetGain() }},
		{"gain_+60.0_dB", // Maximum.
			func(is *InputStrip) { is.Gain = 60.0 },
			func(is *InputStrip) interface{} { return is.GetGain() }},

		{"input_direct_true",
			func(is *InputStrip) { is.InputDirect = true },
			func(is *InputStrip) interface{} { return is.GetInputDirect() }},
		{"input_direct_false",
			func(is *InputStrip) { is.InputDirect = false },
			func(is *InputStrip) interface{} { return is.GetInputDirect() }},

		{"eq_in_true",
			func(is *InputStrip) { is.EqIn = true },
			func(is *InputStrip) interface{} { return is.GetEqIn() }},
		{"eq_in_false",
			func(is *InputStrip) { is.EqIn = false },
			func(is *InputStrip) interface{} { return is.GetEqIn() }},

		{"analog_eq_true",
			func(is *InputStrip) { is.AnalogEq = true },
			func(is *InputStrip) interface{} { return is.GetAnalogEq() }},
		{"analog_eq_false",
			func(is *InputStrip) { is.AnalogEq = false },
			func(is *InputStrip) interface{} { return is.GetAnalogEq() }},

		// EQ High
		{"eq_high_in_true",
			func(is *InputStrip) { is.EqHighIn = true },
			func(is *InputStrip) interface{} { return is.GetEqHighIn() }},
		{"eq_high_in_false",
			func(is *InputStrip) { is.EqHighIn = false },
			func(is *InputStrip) interface{} { return is.GetEqHighIn() }},
		{"eq_high_type_curve",
			func(is *InputStrip) { is.EqHighType = pb.DShowInputChannel_EQ_CURVE },
			func(is *InputStrip) interface{} { return is.GetEqHighType() }},
		{"eq_high_type_shelf",
			func(is *InputStrip) { is.EqHighType = pb.DShowInputChannel_EQ_SHELF },
			func(is *InputStrip) interface{} { return is.GetEqHighType() }},
		{"eq_high_gain_-18.0_dB", // minimum
			func(is *InputStrip) { is.EqHighGain = -18.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighGain() }},
		{"eq_high_gain_0.0_dB",
			func(is *InputStrip) { is.EqHighGain = 0.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighGain() }},
		{"eq_high_gain_+18.0_dB", // maximum
			func(is *InputStrip) { is.EqHighGain = 18.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighGain() }},
		{"eq_high_freq_20_Hz", // minimum
			func(is *InputStrip) { is.EqHighFreq = 20 },
			func(is *InputStrip) interface{} { return is.GetEqHighFreq() }},
		{"eq_high_freq_20,000_Hz", // maximum
			func(is *InputStrip) { is.EqHighFreq = 20000 },
			func(is *InputStrip) interface{} { return is.GetEqHighFreq() }},
		{"eq_high_q_bw_10", // minimum
			func(is *InputStrip) { is.EqHighQBw = 10.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighQBw() }},
		{"eq_high_q_bw_0.1", // maximum
			func(is *InputStrip) { is.EqHighQBw = 0.1 },
			func(is *InputStrip) interface{} { return is.GetEqHighQBw() }},

		// EQ High Mid
		{"eq_high_mid_in_true",
			func(is *InputStrip) { is.EqHighMidIn = true },
			func(is *InputStrip) interface{} { return is.GetEqHighMidIn() }},
		{"eq_high_mid_in_false",
			func(is *InputStrip) { is.EqHighMidIn = false },
			func(is *InputStrip) interface{} { return is.GetEqHighMidIn() }},
		{"eq_high_mid_type_curve",
			func(is *InputStrip) { is.EqHighMidType = pb.DShowInputChannel_EQ_CURVE },
			func(is *InputStrip) interface{} { return is.GetEqHighMidType() }},
		{"eq_high_mid_type_shelf",
			func(is *InputStrip) { is.EqHighMidType = pb.DShowInputChannel_EQ_SHELF },
			func(is *InputStrip) interface{} { return is.GetEqHighMidType() }},
		{"eq_high_mid_gain_-18.0_dB", // minimum
			func(is *InputStrip) { is.EqHighMidGain = -18.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighMidGain() }},
		{"eq_high_mid_gain_0.0_dB",
			func(is *InputStrip) { is.EqHighMidGain = 0.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighMidGain() }},
		{"eq_high_mid_gain_+18.0_dB", // maximum
			func(is *InputStrip) { is.EqHighMidGain = 18.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighMidGain() }},
		{"eq_high_mid_freq_20_Hz", // minimum
			func(is *InputStrip) { is.EqHighMidFreq = 20 },
			func(is *InputStrip) interface{} { return is.GetEqHighMidFreq() }},
		{"eq_high_mid_freq_20,000_Hz", // maximum
			func(is *InputStrip) { is.EqHighMidFreq = 20000 },
			func(is *InputStrip) interface{} { return is.GetEqHighMidFreq() }},
		{"eq_high_mid_q_bw_10", // minimum
			func(is *InputStrip) { is.EqHighMidQBw = 10.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighMidQBw() }},
		{"eq_high_mid_q_bw_0.1", // maximum
			func(is *InputStrip) { is.EqHighMidQBw = 0.1 },
			func(is *InputStrip) interface{} { return is.GetEqHighMidQBw() }},

		// EQ Low Mid
		{"eq_low_mid_in_true",
			func(is *InputStrip) { is.EqLowMidIn = true },
			func(is *InputStrip) interface{} { return is.GetEqLowMidIn() }},
		{"eq_low_mid_in_false",
			func(is *InputStrip) { is.EqLowMidIn = false },
			func(is *InputStrip) interface{} { return is.GetEqLowMidIn() }},
		{"eq_low_mid_type_curve",
			func(is *InputStrip) { is.EqLowMidType = pb.DShowInputChannel_EQ_CURVE },
			func(is *InputStrip) interface{} { return is.GetEqLowMidType() }},
		{"eq_low_mid_type_shelf",
			func(is *InputStrip) { is.EqLowMidType = pb.DShowInputChannel_EQ_SHELF },
			func(is *InputStrip) interface{} { return is.GetEqLowMidType() }},
		{"eq_low_mid_gain_-18.0_dB", // minimum
			func(is *InputStrip) { is.EqLowMidGain = -18.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowMidGain() }},
		{"eq_low_mid_gain_0.0_dB",
			func(is *InputStrip) { is.EqLowMidGain = 0.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowMidGain() }},
		{"eq_low_mid_gain_+18.0_dB", // maximum
			func(is *InputStrip) { is.EqLowMidGain = 18.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowMidGain() }},
		{"eq_low_mid_freq_20_Hz", // minimum
			func(is *InputStrip) { is.EqLowMidFreq = 20 },
			func(is *InputStrip) interface{} { return is.GetEqLowMidFreq() }},
		{"eq_low_mid_freq_20,000_Hz", // maximum
			func(is *InputStrip) { is.EqLowMidFreq = 20000 },
			func(is *InputStrip) interface{} { return is.GetEqLowMidFreq() }},
		{"eq_low_mid_q_bw_10", // minimum
			func(is *InputStrip) { is.EqLowMidQBw = 10.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowMidQBw() }},
		{"eq_low_mid_q_bw_0.1", // maximum
			func(is *InputStrip) { is.EqLowMidQBw = 0.1 },
			func(is *InputStrip) interface{} { return is.GetEqLowMidQBw() }},

		// EQ Low
		{"eq_low_in_true",
			func(is *InputStrip) { is.EqLowIn = true },
			func(is *InputStrip) interface{} { return is.GetEqLowIn() }},
		{"eq_low_in_false",
			func(is *InputStrip) { is.EqLowIn = false },
			func(is *InputStrip) interface{} { return is.GetEqLowIn() }},
		{"eq_low_type_curve",
			func(is *InputStrip) { is.EqLowType = pb.DShowInputChannel_EQ_CURVE },
			func(is *InputStrip) interface{} { return is.GetEqLowType() }},
		{"eq_low_type_shelf",
			func(is *InputStrip) { is.EqLowType = pb.DShowInputChannel_EQ_SHELF },
			func(is *InputStrip) interface{} { return is.GetEqLowType() }},
		{"eq_low_gain_-18.0_dB", // minimum
			func(is *InputStrip) { is.EqLowGain = -18.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowGain() }},
		{"eq_low_gain_0.0_dB",
			func(is *InputStrip) { is.EqLowGain = 0.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowGain() }},
		{"eq_low_gain_+18.0_dB", // maximum
			func(is *InputStrip) { is.EqLowGain = 18.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowGain() }},
		{"eq_low_freq_20_Hz", // minimum
			func(is *InputStrip) { is.EqLowFreq = 20 },
			func(is *InputStrip) interface{} { return is.GetEqLowFreq() }},
		{"eq_low_freq_20,000_Hz", // maximum
			func(is *InputStrip) { is.EqLowFreq = 20000 },
			func(is *InputStrip) interface{} { return is.GetEqLowFreq() }},
		{"eq_low_q_bw_10", // minimum
			func(is *InputStrip) { is.EqLowQBw = 10.0 },
			func(is *InputStrip) interface{} { return is.GetEqLowQBw() }},
		{"eq_low_q_bw_0.1", // maximum
			func(is *InputStrip) { is.EqLowQBw = 0.1 },
			func(is *InputStrip) interface{} { return is.GetEqLowQBw() }},

		{"bus_1_true",
			func(is *InputStrip) { is.Bus1 = true },
			func(is *InputStrip) interface{} { return is.GetBus1() }},
		{"bus_1_false",
			func(is *InputStrip) { is.Bus1 = false },
			func(is *InputStrip) interface{} { return is.GetBus1() }},
		{"bus_8_true",
			func(is *InputStrip) { is.Bus8 = true },
			func(is *InputStrip) interface{} { return is.GetBus8() }},
		{"bus_8_false",
			func(is *InputStrip) { is.Bus8 = false },
			func(is *InputStrip) interface{} { return is.GetBus8() }},

		// Aux 1 and 2
		{"aux_1_in_true",
			func(is *InputStrip) { is.Aux1In = true },
			func(is *InputStrip) interface{} { return is.GetAux1In() }},
		{"aux_1_in_false",
			func(is *InputStrip) { is.Aux1In = false },
			func(is *InputStrip) interface{} { return is.GetAux1In() }},
		{"aux_1_pre_true",
			func(is *InputStrip) { is.Aux1Pre = true },
			func(is *InputStrip) interface{} { return is.GetAux1Pre() }},
		{"aux_1_pre_false",
			func(is *InputStrip) { is.Aux1Pre = false },
			func(is *InputStrip) interface{} { return is.GetAux1Pre() }},
		{"aux_1_level_-144_dB", // minimum
			func(is *InputStrip) { is.Aux1Level = -144.0 },
			func(is *InputStrip) interface{} { return is.GetAux1Level() }},
		{"aux_1_level_0.0_dB",
			func(is *InputStrip) { is.Aux1Level = 0.0 },
			func(is *InputStrip) interface{} { return is.GetAux1Level() }},
		{"aux_1_level_+12.0_dB", // maximum
			func(is *InputStrip) { is.Aux1Level = 12.0 },
			func(is *InputStrip) interface{} { return is.GetAux1Level() }},

		{"aux_2_in_true",
			func(is *InputStrip) { is.Aux2In = true },
			func(is *InputStrip) interface{} { return is.GetAux2In() }},
		{"aux_2_in_false",
			func(is *InputStrip) { is.Aux2In = false },
			func(is *InputStrip) interface{} { return is.GetAux2In() }},
		{"aux_2_pre_true",
			func(is *InputStrip) { is.Aux2Pre = true },
			func(is *InputStrip) interface{} { return is.GetAux2Pre() }},
		{"aux_2_pre_false",
			func(is *InputStrip) { is.Aux2Pre = false },
			func(is *InputStrip) interface{} { return is.GetAux2Pre() }},
		{"aux_2_level_-144_dB", // minimum
			func(is *InputStrip) { is.Aux2Level = -144.0 },
			func(is *InputStrip) interface{} { return is.GetAux2Level() }},
		{"aux_2_level_0.0_dB",
			func(is *InputStrip) { is.Aux2Level = 0.0 },
			func(is *InputStrip) interface{} { return is.GetAux2Level() }},
		{"aux_2_level_+12.0_dB", // maximum
			func(is *InputStrip) { is.Aux2Level = 12.0 },
			func(is *InputStrip) interface{} { return is.GetAux2Level() }},

		{"aux_1_pan_-100", // left
			func(is *InputStrip) { is.Aux1Pan = -100 },
			func(is *InputStrip) interface{} { return is.GetAux1Pan() }},
		{"aux_1_pan_0", // center
			func(is *InputStrip) { is.Aux1Pan = 0 },
			func(is *InputStrip) interface{} { return is.GetAux1Pan() }},
		{"aux_1_pan_100", // right
			func(is *InputStrip) { is.Aux1Pan = 100 },
			func(is *InputStrip) interface{} { return is.GetAux1Pan() }},

		// Aux 11 and 12
		{"aux_11_in_true",
			func(is *InputStrip) { is.Aux11In = true },
			func(is *InputStrip) interface{} { return is.GetAux11In() }},
		{"aux_11_in_false",
			func(is *InputStrip) { is.Aux11In = false },
			func(is *InputStrip) interface{} { return is.GetAux11In() }},
		{"aux_11_pre_true",
			func(is *InputStrip) { is.Aux11Pre = true },
			func(is *InputStrip) interface{} { return is.GetAux11Pre() }},
		{"aux_11_pre_false",
			func(is *InputStrip) { is.Aux11Pre = false },
			func(is *InputStrip) interface{} { return is.GetAux11Pre() }},
		{"aux_11_level_-144_dB", // minimum
			func(is *InputStrip) { is.Aux11Level = -144.0 },
			func(is *InputStrip) interface{} { return is.GetAux11Level() }},
		{"aux_11_level_0.0_dB",
			func(is *InputStrip) { is.Aux11Level = 0.0 },
			func(is *InputStrip) interface{} { return is.GetAux11Level() }},
		{"aux_11_level_+12.0_dB", // maximum
			func(is *InputStrip) { is.Aux11Level = 12.0 },
			func(is *InputStrip) interface{} { return is.GetAux11Level() }},

		{"aux_12_in_true",
			func(is *InputStrip) { is.Aux12In = true },
			func(is *InputStrip) interface{} { return is.GetAux12In() }},
		{"aux_12_in_false",
			func(is *InputStrip) { is.Aux12In = false },
			func(is *InputStrip) interface{} { return is.GetAux12In() }},
		{"aux_12_pre_true",
			func(is *InputStrip) { is.Aux12Pre = true },
			func(is *InputStrip) interface{} { return is.GetAux12Pre() }},
		{"aux_12_pre_false",
			func(is *InputStrip) { is.Aux12Pre = false },
			func(is *InputStrip) interface{} { return is.GetAux12Pre() }},
		{"aux_12_level_-144_dB", // minimum
			func(is *InputStrip) { is.Aux12Level = -144.0 },
			func(is *InputStrip) interface{} { return is.GetAux12Level() }},
		{"aux_12_level_0.0_dB",
			func(is *InputStrip) { is.Aux12Level = 0.0 },
			func(is *InputStrip) interface{} { return is.GetAux12Level() }},
		{"aux_12_level_+12.0_dB", // maximum
			func(is *InputStrip) { is.Aux12Level = 12.0 },
			func(is *InputStrip) interface{} { return is.GetAux12Level() }},

		{"aux_11_pan_-100", // left
			func(is *InputStrip) { is.Aux11Pan = -100 },
			func(is *InputStrip) interface{} { return is.GetAux11Pan() }},
		{"aux_11_pan_0", // center
			func(is *InputStrip) { is.Aux11Pan = 0 },
			func(is *InputStrip) interface{} { return is.GetAux11Pan() }},
		{"aux_11_pan_100", // right
			func(is *InputStrip) { is.Aux11Pan = 100 },
			func(is *InputStrip) interface{} { return is.GetAux11Pan() }},

		{"aux_count_24", // m
			func(is *InputStrip) { is.Aux11Pan = -100 },
			func(is *InputStrip) interface{} { return is.GetAux11Pan() }},
		{"aux_11_pan_0", // center
			func(is *InputStrip) { is.Aux11Pan = 0 },
			func(is *InputStrip) interface{} { return is.GetAux11Pan() }},

		/*
			{"heat_in_true",
				func(is *InputStrip) { is.HeatIn = true },
				func(is *InputStrip) interface{} { return is.GetHeatIn() }},
			{"heat_in_false",
				func(is *InputStrip) { is.HeatIn = false },
				func(is *InputStrip) interface{} { return is.GetHeatIn() }},
			{"drive_-5", // Minimum.
				func(is *InputStrip) { is.Drive = -5 },
				func(is *InputStrip) interface{} { return is.GetDrive() }},
			{"drive_0",
				func(is *InputStrip) { is.Drive = 0 },
				func(is *InputStrip) interface{} { return is.GetDrive() }},
			{"drive_5", // Maximum.
				func(is *InputStrip) { is.Drive = 5 },
				func(is *InputStrip) interface{} { return is.GetDrive() }},
			{"tone-0", // Minimum.
				func(is *InputStrip) { is.Tone = 0 },
				func(is *InputStrip) interface{} { return is.GetTone() }},
			{"tone-6", // Maximum.
				func(is *InputStrip) { is.Tone = 6 },
				func(is *InputStrip) interface{} { return is.GetTone() }},
		*/
	} {
		// Marshal the proto to bytes.
		p := NewInputStrip()
		tt.setFn(p)
		value := tt.getFn(p)
		m, err := p.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
			continue
		}

		// Read bytes back.
		p = NewInputStrip()
		c, err := p.Read(m)
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
		}
		if c == 0 {
			t.Errorf("%s: expected count > 0, got %d", tt.desc, c)
		}
		if err != nil {
			continue
		}

		// Verify the value.
		if got, want := tt.getFn(p), value; got != want {
			t.Errorf("%s: got = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestMicLineStrips(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		setFn func(*MicLineStrips)
		getFn func(*MicLineStrips) interface{}
	}{
		{"hpf_in_true",
			func(p *MicLineStrips) { p.HpfIn = true },
			func(p *MicLineStrips) interface{} { return p.GetHpfIn() }},
		{"hpf_in_false",
			func(p *MicLineStrips) { p.HpfIn = false },
			func(p *MicLineStrips) interface{} { return p.GetHpfIn() }},

		{"hpf_20",
			func(p *MicLineStrips) { p.Hpf = 20 },
			func(p *MicLineStrips) interface{} { return p.GetHpfIn() }},
		{"hpf_20000",
			func(p *MicLineStrips) { p.Hpf = 20000 },
			func(p *MicLineStrips) interface{} { return p.GetHpfIn() }},

		{"eq_post_dyn",
			func(p *MicLineStrips) { p.EqDyn = pb.DShowInputChannel_EQ_POST_DYN },
			func(p *MicLineStrips) interface{} { return p.GetEqDyn() }},
		{"eq_pre_dyn",
			func(p *MicLineStrips) { p.EqDyn = pb.DShowInputChannel_EQ_PRE_DYN },
			func(p *MicLineStrips) interface{} { return p.GetEqDyn() }},

		{"exp_gate_in_true",
			func(p *MicLineStrips) { p.ExpGateIn = true },
			func(p *MicLineStrips) interface{} { return p.GetExpGateIn() }},
		{"exp_gate_in_false",
			func(p *MicLineStrips) { p.ExpGateIn = false },
			func(p *MicLineStrips) interface{} { return p.GetExpGateIn() }},

		{"lpf_in_true",
			func(p *MicLineStrips) { p.LpfIn = true },
			func(p *MicLineStrips) interface{} { return p.GetLpfIn() }},
		{"lpf_in_false",
			func(p *MicLineStrips) { p.LpfIn = false },
			func(p *MicLineStrips) interface{} { return p.GetLpfIn() }},

		{"lpf_20",
			func(p *MicLineStrips) { p.Lpf = 20 },
			func(p *MicLineStrips) interface{} { return p.GetLpfIn() }},
		{"lpf_20000",
			func(p *MicLineStrips) { p.Lpf = 20000 },
			func(p *MicLineStrips) interface{} { return p.GetLpfIn() }},
	} {
		// Marshal the proto to bytes.
		p := NewMicLineStrips()
		tt.setFn(p)
		value := tt.getFn(p)
		m, err := p.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
			continue
		}

		// Read bytes back.
		p = NewMicLineStrips()
		c, err := p.Read(m)
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
		}
		if c == 0 {
			t.Errorf("%s: expected count > 0, got %d", tt.desc, c)
		}
		if err != nil {
			continue
		}

		// Verify the value.
		if got, want := tt.getFn(p), value; got != want {
			t.Errorf("%s: got = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestStrip(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		setFn func(*Strip)
		getFn func(*Strip) interface{}
	}{
		{"mute_true",
			func(s *Strip) { s.Mute = true },
			func(s *Strip) interface{} { return s.GetMute() }},
		{"mute_false",
			func(s *Strip) { s.Mute = false },
			func(s *Strip) interface{} { return s.GetMute() }},

		{"fader_-INF", // Minimum.
			func(s *Strip) { s.Fader = -131.0 },
			func(s *Strip) interface{} { return s.GetFader() }},
		{"fader_0.0_dB",
			func(s *Strip) { s.Fader = 0.0 },
			func(s *Strip) interface{} { return s.GetFader() }},
		{"fader_+12.0_dB",
			func(s *Strip) { s.Fader = 12.0 },
			func(s *Strip) interface{} { return s.GetFader() }},

		{"channel_name_empty",
			func(s *Strip) { s.ChannelName = "" },
			func(s *Strip) interface{} { return s.GetChannelName() }},
		{"channel_name_hw",
			func(s *Strip) { s.ChannelName = "Hello, world!" },
			func(s *Strip) interface{} { return s.GetChannelName() }},
	} {
		// Marshal the proto to bytes.
		p1 := NewStrip()
		tt.setFn(p1)
		value := tt.getFn(p1)
		m1, err := p1.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected Marshal() error %s", tt.desc, err)
			continue
		}

		// Read bytes back.
		p2 := NewStrip()
		c, err := p2.Read(m1)
		if err != nil {
			t.Errorf("%s: unexpected Read() error %s", tt.desc, err)
			continue
		}
		if c == 0 {
			t.Errorf("%s: expected count > 0, got %d", tt.desc, c)
			continue
		}
		m2, err := p2.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected Marshal() error %s", tt.desc, err)
			continue
		}

		// Checks.
		if got, want := tt.getFn(p2), value; got != want {
			t.Errorf("%s: got = %v, want %v", tt.desc, got, want)
		}
		if m1len, m2len := len(m1), len(m2); m1len != m2len {
			t.Errorf("%s: len(m1) = %d != len(m2) = %d", tt.desc, m1len, m2len)
		}
	}
}

func TestClen(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		bytes []byte
		clen  int
	}{
		{"empty", []byte{}, 0},
		{"empty_string", []byte{0x0a, 0x00}, 1},
		{"hello", []byte{'h', 'e', 'l', 'l', 'o', 0x00}, 5},
		{"double", []byte{'1', 0x00, '2', 0x00}, 1},
	} {
		if got, want := clen(tt.bytes), tt.clen; got != want {
			t.Errorf("%s: clen() = %d, want %d", tt.desc, got, want)
		}
	}
}
