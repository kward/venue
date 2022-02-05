package presets

import (
	"bytes"
	"testing"

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

		{"eq_in_true",
			func(is *InputStrip) { is.EqIn = true },
			func(is *InputStrip) interface{} { return is.GetEqIn() }},
		{"eq_in_false",
			func(is *InputStrip) { is.EqIn = false },
			func(is *InputStrip) interface{} { return is.GetEqIn() }},

		// EQ High
		{"eq_high_in_true",
			func(is *InputStrip) { is.EqHighIn = true },
			func(is *InputStrip) interface{} { return is.GetEqHighIn() }},
		{"eq_high_in_false",
			func(is *InputStrip) { is.EqHighIn = false },
			func(is *InputStrip) interface{} { return is.GetEqHighIn() }},
		{"eq_high_type_curve",
			func(is *InputStrip) { is.EqHighType = EQCurve },
			func(is *InputStrip) interface{} { return is.GetEqHighType() }},
		{"eq_high_type_shelf",
			func(is *InputStrip) { is.EqHighType = EQShelf },
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
		{"eq_high_q_10", // minimum
			func(is *InputStrip) { is.EqHighQ = 10.0 },
			func(is *InputStrip) interface{} { return is.GetEqHighQ() }},
		{"eq_high_q_0.1", // maximum
			func(is *InputStrip) { is.EqHighQ = 0.1 },
			func(is *InputStrip) interface{} { return is.GetEqHighQ() }},

		// EQ High Mid
		{"eq_high_mid_in_true",
			func(is *InputStrip) { is.EqHighMidIn = true },
			func(is *InputStrip) interface{} { return is.GetEqHighMidIn() }},
		{"eq_high_mid_in_false",
			func(is *InputStrip) { is.EqHighMidIn = false },
			func(is *InputStrip) interface{} { return is.GetEqHighMidIn() }},

		// EQ Low Mid
		{"eq_low_mid_in_true",
			func(is *InputStrip) { is.EqLowMidIn = true },
			func(is *InputStrip) interface{} { return is.GetEqLowMidIn() }},
		{"eq_low_mid_in_false",
			func(is *InputStrip) { is.EqLowMidIn = false },
			func(is *InputStrip) interface{} { return is.GetEqLowMidIn() }},

		// EQ Low
		{"eq_low_in_true",
			func(is *InputStrip) { is.EqLowIn = true },
			func(is *InputStrip) interface{} { return is.GetEqLowIn() }},
		{"eq_low_in_false",
			func(is *InputStrip) { is.EqLowIn = false },
			func(is *InputStrip) interface{} { return is.GetEqLowIn() }},

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
