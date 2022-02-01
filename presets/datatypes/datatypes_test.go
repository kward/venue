package datatypes

import "bytes"
import "testing"

func TestReadDataType(t *testing.T) {
	for _, tt := range []struct {
		desc   string
		in     []byte
		offset int
		dt     DataType
		count  int
	}{
		{"invalid", []byte{}, 0, Invalid, 1},
		{"string", []byte{0x0a, 0x00}, 0, String, 1},
		{"token_count_at_2", []byte{0x00, 0x0b, 0x01, 0x00, 0x00, 0x00}, 1, TokenCount, 1},
	} {
		dt, count := ReadDataType(tt.in, tt.offset)
		if got, want := dt, tt.dt; got != want {
			t.Errorf("%s: ReadDataType() = %s, want %s", tt.desc, got, want)
		}
		if got, want := count, tt.count; got != want {
			t.Errorf("%s: ReadDataType() count = %d, want %d", tt.desc, got, want)
		}
	}
}

func TestWriteBytes(t *testing.T) {
	for _, tt := range []struct {
		desc string
		in   []byte
		out  []byte
	}{
		{"empty", []byte{},
			[]byte{0x0d, 0x00, 0x00, 0x00, 0x00}},
		{"abc", []byte{0x61, 0x62, 0x63},
			[]byte{0x0d, 0x03, 0x00, 0x00, 0x00, 0x61, 0x62, 0x63}},
	} {
		if got, want := WriteBytes(tt.in), tt.out; !bytes.Equal(got, want) {
			t.Errorf("%s: WriteBytes() = %02x, want %v", tt.desc, got, want)
		}
	}
}

func TestWriteInt32(t *testing.T) {
	for _, tt := range []struct {
		desc string
		in   int32
		out  []byte
	}{
		{"negative_one", -1, []byte{0x06, 0xff, 0xff, 0xff, 0xff}},
		{"zero", 0, []byte{0x06, 0x00, 0x00, 0x00, 0x00}},
		{"one", 1, []byte{0x06, 0x01, 0x00, 0x00, 0x00}},
	} {
		if got, want := WriteInt32(tt.in), tt.out; !bytes.Equal(got, want) {
			t.Errorf("%s: WriteInt32() = %02x, want %v", tt.desc, got, want)
		}
	}
}

func TestWriteString(t *testing.T) {
	for _, tt := range []struct {
		desc string
		in   string
		out  []byte
	}{
		{"empty", "",
			[]byte{0x0a, 0x00}},
		{"hi", "Hi!",
			[]byte{0x0a, 0x48, 0x69, 0x21, 0x00}},
	} {
		if got, want := WriteString(tt.in), tt.out; !bytes.Equal(got, want) {
			t.Errorf("%s: WriteString() = %02x, want %v", tt.desc, got, want)
		}
	}
}

func TestWriteTokenBytes(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		token string
		in    []byte
		out   []byte
	}{
		{"empty", "Empty", []byte{},
			[]byte{0x0a, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x00, 0x0d, 0x00, 0x00, 0x00, 0x00}},
		{"abc", "ABC", []byte{0x61, 0x62, 0x63},
			[]byte{0x0a, 0x41, 0x42, 0x43, 0x00, 0x0d, 0x03, 0x00, 0x00, 0x00, 0x61, 0x62, 0x63}},
	} {
		if got, want := WriteTokenBytes(tt.token, tt.in), tt.out; !bytes.Equal(got, want) {
			t.Errorf("%s: WriteTokenBytes() = %02x, want %v", tt.desc, got, want)
		}
	}
}

func TestWriteTokenCount(t *testing.T) {
	for _, tt := range []struct {
		desc string
		in   int32
		out  []byte
	}{
		{"zero", 0, []byte{0x0b, 0x00, 0x00, 0x00, 0x00}},
		{"one", 1, []byte{0x0b, 0x01, 0x00, 0x00, 0x00}},
	} {
		if got, want := WriteTokenCount(tt.in), tt.out; !bytes.Equal(got, want) {
			t.Errorf("%s: WriteTokenCount() = %02x, want %v", tt.desc, got, want)
		}
	}
}

func TestWriteTokenInt32(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		token string
		in    int32
		out   []byte
	}{
		{"negative_one", "NegOne", -1,
			[]byte{0x0a, 0x4e, 0x65, 0x67, 0x4f, 0x6e, 0x65, 0x00, 0x06, 0xff, 0xff, 0xff, 0xff}},
		{"zero", "Zero", 0,
			[]byte{0x0a, 0x5a, 0x65, 0x72, 0x6f, 0x00, 0x06, 0x00, 0x00, 0x00, 0x00}},
		{"one", "One", 1,
			[]byte{0x0a, 0x4f, 0x6e, 0x65, 0x00, 0x06, 0x01, 0x00, 0x00, 0x00}},
	} {
		if got, want := WriteTokenInt32(tt.token, tt.in), tt.out; !bytes.Equal(got, want) {
			t.Errorf("%s: WriteTokenInt32() = %02x, want %v", tt.desc, got, want)
		}
	}
}

func TestWriteTokenString(t *testing.T) {
	for _, tt := range []struct {
		desc  string
		token string
		in    string
		out   []byte
	}{
		{"empty", "Empty", "",
			[]byte{0x0a, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x00, 0x0a, 0x00}},
		{"hi", "Hey", "Hi!",
			[]byte{0x0a, 0x48, 0x65, 0x79, 0x00, 0x0a, 0x48, 0x69, 0x21, 0x00}},
	} {
		if got, want := WriteTokenString(tt.token, tt.in), tt.out; !bytes.Equal(got, want) {
			t.Errorf("%s: WriteTokenString() = %02x, want %v", tt.desc, got, want)
		}
	}
}
