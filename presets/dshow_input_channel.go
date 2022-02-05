package presets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"sort"

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

func adjusterRead(a Adjuster, params map[string]kvParam, bs []byte) (int, error) {
	for k, pp := range params {
		log.Tracef("%s.readFn() at %d", k, pp.offset)
		if c := pp.readFn(a, bs, pp.offset); c == 0 {
			return 0, fmt.Errorf("%s: error reading %s", a.Name(), k)
		}
	}
	return len(bs), nil
}

func adjusterMarshal(a Adjuster, params map[string]kvParam, size int) ([]byte, error) {
	bs := make([]byte, size)
	for k, pp := range params {
		log.Tracef("%s.marshalFn() at %d", k, pp.offset)
		pp.marshalFn(a, bs, pp.offset)
	}
	return bs, nil
}

func adjusterString(a Adjuster, params map[string]kvParam) (s string) {
	type keyOffset struct {
		key    string
		offset int
	}

	kos := make([]keyOffset, 0, len(params))
	for k, p := range params {
		kos = append(kos, keyOffset{k, p.offset})
	}
	sort.Slice(kos, func(i, j int) bool {
		return kos[i].offset < kos[j].offset
	})

	s += fmt.Sprintf("%s\n", a.Name())
	for _, ko := range kos {
		s += fmt.Sprintf(" %s: %v\n", ko.key, params[ko.key].getFn(a))
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

// Read Adjuster values from a slice of bytes.
func (p *DShowInputChannel) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())

	for i := 0; i < len(bs)-1; {
		dt, c := datatypes.ReadDataType(bs, i)
		i += c
		switch dt {

		case datatypes.String:
			s, c := readString(bs, i)
			log.Tracef(" string: %v", s)
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

// Read Adjuster values from a slice of bytes.
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
			"AudioMasterStrip":     NewAudioMasterStrip(),
			"AudioStrip":           NewAudioStrip(),
			"Aux Busses Options":   NewAuxBussesOptions(),
			"Aux Busses Options 2": NewAuxBussesOptions2(),
			"Bus Config Mode":      NewBusConfigMode(),
			"InputStrip":           NewInputStrip(),
			"MatrixMasterStrip":    NewMatrixMasterStrip(),
			"MicLine Strips":       NewMicLineStrips(),
			"Strip":                NewStrip(),
			"Strip Type":           NewStripType(),
		},
	}
}

