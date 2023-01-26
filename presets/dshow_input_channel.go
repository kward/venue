package presets

import (
  "bytes"
  "encoding/binary"
  "fmt"
  "io"
  "math"
  "sort"

  "github.com/kward/venue/presets/datatypes"
  pb "github.com/kward/venue/presets/proto"
  log "github.com/sirupsen/logrus"
)

//-----------------------------------------------------------------------------
// Adjuster

type Adjuster interface {
  // https://pkg.go.dev/io#Reader
  io.Reader
  fmt.Stringer

  // Marshal the preset data into bytes.
  Marshal() ([]byte, error)

  // Name of the adjuster type.
  Name() string
}

func readAdjusterParams(a Adjuster, params map[string]kvParam, bs []byte) (int, error) {
  for k, pp := range params {
    log.Tracef("%s.readFn() at %d", k, pp.offset)
    if c := pp.readFn(bs, pp.offset, pp.iface); c == 0 {
      return 0, fmt.Errorf("%s: error reading %s", a.Name(), k)
    }
  }
  return len(bs), nil
}

func marshalAdjusterParams(a Adjuster, params map[string]kvParam, size int) ([]byte, error) {
  bs := make([]byte, size)
  for k, pp := range params {
    if pp.offset > size {
      log.Debugf("%s.marshalAdjusterParams(): pp.offset %d > size %d, ending", a.Name(), pp.offset, size)
      return bs, nil
    }
    log.Tracef("%s.marshalFn() at %d", k, pp.offset)
    pp.marshalFn(bs, pp.offset, pp.iface)
  }
  return bs, nil
}

func stringAdjusterParams(a Adjuster, params map[string]kvParam) (s string) {
  type keyOffset struct {
    key    string
    offset int
  }

  kos := make([]keyOffset, 0, len(params))
  for k, p := range params {
    kos = append(kos, keyOffset{k, p.offset})
  }
  sort.Slice(kos, func(i, j int) bool {
    return kos[i].offset < kos[j].offset
  })

  s += fmt.Sprintf("%s\n", a.Name())
  for _, ko := range kos {
    p := params[ko.key]
    s += fmt.Sprintf(" %s: %s\n", ko.key, p.stringFn(p.iface))
  }
  return
}

//-----------------------------------------------------------------------------
// kvParam

type kvParam struct {
  offset    int
  iface     interface{} // Pointer to the parameter value.
  readFn    func(bs []byte, o int, v interface{}) (c int)
  marshalFn func(bs []byte, o int, v interface{})
  stringFn  func(v interface{}) string
}

func readBoolIface(bs []byte, o int, i interface{}) int {
  const boolSize = 1
  if len(bs) < o+boolSize {
    return 0
  }
  v := i.(*bool)
  *v = bs[o] == 1
  return boolSize
}
func marshalBoolIface(bs []byte, o int, i interface{}) {
  v := i.(*bool)
  writeBool(bs, o, *v)
}
func stringBoolIface(i interface{}) string {
  v := i.(*bool)
  return fmt.Sprintf("%v", *v)
}

func readFloat32Iface(bs []byte, o int, i interface{}) int {
  i32, c := readInt32(bs, o)
  v := i.(*float32)
  *v = float32(i32) / 10
  return c
}
func marshalFloat32Iface(bs []byte, o int, i interface{}) {
  v := i.(*float32)
  writeInt32(bs, o, int32((*v)*10))
}
func stringFloat32Iface(i interface{}) string {
  v := i.(*float32)
  return fmt.Sprintf("%0.1f", *v)
}

func readFloat32Iface100(bs []byte, o int, i interface{}) int {
  i32, c := readInt32(bs, o)
  v := i.(*float32)
  *v = float32(i32) / 100
  return c
}
func marshalFloat32Iface100(bs []byte, o int, i interface{}) {
  v := i.(*float32)
  writeInt32(bs, o, int32((*v)*100))
}
func stringFloat32Iface100(i interface{}) string {
  v := i.(*float32)
  return fmt.Sprintf("%0.2f", *v)
}

