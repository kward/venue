package proto

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/kward/venue/presets/datatypes"
)

func NewDShowInputChannel() *DShowInputChannel {
	return &DShowInputChannel{
		// Header
		Header: &Data{},
		// Metadata
		Version:     &Data{},
		FileType:    &Data{},
		UserComment: &Data{},
		// Data
		AudioMasterStrip:  &Data{},
		AudioStrip:        &Data{},
		AuxBussesOptions:  &Data{},
		AuxBussesOptions2: &Data{},
		BusConfigMode:     &Data{},
		InputStrip:        &Data{},
		MatrixMasterStrip: &Data{},
		MicLineStrips:     &Data{},
		Strip:             &Data{},
		StripType:         &Data{},
	}
}

// StripName returns the human readable strip name.
func (p *DShowInputChannel) StripName() string {
	return string(p.Strip.Bytes[7:])
}

func (p *DShowInputChannel) Read(bs []byte) {
	token := ""
	tokenCount := int32(0)
	for i := 0; i < len(bs)-1; i++ {
		fmt.Print(datatypes.DataType(bs[i]), " ")

		// Handle the data type.
		switch datatypes.DataType(bs[i]) {

		case datatypes.Token:
			t, c, err := p.readToken(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			i += c

			switch token { // NOTE: set token = "" if a token was a value.
			case "Digidesign D-Show Input Channel Preset File":
				p.Header.Token = token
				token = ""
			case "File Type":
				p.FileType.Token = token
				p.FileType.Str = t
				token = ""
			case "User Comment":
				p.UserComment.Token = token
				p.UserComment.Str = t
				token = "" // TODO: test me
			default:
				token = t
			}

		case datatypes.TokenCount:
			v, c, err := p.readInt32(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			tokenCount = v
			_ = tokenCount // TODO: do something with the tokenCount
			i += c

		case datatypes.Bytes:
			v, c, err := p.readInt32(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			i += c

			b, c, err := p.readBytes(bs, i, int(v))
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			i += c

			switch token {
			case "AudioMasterStrip":
				p.AudioMasterStrip.Token = token
				p.AudioMasterStrip.Bytes = b
			case "AudioStrip":
				p.AudioStrip.Token = token
				p.AudioStrip.Bytes = b
			case "Aux Busses Options":
				p.AuxBussesOptions.Token = token
				p.AuxBussesOptions.Bytes = b
			case "Aux Busses Options 2":
				p.AuxBussesOptions2.Token = token
				p.AuxBussesOptions2.Bytes = b
			case "Bus Config Mode":
				p.BusConfigMode.Token = token
				p.BusConfigMode.Bytes = b
			case "InputStrip":
				p.InputStrip.Token = token
				p.InputStrip.Bytes = b
			case "MatrixMasterStrip":
				p.MatrixMasterStrip.Token = token
				p.MatrixMasterStrip.Bytes = b
			case "MicLine Strips":
				p.MicLineStrips.Token = token
				p.MicLineStrips.Bytes = b
			case "Strip":
				p.Strip.Token = token
				p.Strip.Bytes = b
			case "Strip Type":
				p.StripType.Token = token
				p.StripType.Bytes = b
			default:
				fmt.Printf("unrecognized token: %s\n", token)
				fmt.Printf("  bytes: %d\n", len(b))
			}

		case datatypes.Int32:
			v, c, err := p.readInt32(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			i += c

			switch token {
			case "Version":
				p.Version.Token = token
				p.Version.Int32 = v
			default:
				fmt.Printf("unrecognized token: %s\n", token)
			}

		default:
			fmt.Printf("unrecognized datatype: %02x\n", bs[i])
		}
	}
}

// readBytes reads byte data.
func (p *DShowInputChannel) readBytes(bs []byte, o, c int) ([]byte, int, error) {
	fmt.Printf("readBytes()\n  offset: 0x%04x\n", o)
	return bs[o : o+c], c, nil
}

func (p *DShowInputChannel) readInt32(bs []byte, o int) (int32, int, error) {
	fmt.Printf("readInt32()\n  offset: 0x%04x\n", o)

	var i int32
	buf := bytes.NewReader(bs[o : o+4])
	if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
		return 0, 0, fmt.Errorf("binary.Read failed: ", err)
	}
	fmt.Printf("  int32: %d\n", i)

	return i, 4, nil
}

// ReadString reads a null prefixed and terminated string.
func (p *DShowInputChannel) ReadString(bs []byte, o int) ([]byte, int, error) {
	fmt.Printf("ReadString()\n  offset: 0x%04x\n", o)
	c := int(binary.LittleEndian.Uint32(bs[o : o+4]))
	fmt.Printf("  string size: %d\n", c)
	o += 4
	return bs[o : o+c+1], 4 + c + 1, nil // TODO: potential bug with int()
}

func (p *DShowInputChannel) readToken(bs []byte, o int) (string, int, error) {
	fmt.Printf("readToken()\n  offset: 0x%04x\n", o)
	t := bs[o : o+clen(bs[o:])]
	fmt.Printf("  token: %q\n", t)
	return string(t), len(t) + 1, nil
}

func clen(bs []byte) int {
	for i := 0; i < len(bs); i++ {
		if bs[i] == 0 {
			return i
		}
	}
	return len(bs)
}
