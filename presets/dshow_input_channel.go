package presets

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/kward/venue/presets/datatypes"
	pb "github.com/kward/venue/presets/proto"
	log "github.com/sirupsen/logrus"
)

type DShowInputChannel struct {
	pb *pb.DShowInputChannel
}

func NewDShowInputChannel() *DShowInputChannel {
	return &DShowInputChannel{
		&pb.DShowInputChannel{
			// Header
			Header: &pb.Data{},
			// Metadata
			Version:     &pb.Data{},
			FileType:    &pb.Data{},
			UserComment: &pb.Data{},
			// Data
			AudioMasterStrip:  &pb.Data{},
			AudioStrip:        &pb.Data{},
			AuxBussesOptions:  &pb.Data{},
			AuxBussesOptions2: &pb.Data{},
			BusConfigMode:     &pb.Data{},
			InputStrip:        &pb.Data{},
			MatrixMasterStrip: &pb.Data{},
			MicLineStrips:     &pb.Data{},
			Strip:             &pb.Data{},
			StripType:         &pb.Data{},
		},
	}
}

//-----------------------------------------------------------------------------
// Preset functions.

// InputStrip_Gain returns the channel input gain.
func (p *DShowInputChannel) InputStrip_Gain() float32 {
	log.Debugf("DShowInputChannel.InputStrip_Gain()")

	if len(p.pb.InputStrip.Bytes) < 5 {
		return 0.0
	}
	b := p.pb.InputStrip.Bytes[3 : 3+2]
	log.Tracef("  len(b): %d bytes: 0x%04x", len(b), b)

	return float32(binary.LittleEndian.Uint16(b)) * 0.1
}

// StripName returns the human readable strip name.
func (p *DShowInputChannel) Strip_Name() string {
	return string(p.pb.Strip.Bytes[6:])
}

//-----------------------------------------------------------------------------
// General functions.

func (p *DShowInputChannel) Read(bs []byte) error {
	log.Debugf("Read()")
	token := ""
	tokenCount := int32(0)
	for i := 0; i < len(bs)-1; i++ {
		log.Tracef("Read()...")
		dt := datatypes.DataType(bs[i])
		log.Tracef("  datatype: %s (0x%02x)", dt, byte(dt))

		// Handle the data type.
		switch dt {

		case datatypes.Token:
			t, c, err := p.readToken(bs, i+1)
			if err != nil {
				log.Errorf("%s error: %s", datatypes.Token, err)
				break // TODO: need to handle better
			}
			i += c

			switch token { // NOTE: set token = "" if a token was a value.
			case "Digidesign D-Show Input Channel Preset File":
				p.pb.Header.Token = token
				token = ""
			case "File Type":
				p.pb.FileType.Token = token
				p.pb.FileType.Str = t
				token = ""
			case "User Comment":
				p.pb.UserComment.Token = token
				p.pb.UserComment.Str = t
				token = "" // TODO: test me
			default:
				token = t
			}

		case datatypes.TokenCount:
			v, c, err := p.readInt32(bs, i+1)
			if err != nil {
				log.Errorf("%s error: %s", datatypes.TokenCount, err)
				break // TODO: need to handle better
			}
			tokenCount = v
			_ = tokenCount // TODO: do something with the tokenCount
			i += c

		case datatypes.Bytes:
			v, c, err := p.readInt32(bs, i+1)
			if err != nil {
				log.Errorf("%s error: %s", datatypes.Bytes, err)
				break // TODO: need to handle better
			}
			i += c

			b, c, err := p.readBytes(bs, i+1, int(v))
			if err != nil {
				log.Errorf("%s error: %s", datatypes.Bytes, err)
				break // TODO: need to handle better
			}
			i += c

			switch token {
			case "AudioMasterStrip":
				p.pb.AudioMasterStrip.Token = token
				p.pb.AudioMasterStrip.Bytes = b
			case "AudioStrip":
				p.pb.AudioStrip.Token = token
				p.pb.AudioStrip.Bytes = b
			case "Aux Busses Options":
				p.pb.AuxBussesOptions.Token = token
				p.pb.AuxBussesOptions.Bytes = b
			case "Aux Busses Options 2":
				p.pb.AuxBussesOptions2.Token = token
				p.pb.AuxBussesOptions2.Bytes = b
			case "Bus Config Mode":
				p.pb.BusConfigMode.Token = token
				p.pb.BusConfigMode.Bytes = b
			case "InputStrip":
				p.pb.InputStrip.Token = token
				p.pb.InputStrip.Bytes = b
			case "MatrixMasterStrip":
				p.pb.MatrixMasterStrip.Token = token
				p.pb.MatrixMasterStrip.Bytes = b
			case "MicLine Strips":
				p.pb.MicLineStrips.Token = token
				p.pb.MicLineStrips.Bytes = b
			case "Strip":
				p.pb.Strip.Token = token
				p.pb.Strip.Bytes = b
			case "Strip Type":
				p.pb.StripType.Token = token
				p.pb.StripType.Bytes = b
			default:
				log.Errorf("unrecognized token: %s", token)
				log.Tracef("  bytes: %d", len(b))
			}

		case datatypes.Int32:
			v, c, err := p.readInt32(bs, i+1)
			if err != nil {
				log.Errorf("%s error: %s", datatypes.Int32, err)
				break // TODO: need to handle better
			}
			i += c

			switch token {
			case "Version":
				p.pb.Version.Token = token
				p.pb.Version.Int32 = v
			default:
				log.Errorf("unrecognized token: %s", token)
			}

		default:
			log.Errorf("unrecognized datatype: %02x", dt)
		}
	}

	// TODO: Deal with errors.
	return nil
}

// String implements fmt.Stringer.
func (p *DShowInputChannel) String() string {
	s := fmt.Sprintf("Strip\n")
	s += fmt.Sprintf("  Name: %s\n", p.Strip_Name())
	s += fmt.Sprintf("InputStrip:\n")
	s += fmt.Sprintf("  Gain: %0.1f\n", p.InputStrip_Gain())
	return s[:len(s)-1] // Strip trailing \n.
}

// readBytes reads byte data.
func (p *DShowInputChannel) readBytes(bs []byte, o, c int) ([]byte, int, error) {
	log.Debugf("readBytes()")
	log.Tracef("  offset: 0x%04x", o)
	log.Tracef("  len: %d", c)
	return bs[o : o+c], c, nil
}

func (p *DShowInputChannel) readInt32(bs []byte, o int) (int32, int, error) {
	log.Debugf("readInt32()")
	log.Tracef("  offset: 0x%04x", o)

	var i int32
	buf := bytes.NewReader(bs[o : o+4])
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		return 0, 0, fmt.Errorf("binary.Read failed: %s", err)
	}

	log.Tracef("  int32: %d", i)
	return i, 4, nil
}

func (p *DShowInputChannel) readToken(bs []byte, o int) (string, int, error) {
	log.Debugf("readToken()")
	log.Tracef("  offset: 0x%04x", o)

	t := bs[o : o+clen(bs[o:])]

	log.Tracef("  token: %q", t)
	return string(t), len(t) + 1, nil
}

func clen(bs []byte) int {
	for i := 0; i < len(bs); i++ {
		if bs[i] == 0x00 {
			return i
		}
	}
	return len(bs)
}
