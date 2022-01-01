// Package presets provides support for VENUE presets.
package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kward/venue/presets/types"
)

// DShowInputChannel describes the preset.
var DShowInputChannel = []struct {
	token  string
	typ    types.Typ
	offset int
	len    int
}{
	{token: "DigidesignStorage", typ: types.NLCString},                 // "Digidesign Storage - 1.0"
	{token: "DigidesignStorageBytes", typ: types.Bytes, len: 5},        // 0b03 0000 00
	{token: "Version", typ: types.NLCString},                           // "Version"
	{token: "VersionBytes", typ: types.Bytes, len: 5},                  // 0601 0000 00
	{token: "FileType", typ: types.NLCString},                          // "File Type"
	{token: "DigidesignDShowInputChannelPreset", typ: types.NLCString}, // "Digidesign D-Show Input Channel Preset File"
	{token: "UserComment", typ: types.NLCString},                       // "User Comment"
	{token: "UserCommentValue", typ: types.NLCString},                  // TBD
	{token: "PostUserCommentBytes", typ: types.Bytes, len: 5},          // 0b 0a00 0000
	{token: "AudioMasterStrip", typ: types.NLCString},                  // "Audio MasterStrip"
	{token: "AudioStrip", typ: types.NLCString},                        // "AudioStrip"
}

func main() {
	b, err := ioutil.ReadFile("/Users/kward/Documents/D-Show/User Data/Effect Presets/testdata/D-Show Input Channel/211231.00 Ch 1 Clear Console.ich")
	if err != nil {
		os.Exit(1)
	}

	o := 0 // offset
	for _, p := range DShowInputChannel {
		fmt.Printf("token: %s typ: %s \n", p.token, p.typ)
		if p.offset > 0 {
			o = p.offset
		}

		switch p.typ {
		case types.Bytes:
			bb := b[o : o+p.len]
			fmt.Printf("  offset: 0x%04x value: %s", o, hex.EncodeToString(bb))
			o += len(bb)
			fmt.Printf(" new_offset: 0x%04x\n", o)
		case types.CString:
			t := b[o : o+clen(b[o:])]
			fmt.Printf("  offset: 0x%04x value: %q", o, t)
			o += len(t) + 1 // +1 to include the null.
			fmt.Printf(" new_offset: 0x%04x\n", o)
		case types.NLCString:
			for i := o; i < len(b)-o; i++ {
				if b[i] == 0x0a {
					o = i + 1
					t := b[o : o+clen(b[o:])]
					fmt.Printf("  offset: 0x%04x value: %q", o, t)
					o += len(t) + 1 // +1 to include the null.
					fmt.Printf(" new_offset: 0x%04x\n", o)
					break
				}
			}
		}
	}
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}