func readInt32Iface(bs []byte, o int, i interface{}) int {
  i32, c := readInt32(bs, o)
  v := i.(*int32)
  *v = i32
  return c
}
func marshalInt32Iface(bs []byte, o int, i interface{}) {
  v := i.(*int32)
  writeInt32(bs, o, *v)
}
func stringInt32Iface(i interface{}) string {
  v := i.(*int32)
  return fmt.Sprintf("%d", *v)
}

//-----------------------------------------------------------------------------
// DShowInputChannel

var _ Adjuster = new(DShowInputChannel)

type DShowInputChannel struct {
  h *Header
  b *Body
}

func NewDShowInputChannel() *DShowInputChannel {
  return &DShowInputChannel{
    NewHeader("Digidesign Storage - 1.0"),
    NewBody(),
  }
}

// Read Adjuster values from a slice of bytes.
func (p *DShowInputChannel) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())

  for i := 0; i < len(bs)-1; {
    dt, c := datatypes.ReadDataType(bs, i)
    i += c
    switch dt {

    case datatypes.String:
      s, c := readString(bs, i)
      log.Tracef(" string: %v", s)
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
func (p *DShowInputChannel) String() string {
  return p.h.String() + p.b.String()
}

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

// Read Adjuster values from a slice of bytes.
func (p *Header) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())

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
      v, c := readInt32(bs, i)
      if c == 0 {
        return i, fmt.Errorf("%s: error reading Int32", p.Name())
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
      v, c := readString(bs, i)
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
      v, c := readInt32(bs, i)
      if c == 0 {
        return i, fmt.Errorf("%s: error reading TokenCount", p.Name())
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
func (p *Header) String() (s string) {
  s += fmt.Sprintf("[%s]\n", p.Name())
  s += fmt.Sprintf(" TokenCount: %d\n", p.TokenCount)
  s += fmt.Sprintf(" Version: %d\n", p.Version)
  s += fmt.Sprintf(" File Type: %s\n", p.FileType)
  s += fmt.Sprintf(" User Comment: %s\n", p.UserComment)
  return
}

//-----------------------------------------------------------------------------
// Body

var _ Adjuster = new(Body)

// Body holds the main preset values. Body does not expose the values directly
// as it has none to expose. Instead it provides access functions.
type Body struct {
  pb.DShowInputChannel_Body
  adjusters map[string]Adjuster
}

func NewBody() *Body {
  return &Body{
    pb.DShowInputChannel_Body{},
    map[string]Adjuster{
      // Original.
      "AudioStrip":         NewAudioStrip(),
      "Aux Busses Options": NewAuxBussesOptions(),
      "Bus Config Mode":    NewBusConfigMode(),
      "InputStrip":         NewInputStrip(),
      "MicLine Strips":     NewMicLineStrips(),
      "Strip":              NewStrip(),
      "Strip Type":         NewStripType(),
      // Extended.
      "AudioMasterStrip":     NewAudioMasterStrip(),
      "Aux Busses Options 2": NewAuxBussesOptions2(),
      "MatrixMasterStrip":    NewMatrixMasterStrip(),
    },
  }
}

// Read Adjuster values from a slice of bytes.
func (p *Body) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())

  token := ""            // The most recently seen token.
  tokensSeen := int32(0) // How many tokens have been seen so far.

  for i := 0; i < len(bs)-1; {
    if token != "" {
      tokensSeen++
    }
    log.Tracef(" i: %d token: %q tokensSeen: %d", i, token, tokensSeen)

    dt, c := datatypes.ReadDataType(bs, i)
    i += c
    switch dt { // Case statements sorted based on first appearance.

    case datatypes.TokenCount:
      v, c := readInt32(bs, i)
      if c == 0 {
        return i, fmt.Errorf("%s: error reading TokenCount", p.Name())
      }
      i += c

      p.TokenCount = v

    case datatypes.String: // Token name.
      v, c := readString(bs, i)
      i += c

      switch token {
      default:
        if token != "" {
          return i, fmt.Errorf("%s: unsupported String token %s", p.Name(), token)
        }
        token = v // Store the token for the next loop.
      } // switch token

    case datatypes.Bytes: // Token data.
      // Determine how many bytes to read.
      v, c := readInt32(bs, i)
      if c == 0 {
        return i, fmt.Errorf("%s: error determining Bytes count", p.Name())
      }
      i += c

      // Read the bytes.
      b, c := readBytes(bs, i, int(v))
      if c != int(v) {
        return i, fmt.Errorf("%s: error reading Bytes", p.Name())
      }
      i += c

      // Hand off the reading of the token data.
      a, ok := p.adjusters[token]
      if !ok {
        return i, fmt.Errorf("%s: unsupported Bytes token %s", p.Name(), token)
      }
      c, err := a.Read(b)
      if err != nil {
        return i, fmt.Errorf("%s: failed to Read %s; %s", p.Name(), a.Name(), err)
      }

      token = ""

    default:
      return i, fmt.Errorf("%s: unsupported datatype 0x%02x", p.Name(), dt)
    } // switch dt

    if tokensSeen == p.TokenCount {
      return i, nil
    }
  }

  return len(bs), fmt.Errorf("%s: expected %d tokens, found %d", p.Name(), p.TokenCount, tokensSeen)
}

// Marshal the Body into a slice of bytes.
func (p *Body) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())

  ks := make([]string, 0, len(p.adjusters))
  for k := range p.adjusters {
    ks = append(ks, k)
  }
  sort.Strings(ks)

  bs := datatypes.WriteTokenCount(int32(len(p.adjusters)))
  for _, k := range ks {
    log.Tracef("%s.Marshal()", k)
    m, err := p.adjusters[k].Marshal()
    if err != nil {
      return nil, err
    }
    bs = append(bs, datatypes.WriteTokenBytes(k, m)...)
  }

  return bs, nil
}

