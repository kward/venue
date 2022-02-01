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
			s, c, err := readString(bs, i)
			if err != nil {
				log.Errorf("%s error: %s", datatypes.String, err)
				break // TODO: need to handle better
			}
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
			v, c, err := readInt32(bs, i)
			if err != nil {
				return i, fmt.Errorf("%s: error reading Int32; %s", p.Name(), err)
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
			v, c, err := readString(bs, i)
			if err != nil {
				return i, fmt.Errorf("%s: error reading String; %s", p.Name(), err)
			}
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
			v, c, err := readInt32(bs, i)
			if err != nil {
				return i, fmt.Errorf("%s: error reading TokenCount; %s", p.Name(), err)
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
			v, c, err := readInt32(bs, i)
			if err != nil {
				return i, fmt.Errorf("%s: error reading TokenCount; %s", p.Name(), err)
			}
			i += c

			p.pb.TokenCount = v

		case datatypes.String:
			v, c, err := readString(bs, i)
			if err != nil {
				return i, fmt.Errorf("%s: error reading String; %s", p.Name(), err)
			}
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
			v, c, err := readInt32(bs, i)
			if err != nil {
				return i, fmt.Errorf("%s: error determining Bytes count; %s", p.Name(), err)
			}
			i += c
			// Read the bytes.
			b, err := readBytes(bs, i, int(v))
			if err != nil {
				return i, fmt.Errorf("%s: error reading Bytes; %s", p.Name(), err)
			}
			i += int(v)

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
	key       string
	readFn    func(*InputStrip, []byte) error
	marshalFn func(*InputStrip) []byte
}

const (
	tInputStrip = "Input Strip"
	pad         = "pad"
	phantom     = "phantom"
)

type InputStrip struct {
	pb.DShowInputChannel_InputStrip
	params []kvParam
}

func NewInputStrip() *InputStrip {
	return &InputStrip{
		pb.DShowInputChannel_InputStrip{},
		[]kvParam{
			{"",
				func(*InputStrip, []byte) error { return nil },
				func(*InputStrip) []byte { return make([]byte, 1) }},
			{phantom,
				func(is *InputStrip, bs []byte) (err error) {
					is.Phantom, err = readBool(bs, 1)
					log.Tracef("phantom = %v", is.GetPhantom())
					return
				},
				func(is *InputStrip) []byte {
					log.Tracef("phantom = %v", is.GetPhantom())
					return writeBool(is.Phantom)
				}},
			{pad,
				func(is *InputStrip, bs []byte) (err error) {
					is.Pad, err = readBool(bs, 2)
					return
				},
				func(is *InputStrip) []byte { return writeBool(is.Pad) }},
		},
	}
}

// Read InputStrip values from a slice of bytes.
func (p *InputStrip) Read(bs []byte) (int, error) {
	log.Tracef("%s.Read()", p.Name())
	for _, pp := range p.params {
		log.Tracef("%s.readFn()", pp.key)
		err := pp.readFn(p, bs)
		if err != nil {
			return 0, fmt.Errorf("%s: error reading %s; %v", p.Name(), pp.key, err)
		}
	}
	return len(bs), nil
}

// Marshal the InputStrip into a slice of bytes.
func (p *InputStrip) Marshal() ([]byte, error) {
	log.Tracef("%s.Marshal()", p.Name())
	bs := []byte{}
	for _, pp := range p.params {
		log.Tracef("%s.marshalFn()", pp.key)
		bs = append(bs, pp.marshalFn(p)...)
	}
	return bs, nil
}

func (p *InputStrip) Name() string { return "InputStrip" }

//-----------------------------------------------------------------------------
// Base functions

func clen(bs []byte) int {
	for i := 0; i < len(bs); i++ {
		if bs[i] == 0x00 {
			return i
		}
	}
	return len(bs)
}

func readBool(bs []byte, offset int) (bool, error) {
	log.Tracef("readBool(%v, %d)", bs, offset)
	const size = 1
	if len(bs) < offset+size {
		return false, fmt.Errorf("readBool() out of range; len(bs) = %d, need %d", len(bs), offset+size)
	}
	return bs[offset] == 1, nil
}

func writeBool(v bool) []byte {
	if v {
		return []byte{0x01}
	}
	return []byte{0x00}
}

func readBytes(bs []byte, offset, size int) ([]byte, error) {
	if len(bs) < offset+size {
		return []byte{}, fmt.Errorf("readBytes() out of range; len(bs) = %d, need %d", len(bs), offset+size)
	}
	return bs[offset : offset+size], nil
}

func readFloat32(bs []byte, offset int) (float32, error) {
	const size = 4
	if len(bs) < offset+size {
		log.Errorf("len(bs): %d offset: %d size: %d", len(bs), offset, size)
		return 0.0, fmt.Errorf("readFloat32() out of range; len(bs) = %d, need %d", len(bs), offset+size)
	}
	return float32(int32(binary.LittleEndian.Uint32(bs[offset : offset+size]))), nil
}

const int32size = 4

func readInt32(bs []byte, offset int) (int32, int, error) {
	log.Tracef("readInt32() offset: %d, bs: %02x", offset, bs)
	var i int32
	buf := bytes.NewReader(bs[offset : offset+int32size])
	log.Tracef("buf: %v", buf)
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		return 0, 0, fmt.Errorf("binary.Read failed: %s", err)
	}
	return i, int32size, nil
}

func readString(bs []byte, o int) (string, int, error) {
	log.Debugf("readString()")
	log.Tracef("  offset: 0x%04x", o)

	t := bs[o : o+clen(bs[o:])]

	log.Tracef("  string: %q", t)
	return string(t), len(t) + 1, nil
}
