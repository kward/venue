package main

import "testing"

func TestClen(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		bytes []byte
		clen  int
	}{
		{"empty", []byte{}, 0},
		{"empty_string", []byte{0x0a, 0x00}, 1},
		{"hello", []byte{'h', 'e', 'l', 'l', 'o', 0x00}, 5},
		{"double", []byte{'1', 0x00, '2', 0x00}, 1},
	} {
		if got, want := clen(tt.bytes), tt.clen; got != want {
			t.Errorf("%s: clen() = %d, want %d", tt.desc, got, want)
		}
	}
}