func (p *Body) Name() string { return "Body" }
func (p *Body) String() (s string) {
  ks := make([]string, 0, len(p.adjusters))
  for k := range p.adjusters {
    ks = append(ks, k)
  }
  sort.Strings(ks)

  s += fmt.Sprintf("[%s]\n", p.Name())
  s += fmt.Sprintf(" TokenCount: %d\n", p.TokenCount)
  for _, k := range ks {
    s += p.adjusters[k].String()
  }
  return
}

// func (p *Body) AudioMasterStrip() *AudioMasterStrip {
//  return p.adjusters["AudioMasterStrip"].(*AudioMasterStrip)
// }

// func (p *Body) AudioStrip() *AudioStrip {
//  return p.adjusters["AudioStrip"].(*AudioStrip)
// }

func (p *Body) InputStrip() *InputStrip {
  return p.adjusters["InputStrip"].(*InputStrip)
}

//-----------------------------------------------------------------------------
// Body > AudioMasterStrip

var _ Adjuster = new(AudioMasterStrip)

type AudioMasterStrip struct {
  pb.DShowInputChannel_AudioMasterStrip
  params map[string]kvParam
}

func NewAudioMasterStrip() *AudioMasterStrip {
  return &AudioMasterStrip{
    pb.DShowInputChannel_AudioMasterStrip{},
    map[string]kvParam{},
  }
}

// Read Adjuster values from a slice of bytes.
func (p *AudioMasterStrip) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

const audioMasterStripSize = 0x8d

// Marshal the Adjuster into a slice of bytes.
func (p *AudioMasterStrip) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, audioMasterStripSize)
}

func (p *AudioMasterStrip) Name() string   { return "AudioMasterStrip" }
func (p *AudioMasterStrip) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > AudioStrip

var _ Adjuster = new(AudioStrip)

type AudioStrip struct {
  pb.DShowInputChannel_AudioStrip
  params map[string]kvParam
}

