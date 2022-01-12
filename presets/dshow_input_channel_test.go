package presets

import "testing"

//-----------------------------------------------------------------------------
// AudioStrip

var audioStrip = []struct {
	desc string
	data []byte

	delayIn   bool
	delay     float32
	directOut bool
	phase     bool
}{
	// phase
	{desc: "phase_off", phase: false,
		data: []byte{0x00}},
	{desc: "phase_on", phase: true,
		data: []byte{0x01}},
	// delay in
	{desc: "delay_out", delayIn: false,
		data: append(make([]byte, delayInOffset), []byte{0x00}...)},
	{desc: "delay_in", delayIn: true,
		data: append(make([]byte, delayInOffset), []byte{0x01}...)},
	// delay
	{desc: "delay_0.0", delay: 0.0,
		data: append(make([]byte, delayOffset), []byte{0x00, 0x00}...)},
	{desc: "delay_0.1", delay: 0.1,
		data: append(make([]byte, delayOffset), []byte{0x0a, 0x00}...)},
	{desc: "delay_1.0", delay: 1.0,
		data: append(make([]byte, delayOffset), []byte{0x60, 0x00}...)},
	{desc: "delay_10.0", delay: 10.0,
		data: append(make([]byte, delayOffset), []byte{0xc0, 0x03}...)},
	{desc: "delay_100", delay: 100.0,
		data: append(make([]byte, delayOffset), []byte{0x80, 0x25}...)},
	{desc: "delay_250", delay: 250.0,
		data: append(make([]byte, delayOffset), []byte{0xc0, 0x5d}...)},
	// direct out
	{desc: "direct_out_off", directOut: false,
		data: append(make([]byte, directOutOffset), []byte{0x00}...)},
	{desc: "direct_out_on", directOut: true,
		data: append(make([]byte, directOutOffset), []byte{0x01}...)},
}