// Read Adjuster values from a slice of bytes.
func (p *Body) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())

	token := ""            // The most recently seen token.
	tokensSeen := int32(0) // How many tokens have been seen so far.

	for i := 0; i < len(bs)-1; {
		if token != "" {
			tokensSeen++
		}
		log.Tracef(" i: %d token: %q tokensSeen: %d", i, token, tokensSeen)

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

	ks := make([]string, 0, len(p.adjusters))
	for k := range p.adjusters {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	bs := datatypes.WriteTokenCount(int32(len(p.adjusters)))
	for _, k := range ks {
		log.Tracef("%s.Marshal()", k)
		m, err := p.adjusters[k].Marshal()
		if err != nil {
			return nil, err
		}
		bs = append(bs, datatypes.WriteTokenBytes(k, m)...)
	}

	return bs, nil
}

func (p *Body) Name() string { return "Body" }
func (p *Body) String() (s string) {
	ks := make([]string, 0, len(p.adjusters))
	for k := range p.adjusters {
		ks = append(ks, k)
	}
	sort.Strings(ks)

	s += fmt.Sprintf("[%s]\n", p.Name())
	s += fmt.Sprintf(" TokenCount: %d\n", p.TokenCount)
	for _, k := range ks {
		s += p.adjusters[k].String()
	}
	return
}

// func (p *Body) AudioMasterStrip() *AudioMasterStrip {
// 	return p.adjusters["AudioMasterStrip"].(*AudioMasterStrip)
// }

// func (p *Body) AudioStrip() *AudioStrip {
// 	return p.adjusters["AudioStrip"].(*AudioStrip)
// }

func (p *Body) InputStrip() *InputStrip {
	return p.adjusters["InputStrip"].(*InputStrip)
}

//-----------------------------------------------------------------------------
// Body > AudioMasterStrip

var _ Adjuster = new(AudioMasterStrip)

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

// Read Adjuster values from a slice of bytes.
func (p *AudioMasterStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const audioMasterStripSize = 0x8d

// Marshal the Adjuster into a slice of bytes.
func (p *AudioMasterStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, audioMasterStripSize)
}

func (p *AudioMasterStrip) Name() string   { return "AudioMasterStrip" }
func (p *AudioMasterStrip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > AudioStrip

var _ Adjuster = new(AudioStrip)

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
			// Unlike most other floats which are multiplied or divided by 10, the
			// delay is adjusted by 96 for some unknown reason.
			"delay": {2,
				func(a Adjuster, bs []byte, o int) (c int) {
					var i32 int32
					i32, c = readInt32(bs, o)
					a.(*AudioStrip).Delay = float32(math.Trunc(float64(i32) / 96))
					return
				},
				func(a Adjuster, bs []byte, o int) {
					writeInt32(bs, o, int32(a.(*AudioStrip).Delay*96))
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

// Read Adjuster values from a slice of bytes.
func (p *AudioStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const audioStripSize = 0x49

// Marshal the Adjuster into a slice of bytes.
func (p *AudioStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, audioStripSize)
}

func (p *AudioStrip) Name() string   { return "AudioStrip" }
func (p *AudioStrip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > AuxBussesOptions

var _ Adjuster = new(AuxBussesOptions)

type AuxBussesOptions struct {
	pb.DShowInputChannel_AuxBussesOptions
	params map[string]kvParam
}

func NewAuxBussesOptions() *AuxBussesOptions {
	return &AuxBussesOptions{
		pb.DShowInputChannel_AuxBussesOptions{},
		map[string]kvParam{},
	}
}

// Read Adjuster values from a slice of bytes.
func (p *AuxBussesOptions) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const auxBussesOptionsSize = 0x18

// Marshal the Adjuster into a slice of bytes.
func (p *AuxBussesOptions) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, auxBussesOptionsSize)
}

func (p *AuxBussesOptions) Name() string   { return "AuxBussesOptions" }
func (p *AuxBussesOptions) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > AuxBussesOptions2

var _ Adjuster = new(AuxBussesOptions2)

type AuxBussesOptions2 struct {
	pb.DShowInputChannel_AuxBussesOptions2
	params map[string]kvParam
}

func NewAuxBussesOptions2() *AuxBussesOptions2 {
	return &AuxBussesOptions2{
		pb.DShowInputChannel_AuxBussesOptions2{},
		map[string]kvParam{},
	}
}

// Read Adjuster values from a slice of bytes.
func (p *AuxBussesOptions2) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const auxBussesOptions2Size = 0x01e8

// Marshal the Adjuster into a slice of bytes.
func (p *AuxBussesOptions2) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, auxBussesOptions2Size)
}

func (p *AuxBussesOptions2) Name() string   { return "AuxBussesOptions2" }
func (p *AuxBussesOptions2) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > BusConfigMode

var _ Adjuster = new(BusConfigMode)

type BusConfigMode struct {
	pb.DShowInputChannel_BusConfigMode
	params map[string]kvParam
}

func NewBusConfigMode() *BusConfigMode {
	return &BusConfigMode{
		pb.DShowInputChannel_BusConfigMode{},
		map[string]kvParam{},
	}
}

// Read Adjuster values from a slice of bytes.
func (p *BusConfigMode) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const busConfigModeSize = 0x0c

// Marshal the Adjuster into a slice of bytes.
func (p *BusConfigMode) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, busConfigModeSize)
}

func (p *BusConfigMode) Name() string   { return "BusConfigMode" }
func (p *BusConfigMode) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > Input Strip

var _ Adjuster = new(InputStrip)

type InputStrip struct {
	pb.DShowInputChannel_InputStrip
	params map[string]kvParam
}

const (
	EQShelf = false
	EQCurve = true
)

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
			"eq_high_in": {16,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).EqHighIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).EqHighIn) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqHighIn },
			},
			"eq_high_type": {17,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).EqHighType, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).EqHighType) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqHighType },
			},
			"eq_high_gain": {18,
				func(a Adjuster, bs []byte, o int) (c int) {
					i32, c := readInt32(bs, o)
					a.(*InputStrip).EqHighGain = float32(i32) / 10
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, int32(a.(*InputStrip).EqHighGain*10)) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqHighGain },
			},
			"eq_high_freq": {22,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).EqHighFreq, c = readInt32(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, a.(*InputStrip).EqHighFreq) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqHighFreq },
			},
			"eq_high_q": {26,
				func(a Adjuster, bs []byte, o int) (c int) {
					i32, c := readInt32(bs, o)
					a.(*InputStrip).EqHighQ = float32(i32) / 100
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, int32(a.(*InputStrip).EqHighQ*100)) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqHighQ },
			},
			"eq_high_mid_in": {30,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).EqHighMidIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).EqHighMidIn) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqHighMidIn },
			},
			"eq_low_mid_in": {44,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).EqLowMidIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).EqLowMidIn) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqLowMidIn },
			},
			"eq_low_in": {58,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*InputStrip).EqLowIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*InputStrip).EqLowIn) },
				func(a Adjuster) interface{} { return a.(*InputStrip).EqLowIn },
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

// Read Adjuster values from a slice of bytes.
func (p *InputStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const inputStripSize = 0x02ea

// Marshal the Adjuster into a slice of bytes.
func (p *InputStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, inputStripSize)
}

