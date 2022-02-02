package presets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/kward/venue/presets/datatypes"
	pb "github.com/kward/venue/presets/proto"
	log "github.com/sirupsen/logrus"
)

// type offset struct {
// 	key    string
// 	offset int
// }

// var offsets = []offset{}

type Adjuster interface {
	// https://pkg.go.dev/io#Reader
	io.Reader

	// Marshal enables presets to marshal themselves into bytes.
	Marshal() ([]byte, error)

	// Name of the adjuster type.
	Name() string
}

//-----------------------------------------------------------------------------
// DShowInputChannel

var _ Adjuster = new(DShowInputChannel)

type DShowInputChannel struct {
	pb *pb.DShowInputChannel
	h  *Header
	b  *Body
}

func NewDShowInputChannel() *DShowInputChannel {
	return &DShowInputChannel{
		&pb.DShowInputChannel{},
		NewHeader("Digidesign Storage - 1.0"),
		NewBody(),
	}
}

func (p *DShowInputChannel) Read(bs []byte) (int, error) {
	log.Tracef("%s.Read()", p.Name())

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
	log.Tracef("%s.Read()", p.Name())

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

//-----------------------------------------------------------------------------
// Body

var _ Adjuster = new(Body)

type Body struct {
	pb *pb.DShowInputChannel_Body
	is *InputStrip
}

func NewBody() *Body {
	return &Body{
		&pb.DShowInputChannel_Body{},
		NewInputStrip(),
	}
}

func (p *Body) Read(bs []byte) (int, error) {
	log.Tracef("%s.Read()", p.Name())

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

			p.pb.TokenCount = v

		case datatypes.String:
			v, c := readString(bs, i)
			i += c

			switch token {
			default:
				if token != "" {
					return i, fmt.Errorf("%s: unsupported String token %s", p.Name(), token)
				}
				token = v // Store the token for the next loop.
			} // switch token

		case datatypes.Bytes:
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

			switch token {
			case "AudioMasterStrip":
			case "AudioStrip":
			case "AuxBussesOptions":
			case "AuxBussesOptions2":
			case "BusConfigMode":
			case tInputStrip:
				c, err := p.is.Read(b)
				if err != nil {
					return i, fmt.Errorf("%s: failed to Read %s; %s", p.Name(), p.is.Name(), err)
				}
				i += c
			case "MatrixMasterStrip":
			case "MicLineStrips":
			case "Strip":
			case "StripType":
			default:
				return i, fmt.Errorf("%s: unsupported Bytes token %s", p.Name(), token)
			} // switch token

			token = ""

		default:
			return i, fmt.Errorf("%s: unsupported datatype 0x%02x", p.Name(), dt)
		} // switch dt

		if tokensSeen == p.pb.TokenCount {
			return i, nil
		}
	}

	return len(bs), fmt.Errorf("%s: expected %d tokens, found %d", p.Name(), p.pb.TokenCount, tokensSeen)
}

// Marshal the Body into a slice of bytes.
func (p *Body) Marshal() ([]byte, error) {
	bs := []byte{}
	bs = append(bs, datatypes.WriteTokenCount(1)...) // TODO should be 10.

	m, err := p.is.Marshal()
	if err != nil {
		return nil, err
	}
	bs = append(bs, datatypes.WriteTokenBytes(tInputStrip, m)...)
	bs = append(bs, m...)

	return bs, nil
}

func (p *Body) Name() string { return "Body" }

func (p *Body) InputStrip() *InputStrip { return p.is }

//-----------------------------------------------------------------------------
// Input Strip

var _ Adjuster = new(InputStrip)

type kvType int

type kvParam struct {
	offset    int
	readFn    func(*InputStrip, []byte, int) int
	marshalFn func(*InputStrip, []byte, int)
}

type InputStrip struct {
	pb.DShowInputChannel_InputStrip
	params map[string]kvParam
}

const (
	tInputStrip    = "Input Strip"
	inputStripSize = 746
)

func NewInputStrip() *InputStrip {
	return &InputStrip{
		pb.DShowInputChannel_InputStrip{},
		map[string]kvParam{
			"phantom": {1,
				func(is *InputStrip, bs []byte, o int) (c int) {
					is.Phantom, c = readBool(bs, o)
					return
				},
				func(is *InputStrip, bs []byte, o int) { writeBool(bs, o, is.Phantom) }},
			"pad": {2,
				func(is *InputStrip, bs []byte, o int) (c int) {
					is.Pad, c = readBool(bs, o)
					return
				},
				func(is *InputStrip, bs []byte, o int) { writeBool(bs, o, is.Pad) }},
			"gain": {3,
				func(is *InputStrip, bs []byte, o int) (c int) {
					i32, c := readInt32(bs, o)
					is.Gain = float32(i32) / 10
					return
				},
				func(is *InputStrip, bs []byte, o int) { writeInt32(bs, o, int32(is.Gain*10)) }},
			"eq_in": {14,
				func(is *InputStrip, bs []byte, o int) (c int) {
					is.EqIn, c = readBool(bs, o)
					return
				},
				func(is *InputStrip, bs []byte, o int) { writeBool(bs, o, is.EqIn) }},
			"heat_in": {737,
				func(is *InputStrip, bs []byte, o int) (c int) {
					is.HeatIn, c = readBool(bs, o)
					return
				},
				func(is *InputStrip, bs []byte, o int) { writeBool(bs, o, is.HeatIn) }},
			"drive": {738,
				func(is *InputStrip, bs []byte, o int) (c int) {
					is.Drive, c = readInt32(bs, o)
					return
				},
				func(is *InputStrip, bs []byte, o int) { writeInt32(bs, o, is.Drive) }},
			"tone": {742,
				func(is *InputStrip, bs []byte, o int) (c int) {
					is.Tone, c = readInt32(bs, o)
					return
				},
				func(is *InputStrip, bs []byte, o int) { writeInt32(bs, o, is.Tone) }},
		},
	}
}

// Read InputStrip values from a slice of bytes.
func (p *InputStrip) Read(bs []byte) (int, error) {
	log.Tracef("%s.Read()", p.Name())
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
	log.Tracef("%s.Marshal()", p.Name())
	bs := make([]byte, inputStripSize)

	for k, pp := range p.params {
		log.Tracef("%s.marshalFn() at %d", k, pp.offset)
		pp.marshalFn(p, bs, pp.offset)
	}
	return bs, nil
}

func (p *InputStrip) Name() string { return "InputStrip" }

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

// const float32size = 4

// func readFloat32(bs []byte, offset int) (float32, error) {
// 	if len(bs) < offset+float32size {
// 		log.Errorf("len(bs): %d offset: %d size: %d", len(bs), offset, float32size)
// 		return 0.0, fmt.Errorf("readFloat32() out of range; len(bs) = %d, need %d", len(bs), offset+float32size)
// 	}
// 	return float32(int32(binary.LittleEndian.Uint32(bs[offset : offset+float32size]))), nil
// }

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
