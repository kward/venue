package presets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/kward/venue/presets/datatypes"
	pb "github.com/kward/venue/presets/proto"
	log "github.com/sirupsen/logrus"
)

type Adjuster interface {
	// https://pkg.go.dev/io#Reader
	io.Reader
	fmt.Stringer

	// Marshal enables presets to marshal themselves into bytes.
	Marshal() ([]byte, error)

	// Name of the adjuster type.
	Name() string
}

func adjusterString(a Adjuster, params map[string]kvParam) (s string) {
	s += fmt.Sprintf("%s\n", a.Name())
	for k, pp := range params {
		s += fmt.Sprintf(" %s: %v\n", k, pp.getFn(a))
	}
	return
}

type kvParam struct {
	offset    int
	readFn    func(Adjuster, []byte, int) int
	marshalFn func(Adjuster, []byte, int)
	getFn     func(Adjuster) interface{}
}

//-----------------------------------------------------------------------------
// DShowInputChannel

var _ Adjuster = new(DShowInputChannel)

type DShowInputChannel struct {
	h *Header
	b *Body
}

func NewDShowInputChannel() *DShowInputChannel {
	return &DShowInputChannel{
		NewHeader("Digidesign Storage - 1.0"),
		NewBody(),
	}
}

func (p *DShowInputChannel) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())

	for i := 0; i < len(bs)-1; {
		dt, c := datatypes.ReadDataType(bs, i)
		i += c
		switch dt {

		case datatypes.String:
			s, c := readString(bs, i)
			i += c

			switch s {
			case "Digidesign Storage - 1.0":
				c, err := p.h.Read(bs[i:])
				if err != nil {
					return i, fmt.Errorf("%s: failed to Read %s; %s", p.Name(), p.h.Name(), err)
				}
				i += c

				c, err = p.b.Read(bs[i:])
				if err != nil {
					return i, fmt.Errorf("%s: failed to Read %s; %s", p.Name(), p.b.Name(), err)
				}
				i += c

			default:
				return i, fmt.Errorf("%s: unsupported String token %s", p.Name(), s)
			} // switch t

		default:
			return i, fmt.Errorf("%s: unsupported datatype 0x%02x", p.Name(), byte(dt))
		} // switch dt
	}

	return len(bs), nil
}

func (p *DShowInputChannel) Marshal() ([]byte, error) {
	bs := datatypes.WriteString("Digidesign Storage - 1.0")

	mbs, err := p.h.Marshal()
	if err != nil {
		return nil, fmt.Errorf("%s: error marshaling %s; %s", p.Name(), p.h.Name(), err)
	}
	bs = append(bs, mbs...)

	mbs, err = p.b.Marshal()
	if err != nil {
		return nil, fmt.Errorf("%s: error marshaling %s; %s", p.Name(), p.b.Name(), err)
	}
	bs = append(bs, mbs...)

	return bs, nil
}

func (p *DShowInputChannel) Name() string { return "DShowInputChannel" }
func (p *DShowInputChannel) String() string {
	return p.h.String() + p.b.String()
}

func (p *DShowInputChannel) Header() *Header { return p.h }
func (p *DShowInputChannel) Body() *Body     { return p.b }

//-----------------------------------------------------------------------------
// Header

var _ Adjuster = new(Header)

type Header struct {
	pb.DShowInputChannel_Header
}

func NewHeader(token string) *Header {
	return &Header{
		pb.DShowInputChannel_Header{
			Token:       token,
			TokenCount:  3,
			FileType:    "Digidesign D-Show Input Channel Preset",
			Version:     1,
			UserComment: "",
		},
	}
}

const (
	tFileType    = "File Type"
	tUserComment = "User Comment"
	tVersion     = "Version"
)

