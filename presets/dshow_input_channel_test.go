package presets

import "testing"

//-----------------------------------------------------------------------------
// AudioStrip

var audioStrip = []struct {
	desc string
	data []byte

	delayIn     bool
	delay       float32
	directOutIn bool
	directOut   float32
	pan         int32
	phase       bool
}{
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
	// direct out in
	{desc: "direct_out_out", directOutIn: false,
		data: append(make([]byte, directOutInOffset), []byte{0x00}...)},
	{desc: "direct_out_in", directOutIn: true,
		data: append(make([]byte, directOutInOffset), []byte{0x01}...)},
	// direct out
	{desc: "direct_out_-INF", directOut: -144.0,
		data: append(make([]byte, directOutOffset), []byte{0x60, 0xfa, 0xff, 0xff}...)},
	{desc: "direct_out_0.0", directOut: 0.0,
		data: append(make([]byte, directOutOffset), []byte{0x00, 0x00, 0x00, 0x00}...)},
	{desc: "direct_out_+12.0", directOut: 12.0,
		data: append(make([]byte, directOutOffset), []byte{0x78, 0x00, 0x00, 0x00}...)},
	// pad
	{desc: "pad_-100", pan: -100,
		data: append(make([]byte, panOffset), []byte{0x9c, 0xff, 0xff, 0xff}...)},
	{desc: "pad_0", pan: 0,
		data: append(make([]byte, panOffset), []byte{0x00, 0x00, 0x00, 0x00}...)},
	{desc: "pad_100", pan: 100,
		data: append(make([]byte, panOffset), []byte{0x64, 0x00, 0x00, 0x00}...)},
	// phase
	{desc: "phase_off", phase: false,
		data: []byte{0x00}},
	{desc: "phase_on", phase: true,
		data: []byte{0x01}},
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

func TestDirectOutIn(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range audioStrip {
		p.pb.AudioStrip.Bytes = tt.data
		if got, want := p.DirectOutIn(), tt.directOutIn; got != want {
			t.Errorf("%s: DirectOutIn() = %v, want %v", tt.desc, got, want)
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

func TestPan(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range audioStrip {
		p.pb.AudioStrip.Bytes = tt.data
		if got, want := p.Pan(), tt.pan; got != want {
			t.Errorf("%s: Pan() = %v, want %v", tt.desc, got, want)
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

	drive   int
	eqIn    bool
	gain    float32
	heatIn  bool
	pad     bool
	phantom bool
	tone    int
}{
	// drive
	{desc: "drive_-5", drive: -5,
		data: append(make([]byte, driveOffset), []byte{0xfb, 0xff, 0xff, 0xff}...)},
	{desc: "drive_0", drive: 0,
		data: append(make([]byte, driveOffset), []byte{0x00, 0x00, 0x00, 0x00}...)},
	{desc: "drive_+5", drive: 5,
		data: append(make([]byte, driveOffset), []byte{0x05, 0x00, 0x00, 0x00}...)},
	// eq in
	{desc: "eq_out", heatIn: false,
		data: append(make([]byte, eqInOffset), []byte{0x00}...)},
	{desc: "eq_in", heatIn: true,
		data: append(make([]byte, eqInOffset), []byte{0x01}...)},
	// gain
	{desc: "gain_+10.0_dB", gain: 10.0,
		data: []byte{0x00, 0x00, 0x00, 0x64, 0x00}},
	{desc: "gain_+10.1_dB", gain: 10.1,
		data: []byte{0x00, 0x00, 0x00, 0x65, 0x00}},
	{desc: "gain_+59.9_dB", gain: 59.9,
		data: []byte{0x00, 0x00, 0x00, 0x57, 0x02}},
	{desc: "gain_+60.0_dB", gain: 60.0,
		data: []byte{0x00, 0x00, 0x00, 0x58, 0x02}},
	// heat in
	{desc: "heat_out", heatIn: false,
		data: append(make([]byte, heatInOffset), []byte{0x00}...)},
	{desc: "heat_in", heatIn: true,
		data: append(make([]byte, heatInOffset), []byte{0x01}...)},
	// pad
	{desc: "pad_off", pad: false,
		data: []byte{0x00, 0x00, 0x00}},
	{desc: "pad_on", pad: true,
		data: []byte{0x00, 0x00, 0x01}},
	// phantom
	{desc: "phantom_off", phantom: false,
		data: []byte{0x00, 0x00}},
	{desc: "phantom_on", phantom: true,
		data: []byte{0x00, 0x01}},
	// tone
	{desc: "tone_0", tone: 0,
		data: append(make([]byte, toneOffset), []byte{0x00, 0x00, 0x00, 0x00}...)},
	{desc: "tone_3", tone: 3,
		data: append(make([]byte, toneOffset), []byte{0x03, 0x00, 0x00, 0x00}...)},
	{desc: "tone_6", tone: 6,
		data: append(make([]byte, toneOffset), []byte{0x06, 0x00, 0x00, 0x00}...)},
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

	hpfIn bool
	hpf   int16
	lpfIn bool
	lpf   int16
}{
	// hpf in
	{desc: "hpf_out", hpfIn: false,
		data: []byte{0x00}},
	{desc: "hpf_in", hpfIn: true,
		data: []byte{0x01}},
	// hpf
	{desc: "hpf_20", hpf: 20,
		data: append(make([]byte, hpfOffset), []byte{0x14, 0x00}...)},
	{desc: "hpf_80", hpf: 80,
		data: append(make([]byte, hpfOffset), []byte{0x50, 0x00}...)},
	{desc: "hpf_20k", hpf: 20000,
		data: append(make([]byte, hpfOffset), []byte{0x20, 0x4e}...)},
	// lpf in
	{desc: "lpf_out", lpfIn: false,
		data: append(make([]byte, lpfInOffset), []byte{0x00}...)},
	{desc: "lpf_in", lpfIn: true,
		data: append(make([]byte, lpfInOffset), []byte{0x01}...)},
	// lpf
	{desc: "lpf_20", lpf: 20,
		data: append(make([]byte, lpfOffset), []byte{0x14, 0x00}...)},
	{desc: "lpf_20k", lpf: 20000,
		data: append(make([]byte, lpfOffset), []byte{0x20, 0x4e}...)},
}

func TestHPFIn(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range micLineStrips {
		p.pb.MicLineStrips.Bytes = tt.data
		if got, want := p.HPFIn(), tt.hpfIn; got != want {
			t.Errorf("%s: HPFIn() = %v, want %v", tt.desc, got, want)
		}
	}
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

func TestLPFIn(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range micLineStrips {
		p.pb.MicLineStrips.Bytes = tt.data
		if got, want := p.LPFIn(), tt.lpfIn; got != want {
			t.Errorf("%s: LPFIn() = %v, want %v", tt.desc, got, want)
		}
	}
}

func TestLPF(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range micLineStrips {
		p.pb.MicLineStrips.Bytes = tt.data
		if got, want := p.LPF(), tt.lpf; got != want {
			t.Errorf("%s: LPF() = %v, want %v", tt.desc, got, want)
		}
	}
}