func (p *InputStrip) Name() string   { return "InputStrip" }
func (p *InputStrip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > MatrixMasterStrip

var _ Adjuster = new(MatrixMasterStrip)

type MatrixMasterStrip struct {
	pb.DShowInputChannel_MatrixMasterStrip
	params map[string]kvParam
}

func NewMatrixMasterStrip() *MatrixMasterStrip {
	return &MatrixMasterStrip{
		pb.DShowInputChannel_MatrixMasterStrip{},
		map[string]kvParam{},
	}
}

// Read Adjuster values from a slice of bytes.
func (p *MatrixMasterStrip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const matrixMasterStripSize = 0x0484

// Marshal the Adjuster into a slice of bytes.
func (p *MatrixMasterStrip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, matrixMasterStripSize)
}

func (p *MatrixMasterStrip) Name() string   { return "MatrixMasterStrip" }
func (p *MatrixMasterStrip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > MicLineStrips

var _ Adjuster = new(MicLineStrips)

type MicLineStrips struct {
	pb.DShowInputChannel_MicLineStrips
	params map[string]kvParam
}

func NewMicLineStrips() *MicLineStrips {
	return &MicLineStrips{
		pb.DShowInputChannel_MicLineStrips{},
		map[string]kvParam{
			"hpfIn": {0,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*MicLineStrips).HpfIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*MicLineStrips).HpfIn) },
				func(a Adjuster) interface{} { return a.(*MicLineStrips).HpfIn },
			},
			"hpf": {1,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*MicLineStrips).Hpf, c = readInt32(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, a.(*MicLineStrips).Hpf) },
				func(a Adjuster) interface{} { return a.(*MicLineStrips).Hpf },
			},
			"lpfIn": {89,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*MicLineStrips).LpfIn, c = readBool(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeBool(bs, o, a.(*MicLineStrips).LpfIn) },
				func(a Adjuster) interface{} { return a.(*MicLineStrips).LpfIn },
			},
			"lpf": {90,
				func(a Adjuster, bs []byte, o int) (c int) {
					a.(*MicLineStrips).Lpf, c = readInt32(bs, o)
					return
				},
				func(a Adjuster, bs []byte, o int) { writeInt32(bs, o, a.(*MicLineStrips).Lpf) },
				func(a Adjuster) interface{} { return a.(*MicLineStrips).Lpf },
			},
		},
	}
}

// Read Adjuster values from a slice of bytes.
func (p *MicLineStrips) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const micLineStripsSize = 0x65

// Marshal the Adjuster into a slice of bytes.
func (p *MicLineStrips) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, micLineStripsSize)
}

func (p *MicLineStrips) Name() string   { return "MicLineStrips" }
func (p *MicLineStrips) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > Strip

var _ Adjuster = new(Strip)

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

// Read Adjuster values from a slice of bytes.
func (p *Strip) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

// Marshal the Adjuster into a slice of bytes. This size of a marshaled strip
// depends on the length of `ChannelName` string.
func (p *Strip) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	cn := p.params["channel_name"]
	bs := make([]byte, cn.offset+len(cn.getFn(p).(string))+1) // +1 for \0.
	for k, pp := range p.params {
		log.Tracef("%s.marshalFn() at %d", k, pp.offset)
		pp.marshalFn(p, bs, pp.offset)
	}
	return bs, nil
}

func (p *Strip) Name() string   { return "Strip" }
func (p *Strip) String() string { return adjusterString(p, p.params) }

//-----------------------------------------------------------------------------
// Body > StripType

var _ Adjuster = new(StripType)

type StripType struct {
	pb.DShowInputChannel_StripType
	params map[string]kvParam
}

func NewStripType() *StripType {
	return &StripType{
		pb.DShowInputChannel_StripType{},
		map[string]kvParam{},
	}
}

// Read Adjuster values from a slice of bytes.
func (p *StripType) Read(bs []byte) (int, error) {
	log.Debugf("%s.Read()", p.Name())
	return adjusterRead(p, p.params, bs)
}

const stripTypeSize = 0x02

// Marshal the Adjuster into a slice of bytes.
func (p *StripType) Marshal() ([]byte, error) {
	log.Debugf("%s.Marshal()", p.Name())
	return adjusterMarshal(p, p.params, stripTypeSize)
}

func (p *StripType) Name() string   { return "StripType" }
func (p *StripType) String() string { return adjusterString(p, p.params) }

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

func hexdump(bs []byte) {
	c := ""
	ascii := ""
	hex := ""

	for i, b := range bs {
		if i > 1 && i%16 == 0 {
			fmt.Printf("%08x:%s  %s\n", i-16, hex, ascii)
			hex = ""
			ascii = ""
		}
		if i%2 == 0 {
			hex += " "
		}
		hex += fmt.Sprintf("%02x", b)

		switch {
		case uint8(b) < 20 || uint8(b) > 128:
			c = "."
		default:
			c = fmt.Sprintf("%c", b)
		}
		ascii += c
	}
	fmt.Printf("%08x:%-40s  %-16s\n", len(bs)/16*16, hex, ascii)
}
