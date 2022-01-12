package presets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"

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

//=============================================================================
// Preset functions.

//-----------------------------------------------------------------------------
// AudioStrip

const phaseOffset = 0
const phaseSize = 1

// Phase returns `true` if phase flip is enabled.
func (p *DShowInputChannel) Phase() bool {
	log.Debugf("DShowInputChannel.Phase()")

	bs := p.pb.AudioStrip.Bytes
	if len(bs) < phaseOffset+phaseSize {
		return false
	}
	return bs[phaseOffset] == 1
}

const directOutOffset = 7
const directOutSize = 1

// DirectOut returns `true` if the direct out is enabled.
func (p *DShowInputChannel) DirectOut() bool {
	log.Debugf("DShowInputChannel.DirectOut()")

	bs := p.pb.AudioStrip.Bytes
	if len(bs) < directOutOffset+directOutSize {
		return false
	}
	return bs[directOutOffset] == 1
}

const delayOffset = 3
const delaySize = 2
const delayAdj = 0.96

// Delay returns the amount of delay in ms from 0.0 to 250.0.
func (p *DShowInputChannel) Delay() float32 {
	log.Debugf("DShowInputChannel.Delay()")

	bs := p.pb.AudioStrip.Bytes
	if len(bs) < delayOffset+delaySize {
		return 0.0
	}
	// Divide by delayAdj to compensate for VENUE storing 96% of the UI value.
	v := float64(binary.LittleEndian.Uint16(bs[delayOffset:delayOffset+delaySize])) / delayAdj
	// Divide by 100 to shift the decimal two places to the left.
	return float32(math.Trunc(v)) / 100
}

const delayInOffset = 2
const delayInSize = 1

// DelayIn returns `true` if the delay is enabled.
func (p *DShowInputChannel) DelayIn() bool {
	log.Debugf("DShowInputChannel.DelayIn()")

	bs := p.pb.AudioStrip.Bytes
	if len(bs) < delayInOffset+delayInSize {
		return false
	}
	return bs[delayInOffset] == 1
}

//-----------------------------------------------------------------------------
// InputStrip

const gainOffset = 3
const gainSize = 2

// Gain returns the amount of gain applied from 10.0 to 60.0.
func (p *DShowInputChannel) Gain() float32 {
	log.Debugf("DShowInputChannel.Gain()")

	bs := p.pb.InputStrip.Bytes
	if len(bs) < gainOffset+gainSize {
		return 0.0
	}
	// Divide by 10 to shift the decimal one place to the left.
	return float32(binary.LittleEndian.Uint16(bs[gainOffset:gainOffset+gainSize])) / 10
}

const heatOffset = 737
const heatSize = 1

// Heat returns `true` if heat is enabled.
func (p *DShowInputChannel) Heat() bool {
	log.Debugf("DShowInputChannel.Heat()")

	bs := p.pb.InputStrip.Bytes
	if len(bs) < heatOffset+heatSize {
		return false
	}
	return bs[heatOffset] == 1
}

const padOffset = 2
const padSize = 1

// Pad returns `true` if the pad is enabled.
func (p *DShowInputChannel) Pad() bool {
	log.Debugf("DShowInputChannel.Pad()")

	bs := p.pb.InputStrip.Bytes
	if len(bs) < padOffset+padSize {
		return false
	}
	return bs[padOffset] == 1
}

const phantomOffset = 1
const phantomSize = 1

// Phantom returns the input strip phantom state.
func (p *DShowInputChannel) Phantom() bool {
	log.Debugf("DShowInputChannel.Phantom()")

	bs := p.pb.InputStrip.Bytes
	if len(bs) < phantomOffset+phantomSize {
		return false
	}
	return bs[phantomOffset] == 1
}

//-----------------------------------------------------------------------------
// Strip

const faderOffset = 2
const faderSize = 4

// Fader returns the fader value in dB.
func (p *DShowInputChannel) Fader() float32 {
	log.Debugf("DShowInputChannel.Fader()")

	bs := p.pb.Strip.Bytes
	if len(bs) < faderOffset+faderSize {
		return 0.0
	}
	// Divide by 10 to shift the decimal one place to the left.
	return float32(int32(binary.LittleEndian.Uint32(bs[faderOffset:faderOffset+faderSize]))) / 10
}

const muteOffset = 0
const muteSize = 1

// Mute returns `true` if the mute is enabled.
func (p *DShowInputChannel) Mute() bool {
	log.Debugf("DShowInputChannel.Mute()")

	bs := p.pb.Strip.Bytes
	if len(bs) < muteOffset+muteSize {
		return false
	}
	return bs[muteOffset] == 1
}

const nameOffset = 6

// Name returns the strip name.
func (p *DShowInputChannel) Name() string {
	log.Debugf("DShowInputChannel.Name()")

	bs := p.pb.Strip.Bytes
	if len(bs) < nameOffset {
		return ""
	}
	// The strip name comprises all remaining data bytes.
	return string(bs[nameOffset:])
}

//-----------------------------------------------------------------------------
// MicLine Strips

const hpfOffset = 0
const hpfSize = 1

// HPF returns `true` if the high-pass filter is enabled.
func (p *DShowInputChannel) HPF() bool {
	log.Debugf("DShowInputChannel.HPF()")

	bs := p.pb.MicLineStrips.Bytes
	if len(bs) < hpfOffset+hpfSize {
		return false
	}
	return bs[hpfOffset] == 1
}

const lpfOffset = 88
const lpfSize = 1

// LPF returns `true` if the low-pass filter is enabled.
func (p *DShowInputChannel) LPF() bool {
	log.Debugf("DShowInputChannel.LPF()")

	bs := p.pb.MicLineStrips.Bytes
	if len(bs) < lpfOffset+lpfSize {
		return false
	}
	return bs[lpfOffset] == 1
}

//=============================================================================
// General functions

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
	s := ""
	s += fmt.Sprintf("AudioStrip:\n")
	s += fmt.Sprintf("  Phase: %v\n", p.Phase())
	s += fmt.Sprintf("InputStrip:\n")
	s += fmt.Sprintf("  Phantom: %v\n", p.Phantom())
	s += fmt.Sprintf("  Pad: %v\n", p.Pad())
	s += fmt.Sprintf("  Gain: %0.1f\n", p.Gain())
	s += fmt.Sprintf("  Heat: %v\n", p.Heat())
	s += fmt.Sprintf("Strip\n")
	s += fmt.Sprintf("  Name: %s\n", p.Name())
	s += fmt.Sprintf("  Mute: %v\n", p.Mute())
	s += fmt.Sprintf("MicLine Strips\n")
	s += fmt.Sprintf("  HPF: %v\n", p.HPF())
	s += fmt.Sprintf("  LPF: %v\n", p.LPF())
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