func (p *Header) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())

	token := ""            // The most recently seen token.
	tokensSeen := int32(0) // How many tokens have been seen so far.

	for i := 0; i < len(bs)-1; {
		if token != "" {
			tokensSeen++
		}
		log.Tracef(" token: %q tokensSeen: %d", token, tokensSeen)

		dt, c := datatypes.ReadDataType(bs, i)
		i += c
		switch dt {

		case datatypes.Int32:
			v, c := readInt32(bs, i)
			if c == 0 {
				return i, fmt.Errorf("%s: error reading Int32", p.Name())
			}
			i += c

			switch token {

			case tVersion:
				p.Version = v
				token = ""
			default:
				return i, fmt.Errorf("%s: unsupported Int32 token %s", p.Name(), token)
			}

		case datatypes.String:
			v, c := readString(bs, i)
			i += c

			switch token {
			case tFileType:
				p.FileType = v
				token = ""
			case tUserComment:
				p.UserComment = v
				token = ""
			default:
				if token != "" {
					return i, fmt.Errorf("%s: unsupported String token %s", p.Name(), token)
				}
				token = v // Store the token for the next loop.
			} // switch token

		case datatypes.TokenCount:
			v, c := readInt32(bs, i)
			if c == 0 {
				return i, fmt.Errorf("%s: error reading TokenCount", p.Name())
			}
			i += c

			p.TokenCount = v

		default:
			return i, fmt.Errorf("%s: unsupported datatype 0x%02x", p.Name(), dt)
		} // switch dt

		if tokensSeen == p.TokenCount {
			return i, nil
		}
	}

	return len(bs), fmt.Errorf("%s: expected %d tokens, found %d", p.Name(), p.TokenCount, tokensSeen)
}

func (p *Header) Marshal() ([]byte, error) {
	bs := []byte{}
	bs = append(bs, datatypes.WriteTokenCount(p.TokenCount)...)
	bs = append(bs, datatypes.WriteTokenInt32(tVersion, p.Version)...)
	bs = append(bs, datatypes.WriteTokenString(tFileType, p.FileType)...)
	bs = append(bs, datatypes.WriteTokenString(tUserComment, p.UserComment)...)
	return bs, nil
}

func (p *Header) Name() string { return "Header" }
func (p *Header) String() (s string) {
	s += fmt.Sprintf("[%s]\n", p.Name())
	s += fmt.Sprintf(" TokenCount: %d\n", p.TokenCount)
	s += fmt.Sprintf(" Version: %d\n", p.Version)
	s += fmt.Sprintf(" File Type: %s\n", p.FileType)
	s += fmt.Sprintf(" User Comment: %s\n", p.UserComment)
	return
}

//-----------------------------------------------------------------------------
// Body

var _ Adjuster = new(Body)

// Body holds the main preset values. Body does not expose the values directly
// as it has none to expose. Instead it provides access functions.
type Body struct {
	pb.DShowInputChannel_Body
	adjusters map[string]Adjuster
}

func NewBody() *Body {
	return &Body{
		pb.DShowInputChannel_Body{},
		map[string]Adjuster{
			"AudioMasterStrip": NewAudioMasterStrip(),
			"AudioStrip":       NewAudioStrip(),
			// "aux_busses_options": {tInputStrip, NewInputStrip()},
			// "aux_busses_options2": {tInputStrip, NewInputStrip()},
			// "bus_config_mode": {tInputStrip, NewInputStrip()},
			"InputStrip": NewInputStrip(),
			// "matrix_master_strip": {tInputStrip, NewInputStrip()},
			// "mic_line_strips": {tInputStrip, NewInputStrip()},
			"Strip": NewStrip(),
			// "strip_type": {tInputStrip, NewInputStrip()},
		},
	}
}

