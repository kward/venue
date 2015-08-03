package venue

import (
	"testing"

	vnc "github.com/kward/go-vnc"
)

func TestIntToKeys(t *testing.T) {
	tests := []struct {
		v    int
		keys []uint32
	}{
		{0, []uint32{vnc.Key0}},
		{-123, []uint32{vnc.KeyMinus, vnc.Key1, vnc.Key2, vnc.Key3}},
		{456, []uint32{vnc.Key4, vnc.Key5, vnc.Key6}},
	}

	for _, tt := range tests {
		if got, want := intToKeys(tt.v), tt.keys; !equalSlices(got, want) {
			t.Errorf("incorrect keys; got = %v, want = %v", got, want)
		}
	}
}
