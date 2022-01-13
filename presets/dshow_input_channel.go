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

const delayInOffset = 2

// DelayIn returns `true` if the delay is enabled.
func (p *DShowInputChannel) DelayIn() bool {
	log.Debugf("DShowInputChannel.DelayIn()")
	return readBool(p.pb.AudioStrip.Bytes, delayInOffset)
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

const directOutInOffset = 7

// DirectOutIn returns `true` if the direct out is enabled.
func (p *DShowInputChannel) DirectOutIn() bool {
	log.Debugf("DShowInputChannel.DirectOutIn()")
	return readBool(p.pb.AudioStrip.Bytes, directOutInOffset)
}

const directOutOffset = 11
const directOutSize = 4

// DirectOut returns the direct out value.
func (p *DShowInputChannel) DirectOut() float32 {
	log.Debugf("DShowInputChannel.DirectOut()")

	bs := p.pb.AudioStrip.Bytes
	if len(bs) < directOutOffset+directOutSize {
		return 0
	}
	// Divide by 10 to shift the decimal one place to the left.
	return float32(int32(binary.LittleEndian.Uint32(bs[directOutOffset:directOutOffset+directOutSize]))) / 10
}

const panOffset = 17
const panSize = 4

// Pan returns the pan value.
func (p *DShowInputChannel) Pan() int32 {
	log.Debugf("DShowInputChannel.Pan()")

	bs := p.pb.AudioStrip.Bytes
	if len(bs) < panOffset+panSize {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(bs[panOffset : panOffset+panSize]))
}

//-----------------------------------------------------------------------------
// InputStrip

const phantomOffset = 1

// Phantom returns the input strip phantom state.
func (p *DShowInputChannel) Phantom() bool {
	log.Debugf("DShowInputChannel.Phantom()")
	return readBool(p.pb.InputStrip.Bytes, phantomOffset)
}

const padOffset = 2

// Pad returns `true` if the pad is enabled.
func (p *DShowInputChannel) Pad() bool {
	log.Debugf("DShowInputChannel.Pad()")
	return readBool(p.pb.InputStrip.Bytes, padOffset)
}

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

const eqInOffset = 14

func (p *DShowInputChannel) EQIn() bool {
	log.Debugf("DShowInputChannel.EQIn()")
	return readBool(p.pb.InputStrip.Bytes, eqInOffset)
}

const heatInOffset = 737

// HeatIn returns `true` if heat is enabled.
func (p *DShowInputChannel) HeatIn() bool {
	log.Debugf("DShowInputChannel.HeatIn()")
	return readBool(p.pb.InputStrip.Bytes, heatInOffset)
}

const driveOffset = 738
const driveSize = 4

// Drive returns the drive value.
func (p *DShowInputChannel) Drive() int16 {
	log.Debugf("DShowInputChannel.Drive()")

	bs := p.pb.InputStrip.Bytes
	if len(bs) < driveOffset+driveSize {
		return 0
	}
	return int16(binary.LittleEndian.Uint32(bs[driveOffset : driveOffset+driveSize]))
}

const toneOffset = 742
const toneSize = 4

// Tone returns the tone value.
func (p *DShowInputChannel) Tone() int32 {
	log.Debugf("DShowInputChannel.Tone()")

	bs := p.pb.InputStrip.Bytes
	if len(bs) < toneOffset+toneSize {
		return 0
	}
	return int32(binary.LittleEndian.Uint32(bs[toneOffset : toneOffset+toneSize]))
}

//-----------------------------------------------------------------------------
// Strip

const muteOffset = 0

// Mute returns `true` if the mute is enabled.
func (p *DShowInputChannel) Mute() bool {
	log.Debugf("DShowInputChannel.Mute()")
	return readBool(p.pb.Strip.Bytes, muteOffset)
}

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

const hpfInOffset = 0

// HPFIn returns `true` if the high-pass filter is enabled.
func (p *DShowInputChannel) HPFIn() bool {
	log.Debugf("DShowInputChannel.HPFIn()")
	return readBool(p.pb.MicLineStrips.Bytes, hpfInOffset)
}

const hpfOffset = 1
const hpfSize = 2

// HPF returns the high-pass filter value.
func (p *DShowInputChannel) HPF() int16 {
	log.Debugf("DShowInputChannel.HPF()")

	bs := p.pb.MicLineStrips.Bytes
	if len(bs) < hpfOffset+hpfSize {
		return 0
	}
	return int16(binary.LittleEndian.Uint16(bs[hpfOffset : hpfOffset+hpfSize]))
}

const lpfInOffset = 88

// LPFIn returns `true` if the low-pass filter is enabled.
func (p *DShowInputChannel) LPFIn() bool {
	log.Debugf("DShowInputChannel.LPFIn()")
	return readBool(p.pb.MicLineStrips.Bytes, lpfInOffset)
}

const lpfOffset = 89
const lpfSize = 2

// LPF returns the low-pass filter value.
func (p *DShowInputChannel) LPF() int16 {
	log.Debugf("DShowInputChannel.LPF()")

	bs := p.pb.MicLineStrips.Bytes
	if len(bs) < lpfOffset+lpfSize {
		return 0
	}
	return int16(binary.LittleEndian.Uint16(bs[lpfOffset : lpfOffset+lpfSize]))
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
			t, c, err := readToken(bs, i+1)
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
			v, c, err := readInt32(bs, i+1)
			if err != nil {
				log.Errorf("%s error: %s", datatypes.TokenCount, err)
				break // TODO: need to handle better
			}
			tokenCount = v
			_ = tokenCount // TODO: do something with the tokenCount
			i += c

		case datatypes.Bytes:
			v, c, err := readInt32(bs, i+1)
			if err != nil {
				log.Errorf("%s error: %s", datatypes.Bytes, err)
				break // TODO: need to handle better
			}
			i += c

			b, c, err := readBytes(bs, i+1, int(v))
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
			v, c, err := readInt32(bs, i+1)
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
	s += fmt.Sprintf("  Delay In: %v\n", p.DelayIn())
	s += fmt.Sprintf("  Delay: %v\n", p.Delay())
	s += fmt.Sprintf("  DirectOut In: %v\n", p.DirectOutIn())
	s += fmt.Sprintf("  DirectOut: %v\n", p.DirectOut())
	s += fmt.Sprintf("  Pan: %v\n", p.Pan())
	s += fmt.Sprintf("InputStrip:\n")
	s += fmt.Sprintf("  Phantom: %v\n", p.Phantom())
	s += fmt.Sprintf("  Pad: %v\n", p.Pad())
	s += fmt.Sprintf("  Gain: %0.1f\n", p.Gain())
	s += fmt.Sprintf("  EQ In: %v\n", p.EQIn())
	s += fmt.Sprintf("  Heat In: %v\n", p.HeatIn())
	s += fmt.Sprintf("  Drive: %v\n", p.Drive())
	s += fmt.Sprintf("  Tone: %v\n", p.Tone())
	s += fmt.Sprintf("Strip\n")
	s += fmt.Sprintf("  Mute: %v\n", p.Mute())
	s += fmt.Sprintf("  Fader: %v\n", p.Fader())
	s += fmt.Sprintf("  Name: %s\n", p.Name())
	s += fmt.Sprintf("MicLine Strips\n")
	s += fmt.Sprintf("  HPF In: %v\n", p.HPFIn())
	s += fmt.Sprintf("  HPF: %v\n", p.HPF())
	s += fmt.Sprintf("  LPF In: %v\n", p.LPFIn())
	s += fmt.Sprintf("  LPF: %v\n", p.LPF())
	return s[:len(s)-1] // Strip trailing \n.
}

//=============================================================================
// Local functions.

func clen(bs []byte) int {
	for i := 0; i < len(bs); i++ {
		if bs[i] == 0x00 {
			return i
		}
	}
	return len(bs)
}

func readBool(bs []byte, offset int) bool {
	if len(bs) < offset+1 {
		return false
	}
	return bs[offset] == 1
}

func readBytes(bs []byte, o, c int) ([]byte, int, error) {
	log.Debugf("readBytes()")
	log.Tracef("  offset: 0x%04x", o)
	log.Tracef("  len: %d", c)
	return bs[o : o+c], c, nil
}

func readInt32(bs []byte, o int) (int32, int, error) {
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

func readToken(bs []byte, o int) (string, int, error) {
	log.Debugf("readToken()")
	log.Tracef("  offset: 0x%04x", o)

	t := bs[o : o+clen(bs[o:])]

	log.Tracef("  token: %q", t)
	return string(t), len(t) + 1, nil
}