func (p *Body) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())

	token := ""            // The most recently seen token.
	tokensSeen := int32(0) // How many tokens have been seen so far.

	for i := 0; i < len(bs)-1; {
		if token != "" {
			tokensSeen++
		}
		log.Tracef(" token: %q tokensSeen: %d", token, tokensSeen)

		dt, c := datatypes.ReadDataType(bs, i)
		i += c
		switch dt { // Case statements sorted based on first appearance.

		case datatypes.TokenCount:
			v, c := readInt32(bs, i)
			if c == 0 {
				return i, fmt.Errorf("%s: error reading TokenCount", p.Name())
			}
			i += c

			p.TokenCount = v

		case datatypes.String: // Token name.
			v, c := readString(bs, i)
			i += c

			switch token {
			default:
				if token != "" {
					return i, fmt.Errorf("%s: unsupported String token %s", p.Name(), token)
				}
				token = v // Store the token for the next loop.
			} // switch token

		case datatypes.Bytes: // Token data.
			// Determine how many bytes to read.
			v, c := readInt32(bs, i)
			if c == 0 {
				return i, fmt.Errorf("%s: error determining Bytes count", p.Name())
			}
			i += c

			// Read the bytes.
			b, c := readBytes(bs, i, int(v))
			if c != int(v) {
				return i, fmt.Errorf("%s: error reading Bytes", p.Name())
			}
			i += c

			// Hand off the reading of the token data.
			a, ok := p.adjusters[token]
			if !ok {
				return i, fmt.Errorf("%s: unsupported Bytes token %s", p.Name(), token)
			}
			c, err := a.Read(b)
			if err != nil {
				return i, fmt.Errorf("%s: failed to Read %s; %s", p.Name(), a.Name(), err)
			}
			i += c

			token = ""

		default:
			return i, fmt.Errorf("%s: unsupported datatype 0x%02x", p.Name(), dt)
		} // switch dt

		if tokensSeen == p.TokenCount {
			return i, nil
		}
	}

	return len(bs), fmt.Errorf("%s: expected %d tokens, found %d", p.Name(), p.TokenCount, tokensSeen)
}

// Marshal the Body into a slice of bytes.
func (p *Body) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())

	bs := []byte{}
	// TODO should be 10.
	bs = append(bs, datatypes.WriteTokenCount(int32(len(p.adjusters)))...)

	for k, a := range p.adjusters {
		log.Tracef("%s.Marshal()", k)
		m, err := a.Marshal()
		if err != nil {
			return nil, err
		}
		bs = append(bs, datatypes.WriteTokenBytes(k, m)...)
		bs = append(bs, m...)
	}

	return bs, nil
}

func (p *Body) Name() string { return "Body" }
func (p *Body) String() (s string) {
	s += fmt.Sprintf("[%s]\n", p.Name())
	s += fmt.Sprintf(" TokenCount: %d\n", p.TokenCount)
	for _, a := range p.adjusters {
		s += a.String()
	}
	return
}

func (p *Body) AudioMasterStrip() *AudioMasterStrip {
	return p.adjusters["AudioMasterStrip"].(*AudioMasterStrip)
}

func (p *Body) AudioStrip() *AudioStrip {
	return p.adjusters["AudioStrip"].(*AudioStrip)
}

func (p *Body) InputStrip() *InputStrip {
	return p.adjusters["InputStrip"].(*InputStrip)
}

//-----------------------------------------------------------------------------
// Body > AudioMasterStrip

var _ Adjuster = new(AudioMasterStrip)

const audioMasterStripSize = 0x8d

type AudioMasterStrip struct {
	pb.DShowInputChannel_AudioMasterStrip
	params map[string]kvParam
}

func NewAudioMasterStrip() *AudioMasterStrip {
	return &AudioMasterStrip{
		pb.DShowInputChannel_AudioMasterStrip{},
		map[string]kvParam{},
	}
}

// Read AudioMasterStrip values from a slice of bytes.
func (p *AudioMasterStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	for k, pp := range p.params {
		log.Tracef("%s.readFn() at %d", k, pp.offset)
		if c := pp.readFn(p, bs, pp.offset); c == 0 {
			return 0, fmt.Errorf("%s: error reading %s", p.Name(), k)
		}
	}
	return len(bs), nil
}

// Marshal the AudioMasterStrip into a slice of bytes.
func (p *AudioMasterStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	bs := make([]byte, audioMasterStripSize)

	for k, pp := range p.params {
		log.Tracef("%s.marshalFn() at %d", k, pp.offset)
		pp.marshalFn(p, bs, pp.offset)
	}
	return bs, nil
}

