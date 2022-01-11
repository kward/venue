package presets

import "testing"

//-----------------------------------------------------------------------------
// AudioStrip

var audioStrip = []struct {
	desc string
	data []byte

	phase bool
}{
	// phase
	{desc: "phase_off", phase: false,
		data: []byte{0x00}},
	{desc: "phase_on", phase: true,
		data: []byte{0x01}},
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

	mute bool
	name string
}{
	// mute
	{desc: "mute_off", mute: false,
		data: []byte{0x00}},
	{desc: "mute_on", mute: true,
		data: []byte{0x01}},
	// name
	{desc: "ch_1", name: "Ch 1",
		data: append(make([]byte, nameOffset), []byte{0x43, 0x68, 0x20, 0x31}...)},
	// general
	{desc: "empty",
		data: []byte{}},
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
