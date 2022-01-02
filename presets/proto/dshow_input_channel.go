package proto

import (
	"encoding/binary"
	"fmt"

	"github.com/kward/venue/presets/datatypes"
)

func NewDShowInputChannel() *DShowInputChannel {
	return &DShowInputChannel{
		Version:     &DShowInputChannel_Version{},
		UserComment: &DShowInputChannel_UserComment{},
		Strip:       &DShowInputChannel_Strip{},
	}
}

func (p *DShowInputChannel) Read(bs []byte) {
	token := ""
	for i := 0; i < len(bs)-1; i++ {
		// Handle the data type.
		switch datatypes.DataType(bs[i]) {
		case datatypes.Bytes:
			d, c, err := p.ReadBytes(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			switch token {
			case "Digidesign Storage - 1.0":
				p.Strip.Data = d
			}
			i += c
		case datatypes.Token:
			t, c, err := p.ReadToken(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			token = t
			i += c
		case datatypes.String:
			s, c, err := p.ReadString(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			switch token {
			case "Strip":
				p.Strip.Data = s
			}
			i += c
		case datatypes.UInt32:
			u, c, err := p.ReadUInt32(bs, i+1)
			if err != nil {
				fmt.Println(err)
				break // TODO: need to handle better
			}
			switch token {
			case "Version":
				p.Version.Version = u
			}
			i += c
		default:
			fmt.Printf("%02x ", bs[i])
		}

		// Process data for the current token. This assumes that a token is only
		// ever provided once in the data file.
		switch token {
		case "Digidesign Storage - 1.0":
		case "Version":
		case "File Type":
		case "Digidesign D-Show Input Channel Preset File":
		case "User Comment":
		case "AudioMasterStrip":
		case "AudioStrip":
		case "Aux Busses Options":
		case "Aux Busses Options 2":
		case "Bus Config Mode":
		case "InputStrip":
		case "MatrixMasterStrip":
		case "MicLine Strips":
		case "Strip Type":
		default:
			fmt.Printf("unrecognized token: %s\n", token)
		}
	}

	fmt.Println()
	fmt.Printf("Version: offset: %d data: %q\n", p.Version.Offset, p.Version.Version)
	fmt.Printf("User Comment: offset: %d data: %q\n", p.UserComment.Offset, p.UserComment.Data)
}

// ReadBytes reads byte data.
func (p *DShowInputChannel) ReadBytes(bs []byte, o int) ([]byte, int, error) {
	fmt.Printf("ReadBytes() offset: 0x%04x\n", o)
	c := int(bs[o])
	fmt.Printf("  byte count: %d\n", bs[o])
	return bs[o : o+c], 1 + c, nil
}

// ReadString reads a null prefixed and terminated string.
func (p *DShowInputChannel) ReadString(bs []byte, o int) ([]byte, int, error) {
	fmt.Printf("ReadString() offset: 0x%04x\n", o)
	c := int(binary.LittleEndian.Uint32(bs[o : o+4]))
	fmt.Printf("  string size: %d\n", c)
	o += 4
	return bs[o : o+c+1], 4 + c + 1, nil // TODO: potential bug with int()
}

func (p *DShowInputChannel) ReadToken(bs []byte, o int) (string, int, error) {
	fmt.Printf("ReadToken() offset: 0x%04x\n", o)
	t := string(bs[o : o+clen(bs[o:])])
	fmt.Printf("  token: %s\n", t)
	return t, len(t) + 1, nil
}

func (p *DShowInputChannel) ReadUInt32(bs []byte, o int) (uint32, int, error) {
	fmt.Printf("ReadUInt32() offset: 0x%04x\n", o)
	u := binary.LittleEndian.Uint32(bs[o : o+4])
	fmt.Printf("  uint32: %d\n", u)
	return u, 4, nil
}

func clen(bs []byte) int {
	for i := 0; i < len(bs); i++ {
		if bs[i] == 0 {
			return i
		}
	}
	return len(bs)
}