func NewAudioStrip() *AudioStrip {
  a := &AudioStrip{}
  a.params = map[string]kvParam{
    "phaseIn": {0x00, &a.PhaseIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "delayIn": {0x01, &a.DelayIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "delay": {0x02, &a.Delay,
      // Unlike most other floats which are multiplied or divided by 10, the
      // delay is adjusted by 96 for some unknown reason.
      func(bs []byte, o int, i interface{}) int {
        v := i.(*float32)
        i32, c := readInt32(bs, o)
        *v = float32(math.Trunc(float64(i32) / 96))
        return c
      },
      func(bs []byte, o int, i interface{}) {
        v := i.(*float32)
        writeInt32(bs, o, int32(*v)*96)
      },
      stringFloat32Iface,
    },
    "directOutIn": {0x07, &a.DirectOutIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "directOut":   {0x0b, &a.DirectOut, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "pan":         {0x11, &a.Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "left_right":  {0x19, &a.LeftRight, readBoolIface, marshalBoolIface, stringBoolIface},
    "center_mono": {0x1a, &a.CenterMono, readBoolIface, marshalBoolIface, stringBoolIface},
  }
  return a
}

// Read Adjuster values from a slice of bytes.
func (p *AudioStrip) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

const audioStripSize = 0x49

// Marshal the Adjuster into a slice of bytes.
func (p *AudioStrip) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, audioStripSize)
}

func (p *AudioStrip) Name() string   { return "AudioStrip" }
func (p *AudioStrip) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > AuxBussesOptions

var _ Adjuster = new(AuxBussesOptions)

type AuxBussesOptions struct {
  pb.DShowInputChannel_AuxBussesOptions
  params map[string]kvParam
}

func NewAuxBussesOptions() *AuxBussesOptions {
  return &AuxBussesOptions{
    pb.DShowInputChannel_AuxBussesOptions{},
    map[string]kvParam{},
  }
}

// Read Adjuster values from a slice of bytes.
func (p *AuxBussesOptions) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

const auxBussesOptionsSize = 0x18

// Marshal the Adjuster into a slice of bytes.
func (p *AuxBussesOptions) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, auxBussesOptionsSize)
}

func (p *AuxBussesOptions) Name() string   { return "AuxBussesOptions" }
func (p *AuxBussesOptions) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > AuxBussesOptions2

var _ Adjuster = new(AuxBussesOptions2)

type AuxBussesOptions2 struct {
  pb.DShowInputChannel_AuxBussesOptions2
  params map[string]kvParam
}

func NewAuxBussesOptions2() *AuxBussesOptions2 {
  return &AuxBussesOptions2{
    pb.DShowInputChannel_AuxBussesOptions2{},
    map[string]kvParam{},
  }
}

// Read Adjuster values from a slice of bytes.
func (p *AuxBussesOptions2) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

const auxBussesOptions2Size = 0x01e8

// Marshal the Adjuster into a slice of bytes.
func (p *AuxBussesOptions2) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, auxBussesOptions2Size)
}

func (p *AuxBussesOptions2) Name() string   { return "AuxBussesOptions2" }
func (p *AuxBussesOptions2) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > BusConfigMode

var _ Adjuster = new(BusConfigMode)

type BusConfigMode struct {
  pb.DShowInputChannel_BusConfigMode
  params map[string]kvParam
}

func NewBusConfigMode() *BusConfigMode {
  return &BusConfigMode{
    pb.DShowInputChannel_BusConfigMode{},
    map[string]kvParam{},
  }
}

// Read Adjuster values from a slice of bytes.
func (p *BusConfigMode) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

const busConfigModeSize = 0x0c

// Marshal the Adjuster into a slice of bytes.
func (p *BusConfigMode) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, busConfigModeSize)
}

func (p *BusConfigMode) Name() string   { return "BusConfigMode" }
func (p *BusConfigMode) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > Input Strip

var _ Adjuster = new(InputStrip)

type InputStrip struct {
  pb.DShowInputChannel_InputStrip
  params map[string]kvParam
}

const (
  EQShelf = false
  EQCurve = true
)

func readEqTypeIface(bs []byte, o int, i interface{}) int {
  v := i.(*pb.DShowInputChannelEqType)
  b, c := readBool(bs, o)
  if b {
    *v = pb.DShowInputChannel_EQ_CURVE
  } else {
    *v = pb.DShowInputChannel_EQ_SHELF
  }
  return c
}
func marshalEqTypeIface(bs []byte, o int, i interface{}) {
  v := i.(*pb.DShowInputChannelEqType)
  writeBool(bs, o, (*v).Number() == 1)
}
func stringEqTypeIface(i interface{}) string {
  v := i.(*pb.DShowInputChannelEqType)
  return fmt.Sprintf("%v", *v)
}

