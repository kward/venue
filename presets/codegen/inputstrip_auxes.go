package main

import (
  "fmt"
  "os"
  "text/template"
)

type Entry struct {
  AuxOdd    int
  AuxEven   int
  AuxOffset int
  ProtoNum  int
  GenTest   bool
}
type Data struct {
  Entries []Entry
}

var str = `
// --- proto ---
{{with .Entries}}
{{range .}}bool aux{{.AuxOdd}}In = {{.ProtoNum}};
bool aux{{.AuxOdd}}Pre = {{add .ProtoNum 1}};
float aux{{.AuxOdd}}Level = {{add .ProtoNum 2}};
bool aux{{.AuxEven}}In = {{add .ProtoNum 3}};
// {{add .ProtoNum 4}};
bool aux{{.AuxEven}}Pre = {{add .ProtoNum 5}};
float aux{{.AuxEven}}Level = {{add .ProtoNum 6}};
int32 aux{{.AuxOdd}}Pan = {{add .ProtoNum 7}};

{{end}}{{end}}

// --- code ---
{{with .Entries}}
{{range .}}"aux_{{.AuxOdd}}_in": { {{add .AuxOffset 0 |printf "0x%02x"}}, &a.Aux{{.AuxOdd}}In, readBoolIface, marshalBoolIface, stringBoolIface },
"aux_{{.AuxOdd}}_pre": { {{add .AuxOffset 1 |printf "0x%02x"}}, &a.Aux{{.AuxOdd}}Pre, readBoolIface, marshalBoolIface, stringBoolIface },
"aux_{{.AuxOdd}}_level_db": { {{add .AuxOffset 2 |printf "0x%02x"}}, &a.Aux{{.AuxOdd}}Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface },
// byte
"aux_{{.AuxEven}}_in": { {{add .AuxOffset 7 |printf "0x%02x"}}, &a.Aux{{.AuxEven}}In, readBoolIface, marshalBoolIface, stringBoolIface },
"aux_{{.AuxEven}}_pre": { {{add .AuxOffset 8 |printf "0x%02x"}}, &a.Aux{{.AuxEven}}Pre, readBoolIface, marshalBoolIface, stringBoolIface },
"aux_{{.AuxEven}}_level_db": { {{add .AuxOffset 9 |printf "0x%02x"}}, &a.Aux{{.AuxEven}}Level, readFloat32Iface, marshalFloat32Iface, stringFloat32Iface },
"aux_{{.AuxOdd}}_pan": { {{add .AuxOffset 13 |printf "0x%02x"}}, &a.Aux{{.AuxOdd}}Pan, readInt32Iface, marshalInt32Iface, stringInt32Iface },

{{end}}{{end}}

// --- test ---
{{with .Entries}}
{{range .}}{{if .GenTest}}
// Aux {{.AuxOdd}} and {{.AuxEven}}
{"aux_{{.AuxOdd}}_in_true",
  func(is *InputStrip) { is.Aux{{.AuxOdd}}In = true },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}In() }},
{"aux_{{.AuxOdd}}_in_false",
  func(is *InputStrip) { is.Aux{{.AuxOdd}}In = false },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}In() }},
{"aux_{{.AuxOdd}}_pre_true",
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Pre = true },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Pre() }},
{"aux_{{.AuxOdd}}_pre_false",
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Pre = false },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Pre() }},
{"aux_{{.AuxOdd}}_level_-144_dB", // minimum
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Level = -144.0 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Level() }},
{"aux_{{.AuxOdd}}_level_0.0_dB",
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Level = 0.0 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Level() }},
{"aux_{{.AuxOdd}}_level_+12.0_dB", // maximum
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Level = 12.0 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Level() }},

{"aux_{{.AuxEven}}_in_true",
  func(is *InputStrip) { is.Aux{{.AuxEven}}In = true },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxEven}}In() }},
{"aux_{{.AuxEven}}_in_false",
  func(is *InputStrip) { is.Aux{{.AuxEven}}In = false },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxEven}}In() }},
{"aux_{{.AuxEven}}_pre_true",
  func(is *InputStrip) { is.Aux{{.AuxEven}}Pre = true },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxEven}}Pre() }},
{"aux_{{.AuxEven}}_pre_false",
  func(is *InputStrip) { is.Aux{{.AuxEven}}Pre = false },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxEven}}Pre() }},
{"aux_{{.AuxEven}}_level_-144_dB", // minimum
  func(is *InputStrip) { is.Aux{{.AuxEven}}Level = -144.0 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxEven}}Level() }},
{"aux_{{.AuxEven}}_level_0.0_dB",
  func(is *InputStrip) { is.Aux{{.AuxEven}}Level = 0.0 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxEven}}Level() }},
{"aux_{{.AuxEven}}_level_+12.0_dB", // maximum
  func(is *InputStrip) { is.Aux{{.AuxEven}}Level = 12.0 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxEven}}Level() }},

{"aux_{{.AuxOdd}}_pan_-100", // left
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Pan = -100 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Pan() }},
{"aux_{{.AuxOdd}}_pan_0", // center
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Pan = 0 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Pan() }},
{"aux_{{.AuxOdd}}_pan_100", // right
  func(is *InputStrip) { is.Aux{{.AuxOdd}}Pan = 100 },
  func(is *InputStrip) interface{} { return is.GetAux{{.AuxOdd}}Pan() }},
{{end}}{{end}}{{end}}
`

func main() {
  funcMap := template.FuncMap{
    "add": func(x int, y int) int { return x + y },
  }
  tmpl, err := template.New("gen").Funcs(funcMap).Parse(str)
  if err != nil {
    fmt.Printf("error; %s", err)
    return
  }

  es := []Entry{}
  auxOffset := int(0x5c)
  protoNum := 40
  for aux := 1; aux <= 12; aux += 2 {
    e := Entry{
      AuxOdd:    aux,
      AuxEven:   aux + 1,
      AuxOffset: auxOffset,
      ProtoNum:  protoNum,
    }
    switch aux {
    case 1, 11, 23:
      e.GenTest = true
    }
    es = append(es, e)
    auxOffset += 17
    protoNum += 8
  }
  // if err := tmpl.Execute(os.Stdout, Data{Entries: es}); err != nil {
  //   fmt.Printf("error; %s", err)
  //   return
  // }

  // // es = []Entry{}
  // auxOffset = int(0x12d)
  // // protoNum = 40
  // for aux := 13; aux <= 24; aux++ {
  //   e := Entry{
  //     Aux:            aux,
  //     AuxOffsetBase:  fmt.Sprintf("0x%02x", auxOffset+0),
  //     AuxOffsetPre:   fmt.Sprintf("0x%02x", auxOffset+1),
  //     AuxOffsetLevel: fmt.Sprintf("0x%02x", auxOffset+2),
  //     AuxOffsetPan:   fmt.Sprintf("0x%02x", auxOffset+6),
  //     ProtoNumBase:   protoNum + 0,
  //     ProtoNumPre:    protoNum + 1,
  //     ProtoNumLevel:  protoNum + 2,
  //     ProtoNumPan:    protoNum + 3,
  //   }
  //   switch aux {
  //   case 1, 12, 24:
  //     e.GenTest = true
  //   }
  //   es = append(es, e)
  //   auxOffset += 9
  //   protoNum += 4
  // }

  if err := tmpl.Execute(os.Stdout, Data{Entries: es}); err != nil {
    fmt.Printf("error; %s", err)
    return
  }
}
