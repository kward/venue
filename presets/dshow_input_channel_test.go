package presets

import (
	"testing"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetLevel(log.TraceLevel)
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
		p := NewDShowInputChannel()
		tt.setFn(p)
		value := tt.getFn(p)
		m, err := p.Marshal()
		if err != nil {
			t.Errorf("%s: unexpected error %s", tt.desc, err)
			continue
		}

		// Read bytes back.
		p = NewDShowInputChannel()
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