// TODO – a ton of things change when analog_eq is enabled :-(
func NewInputStrip() *InputStrip {
  a := &InputStrip{}
  a.params = map[string]kvParam{
    "patched":      {0x00, &a.Patched, readBoolIface, marshalBoolIface, stringBoolIface},
    "phantom":      {0x01, &a.Phantom, readBoolIface, marshalBoolIface, stringBoolIface},
    "pad":          {0x02, &a.Pad, readBoolIface, marshalBoolIface, stringBoolIface},
    "gain_db":      {0x03, &a.Gain, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "input_direct": {0x0d, &a.InputDirect, readBoolIface, marshalBoolIface, stringBoolIface},

    "eq_in":               {0x0e, &a.EqIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "analog_eq":           {0x0f, &a.AnalogEq, readBoolIface, marshalBoolIface, stringBoolIface},
    "eq_high_in":          {0x10, &a.EqHighIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "eq_high_type":        {0x11, &a.EqHighType, readEqTypeIface, marshalEqTypeIface, stringEqTypeIface},
    "eq_high_gain_db":     {0x12, &a.EqHighGain, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "eq_high_freq_hz":     {0x16, &a.EqHighFreq, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "eq_high_q_bw":        {0x1a, &a.EqHighQBw, readFloat32Iface100, marshalFloat32Iface100, stringFloat32Iface100},
    "eq_high_mid_in":      {0x1e, &a.EqHighMidIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "eq_high_mid_type":    {0x1f, &a.EqHighMidType, readEqTypeIface, marshalEqTypeIface, stringEqTypeIface},
    "eq_high_mid_gain_db": {0x20, &a.EqHighMidGain, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "eq_high_mid_freq_hz": {0x24, &a.EqHighMidFreq, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "eq_high_mid_q_bw":    {0x28, &a.EqHighMidQBw, readFloat32Iface100, marshalFloat32Iface100, stringFloat32Iface100},
    "eq_low_mid_in":       {0x2c, &a.EqLowMidIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "eq_low_mid_type":     {0x2d, &a.EqLowMidType, readEqTypeIface, marshalEqTypeIface, stringEqTypeIface},
    "eq_low_mid_gain_db":  {0x2e, &a.EqLowMidGain, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "eq_low_mid_freq_hz":  {0x32, &a.EqLowMidFreq, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "eq_low_mid_q_bw":     {0x36, &a.EqLowMidQBw, readFloat32Iface100, marshalFloat32Iface100, stringFloat32Iface100},
    "eq_low_in":           {0x3a, &a.EqLowIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "eq_low_type":         {0x3b, &a.EqLowType, readEqTypeIface, marshalEqTypeIface, stringEqTypeIface},
    "eq_low_gain_db":      {0x3c, &a.EqLowGain, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "eq_low_freq_hz":      {0x40, &a.EqLowFreq, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "eq_low_q_bw":         {0x44, &a.EqLowQBw, readFloat32Iface100, marshalFloat32Iface100, stringFloat32Iface100},

    "bus_1": {0x48, &a.Bus1, readBoolIface, marshalBoolIface, stringBoolIface},
    "bus_2": {0x4a, &a.Bus2, readBoolIface, marshalBoolIface, stringBoolIface},
    "bus_3": {0x4c, &a.Bus3, readBoolIface, marshalBoolIface, stringBoolIface},
    "bus_4": {0x4e, &a.Bus4, readBoolIface, marshalBoolIface, stringBoolIface},
    "bus_5": {0x50, &a.Bus5, readBoolIface, marshalBoolIface, stringBoolIface},
    "bus_6": {0x52, &a.Bus6, readBoolIface, marshalBoolIface, stringBoolIface},
    "bus_7": {0x54, &a.Bus7, readBoolIface, marshalBoolIface, stringBoolIface},
    "bus_8": {0x56, &a.Bus8, readBoolIface, marshalBoolIface, stringBoolIface},

    "aux_1_in":       {0x5c, &a.Aux1In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_1_pre":      {0x5d, &a.Aux1Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_1_level_db": {0x5e, &a.Aux1Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    // byte
    "aux_2_in":       {0x63, &a.Aux2In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_2_pre":      {0x64, &a.Aux2Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_2_level_db": {0x65, &a.Aux2Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "aux_1_pan":      {0x69, &a.Aux1Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface},

    "aux_3_in":       {0x6d, &a.Aux3In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_3_pre":      {0x6e, &a.Aux3Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_3_level_db": {0x6f, &a.Aux3Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    // byte
    "aux_4_in":       {0x74, &a.Aux4In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_4_pre":      {0x75, &a.Aux4Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_4_level_db": {0x76, &a.Aux4Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "aux_3_pan":      {0x7a, &a.Aux3Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface},

    "aux_5_in":       {0x7e, &a.Aux5In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_5_pre":      {0x7f, &a.Aux5Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_5_level_db": {0x80, &a.Aux5Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    // byte
    "aux_6_in":       {0x85, &a.Aux6In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_6_pre":      {0x86, &a.Aux6Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_6_level_db": {0x87, &a.Aux6Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "aux_5_pan":      {0x8b, &a.Aux5Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface},

    "aux_7_in":       {0x8f, &a.Aux7In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_7_pre":      {0x90, &a.Aux7Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_7_level_db": {0x91, &a.Aux7Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    // byte
    "aux_8_in":       {0x96, &a.Aux8In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_8_pre":      {0x97, &a.Aux8Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_8_level_db": {0x98, &a.Aux8Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "aux_7_pan":      {0x9c, &a.Aux7Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface},

    "aux_9_in":       {0xa0, &a.Aux9In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_9_pre":      {0xa1, &a.Aux9Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_9_level_db": {0xa2, &a.Aux9Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    // byte
    "aux_10_in":       {0xa7, &a.Aux10In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_10_pre":      {0xa8, &a.Aux10Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_10_level_db": {0xa9, &a.Aux10Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "aux_9_pan":       {0xad, &a.Aux9Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface},

    "aux_11_in":       {0xb1, &a.Aux11In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_11_pre":      {0xb2, &a.Aux11Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_11_level_db": {0xb3, &a.Aux11Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    // byte
    "aux_12_in":       {0xb8, &a.Aux12In, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_12_pre":      {0xb9, &a.Aux12Pre, readBoolIface, marshalBoolIface, stringBoolIface},
    "aux_12_level_db": {0xba, &a.Aux12Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "aux_11_pan":      {0xbe, &a.Aux11Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface},

    //"aux_count": {0x128, &a.AuxCount, readInt32Iface, marshalInt32Iface, stringInt32Iface},

    // "heat_in": {0x2e1, &a.HeatIn, readBoolIface, marshalBoolIface, stringBoolIface},
    // "drive":   {0x2e2, &a.Drive, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    // "tone":    {0x2e6, &a.Tone, readInt32Iface, marshalInt32Iface, stringInt32Iface},
  }
  // Set default values.
  a.EqHighMidType = pb.DShowInputChannel_EQ_CURVE
  a.EqLowMidType = pb.DShowInputChannel_EQ_CURVE
  return a
}

// Read Adjuster values from a slice of bytes.
func (p *InputStrip) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

// const inputStripSize = 0x02ea
const inputStripSize = 0x0128 // Old style preset.

// Marshal the Adjuster into a slice of bytes.
func (p *InputStrip) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, inputStripSize)
}

func (p *InputStrip) Name() string   { return "InputStrip" }
func (p *InputStrip) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > MatrixMasterStrip

var _ Adjuster = new(MatrixMasterStrip)

type MatrixMasterStrip struct {
  pb.DShowInputChannel_MatrixMasterStrip
  params map[string]kvParam
}

func NewMatrixMasterStrip() *MatrixMasterStrip {
  return &MatrixMasterStrip{
    pb.DShowInputChannel_MatrixMasterStrip{},
    map[string]kvParam{},
  }
}

// Read Adjuster values from a slice of bytes.
func (p *MatrixMasterStrip) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

const matrixMasterStripSize = 0x0484

// Marshal the Adjuster into a slice of bytes.
func (p *MatrixMasterStrip) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, matrixMasterStripSize)
}

func (p *MatrixMasterStrip) Name() string   { return "MatrixMasterStrip" }
func (p *MatrixMasterStrip) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > MicLineStrips

var _ Adjuster = new(MicLineStrips)

type MicLineStrips struct {
  pb.DShowInputChannel_MicLineStrips
  size   int
  params map[string]kvParam
}

func readEqDynIface(bs []byte, o int, i interface{}) int {
  v := i.(*pb.DShowInputChannelEqDyn)
  b, c := readBool(bs, o)
  if b {
    *v = pb.DShowInputChannel_EQ_POST_DYN
  } else {
    *v = pb.DShowInputChannel_EQ_PRE_DYN
  }
  return c
}
func marshalEqDynIface(bs []byte, o int, i interface{}) {
  v := i.(*pb.DShowInputChannelEqDyn)
  writeBool(bs, o, (*v).Number() == 1)
}
func stringEqDynIface(i interface{}) string {
  v := i.(*pb.DShowInputChannelEqDyn)
  return fmt.Sprintf("%v", *v)
}

func NewMicLineStrips() *MicLineStrips {
  a := &MicLineStrips{}
  a.params = map[string]kvParam{
    "hpfIn":                 {0x00, &a.HpfIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "hpf":                   {0x01, &a.Hpf, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "eq_dyn":                {0x05, &a.EqDyn, readEqDynIface, marshalEqDynIface, stringEqDynIface},
    "comp_lim_in":           {0x06, &a.CompLimIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "comp_lim_threshold_db": {0x07, &a.CompLimThreshold, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "comp_lim_ratio":        {0x0b, &a.CompLimRatio, readFloat32Iface100, marshalFloat32Iface100, stringFloat32Iface100},
    "comp_lim_attack_us":    {0x0f, &a.CompLimAttack, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "comp_lim_release_ms":   {0x13, &a.CompLimRelease, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "comp_lim_knee":         {0x17, &a.CompLimKnee, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "comp_lim_gain_db":      {0x1a, &a.CompLimGain, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "exp_gate_in":           {0x2f, &a.ExpGateIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "exp_gate_threshold_db": {0x30, &a.ExpGateThreshold, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "exp_gate_attack":       {0x34, &a.ExpGateAttack, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "exp_gate_ratio":        {0x38, &a.ExpGateRatio, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "exp_gate_release_ms":   {0x3c, &a.ExpGateRelease, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "exp_gate_hold_ms":      {0x40, &a.ExpGateHold, readInt32Iface, marshalInt32Iface, stringInt32Iface},
    "exp_gate_range_db":     {0x44, &a.ExpGateRange, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "exp_gate_sidechain_in": {0x52, &a.ExpGateSidechainIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "lpfIn":                 {0x58, &a.LpfIn, readBoolIface, marshalBoolIface, stringBoolIface},
    "lpf":                   {0x59, &a.Lpf, readInt32Iface, marshalInt32Iface, stringInt32Iface},
  }
  return a
}

// Read Adjuster values from a slice of bytes.
func (p *MicLineStrips) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  p.size = len(bs)
  return readAdjusterParams(p, p.params, bs)
}

const micLineStripsSize = 0x65

// Marshal the Adjuster into a slice of bytes.
func (p *MicLineStrips) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  if p.size == 0 {
    p.size = micLineStripsSize
  }
  return marshalAdjusterParams(p, p.params, p.size)
}

func (p *MicLineStrips) Name() string   { return "MicLineStrips" }
func (p *MicLineStrips) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > Strip

var _ Adjuster = new(Strip)

type Strip struct {
  pb.DShowInputChannel_Strip
  params map[string]kvParam
}

func NewStrip() *Strip {
  a := &Strip{}
  a.params = map[string]kvParam{
    "mute":  {0, &a.Mute, readBoolIface, marshalBoolIface, stringBoolIface},
    "fader": {2, &a.Fader, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface},
    "channel_name": {6, &a.ChannelName,
      func(bs []byte, o int, i interface{}) (c int) {
        v := i.(*string)
        *v, c = readString(bs, o)
        return
      },
      func(bs []byte, o int, i interface{}) {
        v := i.(*string)
        writeString(bs, o, *v)
      },
      func(i interface{}) string {
        v := i.(*string)
        return fmt.Sprintf("%s", *v)
      },
    },
  }
  return a
}

// Read Adjuster values from a slice of bytes.
func (p *Strip) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

// Marshal the Adjuster into a slice of bytes. This size of a marshaled strip
// depends on the length of `ChannelName` string.
func (p *Strip) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  cn := p.params["channel_name"]
  v := cn.iface.(*string)
  bs := make([]byte, cn.offset+len(*v)+1) // +1 for \0.
  for k, pp := range p.params {
    log.Tracef("%s.marshalFn() at %d", k, pp.offset)
    pp.marshalFn(bs, pp.offset, pp.iface)
  }
  return bs, nil
}

func (p *Strip) Name() string   { return "Strip" }
func (p *Strip) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Body > StripType

var _ Adjuster = new(StripType)

type StripType struct {
  pb.DShowInputChannel_StripType
  params map[string]kvParam
}

func NewStripType() *StripType {
  return &StripType{
    pb.DShowInputChannel_StripType{},
    map[string]kvParam{},
  }
}

// Read Adjuster values from a slice of bytes.
func (p *StripType) Read(bs []byte) (int, error) {
  log.Debugf("%s.Read()", p.Name())
  return readAdjusterParams(p, p.params, bs)
}

const stripTypeSize = 0x02

// Marshal the Adjuster into a slice of bytes.
func (p *StripType) Marshal() ([]byte, error) {
  log.Debugf("%s.Marshal()", p.Name())
  return marshalAdjusterParams(p, p.params, stripTypeSize)
}

func (p *StripType) Name() string   { return "StripType" }
func (p *StripType) String() string { return stringAdjusterParams(p, p.params) }

//-----------------------------------------------------------------------------
// Base functions

func readBool(bs []byte, offset int) (bool, int) {
  const size = 1
  if len(bs) < offset+size {
    return false, 0
  }
  return bs[offset] == 1, size
}
func writeBool(bs []byte, offset int, v bool) {
  b := []byte{0x00}
  if v {
    b = []byte{0x01}
  }
  copy(bs[offset:], b)
}

func readBytes(bs []byte, offset, size int) ([]byte, int) {
  if len(bs) < offset+size {
    return []byte{}, 0
  }
  return bs[offset : offset+size], size
}

func readInt32(bs []byte, offset int) (int32, int) {
  const size = 4
  var i int32
  buf := bytes.NewReader(bs[offset : offset+size])
  if err := binary.Read(buf, binary.LittleEndian, &i); err != nil {
    return 0, 0
  }
  return i, size
}
func writeInt32(bs []byte, offset int, v int32) {
  buf := new(bytes.Buffer)
  if err := binary.Write(buf, binary.LittleEndian, v); err != nil {
    return
  }
  copy(bs[offset:], buf.Bytes())
}

func clen(bs []byte) int {
  for i := 0; i < len(bs); i++ {
    if bs[i] == 0x00 {
      return i
    }
  }
  return len(bs)
}

func readString(bs []byte, o int) (string, int) {
  s := bs[o : o+clen(bs[o:])]
  return string(s), len(s) + 1
}
func writeString(bs []byte, o int, v string) {
  copy(bs[o:], v)
}

func hexdump(bs []byte) {
  c := ""
  ascii := ""
  hex := ""

  for i, b := range bs {
    if i > 1 && i%16 == 0 {
      fmt.Printf("%08x:%s  %s\n", i-16, hex, ascii)
      hex = ""
      ascii = ""
    }
    if i%2 == 0 {
      hex += " "
    }
    hex += fmt.Sprintf("%02x", b)

    switch {
    case uint8(b) < 20 || uint8(b) > 128:
      c = "."
    default:
      c = fmt.Sprintf("%c", b)
    }
    ascii += c
  }
  fmt.Printf("%08x:%-40s  %-16s\n", len(bs)/16*16, hex, ascii)
}