func TestDelay(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range audioStrip {
		p.pb.AudioStrip.Bytes = tt.data
		if got, want := p.Delay(), tt.delay; got != want {
			t.Errorf("%s: Delay() = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestDirectOut(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range audioStrip {
		p.pb.AudioStrip.Bytes = tt.data
		if got, want := p.DirectOut(), tt.directOut; got != want {
			t.Errorf("%s: DirectOut() = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestPhase(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range audioStrip {
		p.pb.AudioStrip.Bytes = tt.data
		if got, want := p.Phase(), tt.phase; got != want {
			t.Errorf("%s: Phase() = %v, want %v", tt.desc, got, want)
		}
	}
}

//-----------------------------------------------------------------------------
// InputStrip

var inputStrip = []struct {
	desc string
	data []byte

	gain    float32
	heat    bool
	pad     bool
	phantom bool
}{
	// phantom
	{desc: "phantom_off", phantom: false,
		data: []byte{0x00, 0x00}},
	{desc: "phantom_on", phantom: true,
		data: []byte{0x00, 0x01}},
	// pad
	{desc: "pad_off", pad: false,
		data: []byte{0x00, 0x00, 0x00}},
	{desc: "pad_on", pad: true,
		data: []byte{0x00, 0x00, 0x01}},
	// gain
	{desc: "gain_+10.0_dB", gain: 10.0,
		data: []byte{0x00, 0x00, 0x00, 0x64, 0x00}},
	{desc: "gain_+10.1_dB", gain: 10.1,
		data: []byte{0x00, 0x00, 0x00, 0x65, 0x00}},
	{desc: "gain_+59.9_dB", gain: 59.9,
		data: []byte{0x00, 0x00, 0x00, 0x57, 0x02}},
	{desc: "gain_+60.0_dB", gain: 60.0,
		data: []byte{0x00, 0x00, 0x00, 0x58, 0x02}},
	// heat
	{desc: "heat_off", heat: false,
		data: append(make([]byte, heatOffset), []byte{0x00}...)},
	{desc: "heat_on", heat: true,
		data: append(make([]byte, heatOffset), []byte{0x01}...)},
	// general
	{desc: "clear",
		data: []byte{0x00, 0x00, 0x00, 0x00, 0x00}},
	{desc: "empty",
		data: []byte{}},
}

func TestGain(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range inputStrip {
		p.pb.InputStrip.Bytes = tt.data
		if got, want := p.Gain(), tt.gain; got != want {
			t.Errorf("%s: Gain() = %f, want %f", tt.desc, got, want)
		}
	}
}

func TestPad(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range inputStrip {
		p.pb.InputStrip.Bytes = tt.data
		if got, want := p.Pad(), tt.pad; got != want {
			t.Errorf("%s: Pad() = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestPhantom(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range inputStrip {
		p.pb.InputStrip.Bytes = tt.data
		if got, want := p.Phantom(), tt.phantom; got != want {
			t.Errorf("%s: Phantom() = %v, want %v", tt.desc, got, want)
		}

	}
}

//-----------------------------------------------------------------------------
// Strip

var strip = []struct {
	desc string
	data []byte

	fader float32
	mute  bool
	name  string
}{
	// mute
	{desc: "mute_off", mute: false,
		data: []byte{0x00}},
	{desc: "mute_on", mute: true,
		data: []byte{0x01}},
	// fader
	{desc: "fader_-INF", fader: -144.0,
		data: append(make([]byte, faderOffset), []byte{0x60, 0xfa, 0xff, 0xff}...)},
	{desc: "fader_-131", fader: -131.2,
		data: append(make([]byte, faderOffset), []byte{0xe0, 0xfa, 0xff, 0xff}...)},
	{desc: "fader_-114", fader: -114.2,
		data: append(make([]byte, faderOffset), []byte{0x8a, 0xfb, 0xff, 0xff}...)},
	{desc: "fader_0.0", fader: 0.0,
		data: append(make([]byte, faderOffset), []byte{0x00, 0x00, 0x00, 0x00}...)},
	{desc: "fader_12.0", fader: 12.0,
		data: append(make([]byte, faderOffset), []byte{0x78, 0x00, 0x00, 0x00}...)},
	// name
	{desc: "name_ch_1", name: "Ch 1",
		data: append(make([]byte, nameOffset), []byte{0x43, 0x68, 0x20, 0x31}...)},
	// general
	{desc: "empty",
		data: []byte{}},
}

func TestFader(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range strip {
		p.pb.Strip.Bytes = tt.data
		if got, want := p.Fader(), tt.fader; got != want {
			t.Errorf("%s: Fader() = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestMute(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range strip {
		p.pb.Strip.Bytes = tt.data
		if got, want := p.Mute(), tt.mute; got != want {
			t.Errorf("%s: Mute() = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestName(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range strip {
		p.pb.Strip.Bytes = tt.data
		if got, want := p.Name(), tt.name; got != want {
			t.Errorf("%s: Name() = %v, want %v", tt.desc, got, want)
		}
	}
}

//-----------------------------------------------------------------------------
// MicLine Strips

var micLineStrips = []struct {
	desc string
	data []byte

	hpf bool
	lpf bool
}{
	// hpf
	{desc: "hpf_off", hpf: false,
		data: []byte{0x00}},
	{desc: "hpf_on", hpf: true,
		data: []byte{0x01}},
	// lpf
	{desc: "lpf_off", lpf: false,
		data: append(make([]byte, lpfOffset), []byte{0x00}...)},
	{desc: "lpf_on", lpf: true,
		data: append(make([]byte, lpfOffset), []byte{0x01}...)},
}

func TestHPF(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range micLineStrips {
		p.pb.MicLineStrips.Bytes = tt.data
		if got, want := p.HPF(), tt.hpf; got != want {
			t.Errorf("%s: HPF() = %v, want %v", tt.desc, got, want)
		}
	}
}