func (p *AudioMasterStrip) Name() string   { return "AudioMasterStrip" }
func (p *AudioMasterStrip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > AudioStrip

var _ Adjuster = new(AudioStrip)

const (
	audioStripSize = 0x49
	delayAdj       = 0.96
)

type AudioStrip struct {
	pb.DShowInputChannel_AudioStrip
	params map[string]kvParam
}

func NewAudioStrip() *AudioStrip {
	return &AudioStrip{
		pb.DShowInputChannel_AudioStrip{},
		map[string]kvParam{
			"phaseIn": {0,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*AudioStrip).PhaseIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) {
					writeBool(bs, o, a.(*AudioStrip).PhaseIn)
				},
				func(a Adjuster) interface{} { return a.(*AudioStrip).PhaseIn },
			},
			"delayIn": {1,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*AudioStrip).DelayIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) {
					writeBool(bs, o, a.(*AudioStrip).DelayIn)
				},
				func(a Adjuster) interface{} { return a.(*AudioStrip).DelayIn },
			},
			"delay": {2,
				func(a Adjuster, bs []byte, o int) (c int) {
					var i32 int32
					i32, c = readInt32(bs, o)
					a.(*AudioStrip).Delay = float32(math.Trunc(float64(i32)/delayAdj)) / 100
					return
				},
				func(a Adjuster, bs []byte, o int) {
					writeInt32(bs, o, int32(a.(*AudioStrip).Delay*100*delayAdj))
				},
				func(a Adjuster) interface{} { return a.(*AudioStrip).Delay },
			},
			"directOutIn": {7,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*AudioStrip).DirectOutIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) {
					writeBool(bs, o, a.(*AudioStrip).DirectOutIn)
				},
				func(a Adjuster) interface{} { return a.(*AudioStrip).DirectOutIn },
			},
			"directOut": {11,
				func(a Adjuster, bs []byte, o int) (c int) {
					var i32 int32
					i32, c = readInt32(bs, o)
					a.(*AudioStrip).DirectOut = float32(i32) / 10
					return c
				},
				func(a Adjuster, bs []byte, o int) {
					writeInt32(bs, o, int32(a.(*AudioStrip).DirectOut*10))
				},
				func(a Adjuster) interface{} { return a.(*AudioStrip).DirectOut },
			},
			"pan": {17,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*AudioStrip).Pan, c = readInt32(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) {
					writeInt32(bs, o, a.(*AudioStrip).Pan)
				},
				func(a Adjuster) interface{} { return a.(*AudioStrip).Pan },
			},
		},
	}
}

// Read AudioStrip values from a slice of bytes.
func (p *AudioStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	for k, pp := range p.params {
		log.Tracef("%s.readFn() at %d", k, pp.offset)
		if c := pp.readFn(p, bs, pp.offset); c == 0 {
			return 0, fmt.Errorf("%s: error reading %s", p.Name(), k)
		}
	}
	return len(bs), nil
}

// Marshal the AudioStrip into a slice of bytes.
func (p *AudioStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	bs := make([]byte, audioStripSize)

	for k, pp := range p.params {
		log.Tracef("%s.marshalFn() at %d", k, pp.offset)
		pp.marshalFn(p, bs, pp.offset)
	}
	return bs, nil
}

func (p *AudioStrip) Name() string   { return "AudioStrip" }
func (p *AudioStrip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > Input Strip

var _ Adjuster = new(InputStrip)

const inputStripSize = 0x2ea

type InputStrip struct {
	pb.DShowInputChannel_InputStrip
	params map[string]kvParam
}

func NewInputStrip() *InputStrip {
	return &InputStrip{
		pb.DShowInputChannel_InputStrip{},
		map[string]kvParam{
			"phantom": {1,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).Phantom, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).Phantom) },
				func(a Adjuster) interface{} { return a.(*InputStrip).Phantom },
			},
			"pad": {2,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).Pad, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).Pad) },
				func(a Adjuster) interface{} { return a.(*InputStrip).Pad },
			},
			"gain": {3,
				func(a Adjuster, bs []byte, o int) (c int) {
					i32, c := readInt32(bs, o)
					a.(*InputStrip).Gain = float32(i32) / 10
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, int32(a.(*InputStrip).Gain*10)) },
				func(a Adjuster) interface{} { return a.(*InputStrip).Gain },
			},
			"eq_in": {14,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).EqIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).EqIn) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqIn },
			},
			"heat_in": {737,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).HeatIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).HeatIn) },
				func(a Adjuster) interface{} { return a.(*InputStrip).HeatIn },
			},
			"drive": {738,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).Drive, c = readInt32(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, a.(*InputStrip).Drive) },
				func(a Adjuster) interface{} { return a.(*InputStrip).Drive },
			},
			"tone": {742,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).Tone, c = readInt32(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, a.(*InputStrip).Tone) },
				func(a Adjuster) interface{} { return a.(*InputStrip).Tone },
			},
		},
	}
}

