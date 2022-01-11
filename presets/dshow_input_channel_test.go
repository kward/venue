package presets

import "testing"

func TestInputStrip_Gain(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range []struct {
		desc  string
		bytes []byte
		gain  float32
	}{
		{"+10.0_dB", []byte{0x01, 0x00, 0x00, 0x64, 0x00}, 10.0},
		{"+10.1_dB", []byte{0x01, 0x00, 0x00, 0x65, 0x00}, 10.1},
		{"+11.0_dB", []byte{0x01, 0x00, 0x00, 0x6e, 0x00}, 11.0},
		{"+20.0_dB", []byte{0x01, 0x00, 0x00, 0xc8, 0x00}, 20.0},
		{"+59.9_dB", []byte{0x01, 0x00, 0x00, 0x57, 0x02}, 59.9},
		{"+60.0_dB", []byte{0x01, 0x00, 0x00, 0x58, 0x02}, 60.0},
		{"clear", []byte{0x00, 0x00, 0x00, 0x00, 0x00}, 0.0},
		{"empty", []byte{}, 0.0},
	} {
		p.pb.InputStrip.Bytes = tt.bytes
		if got, want := p.InputStrip_Gain(), tt.gain; got != want {
			t.Errorf("%s: InputStrip_Gain() = %f, want %f", tt.desc, got, want)
		}
	}
}

func TestStrip_Name(t *testing.T) {
	p := NewDShowInputChannel()
	for _, tt := range []struct {
		desc  string
		bytes []byte
		name  string
	}{
		{"ch_1", []byte{0x00, 0x00, 0x60, 0xfa, 0xff, 0xff, 0x43, 0x68, 0x20, 0x31}, "Ch 1"},
	} {
		p.pb.Strip.Bytes = tt.bytes
		if got, want := p.Strip_Name(), tt.name; got != want {
			t.Errorf("%s: Strip_Name() = %q, want %q", tt.desc, got, want)
		}

	}
}