// Read InputStrip values from a slice of bytes.
func (p *InputStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	for k, pp := range p.params {
		log.Tracef("%s.readFn() at %d", k, pp.offset)
		if c := pp.readFn(p, bs, pp.offset); c == 0 {
			return 0, fmt.Errorf("%s: error reading %s", p.Name(), k)
		}
	}
	return len(bs), nil
}

// Marshal the InputStrip into a slice of bytes.
func (p *InputStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	bs := make([]byte, inputStripSize)

	for k, pp := range p.params {
		log.Tracef("%s.marshalFn() at %d", k, pp.offset)
		pp.marshalFn(p, bs, pp.offset)
	}
	return bs, nil
}

func (p *InputStrip) Name() string   { return "InputStrip" }
func (p *InputStrip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > Strip

var _ Adjuster = new(Strip)

// const stripSize = 0x2ea

type Strip struct {
	pb.DShowInputChannel_Strip
	params map[string]kvParam
}

func NewStrip() *Strip {
	return &Strip{
		pb.DShowInputChannel_Strip{},
		map[string]kvParam{
			"mute": {0,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*Strip).Mute, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*Strip).Mute) },
				func(a Adjuster) interface{} { return a.(*Strip).Mute },
			},
			"fader": {2,
				func(a Adjuster, bs []byte, o int) (c int) {
					i32, c := readInt32(bs, o)
					a.(*Strip).Fader = float32(i32) / 10
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, int32(a.(*Strip).Fader*10)) },
				func(a Adjuster) interface{} { return a.(*Strip).Fader },
			},
			"channel_name": {6,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*Strip).ChannelName, c = readString(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeString(bs, o, a.(*Strip).ChannelName) },
				func(a Adjuster) interface{} { return a.(*Strip).ChannelName },
			},
		},
	}
}

// Read values into the Adjuster from a slice of bytes.
func (p *Strip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	for k, pp := range p.params {
		log.Tracef("%s.readFn() at %d", k, pp.offset)
		if c := pp.readFn(p, bs, pp.offset); c == 0 {
			return 0, fmt.Errorf("%s: error reading %s", p.Name(), k)
		}
	}
	return len(bs), nil
}

// Marshal the Adjuster into a slice of bytes.
func (p *Strip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())

	cn := p.params["channel_name"]
	bs := make([]byte, cn.offset+len(cn.getFn(p).(string)))

	for k, pp := range p.params {
		log.Tracef("%s.marshalFn() at %d", k, pp.offset)
		pp.marshalFn(p, bs, pp.offset)
	}
	return bs, nil
}

func (p *Strip) Name() string   { return "Strip" }
func (p *Strip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Base functions

const boolSize = 1

func readBool(bs []byte, offset int) (bool, int) {
	if len(bs) < offset+boolSize {
		return false, 0
	}
	return bs[offset] == 1, boolSize
}
func writeBool(bs []byte, offset int, v bool) {
	b := []byte{0x00}
	if v {
		b = []byte{0x01}
	}
	copy(bs[offset:], b)
}

func readBytes(bs []byte, offset, size int) ([]byte, int) {
	if len(bs) < offset+size {
		return []byte{}, 0
	}
	return bs[offset : offset+size], size
}

const int32size = 4

func readInt32(bs []byte, offset int) (int32, int) {
	var i int32
	buf := bytes.NewReader(bs[offset : offset+int32size])
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		return 0, 0
	}
	return i, int32size
}
func writeInt32(bs []byte, offset int, v int32) {
	buf := new(bytes.Buffer)
	if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
		return
	}
	copy(bs[offset:], buf.Bytes())
}

func clen(bs []byte) int {
	for i := 0; i < len(bs); i++ {
		if bs[i] == 0x00 {
			return i
		}
	}
	return len(bs)
}

func readString(bs []byte, o int) (string, int) {
	s := bs[o : o+clen(bs[o:])]
	return string(s), len(s) + 1
}
func writeString(bs []byte, o int, v string) {
	copy(bs[o:], v)
}
