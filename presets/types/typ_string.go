// Code generated by "stringer -type=Typ"; DO NOT EDIT.

package types

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unknown-0]
	_ = x[NLCString-1]
	_ = x[Bytes-2]
}

const _Typ_name = "UnknownNLCStringBytes"

var _Typ_index = [...]uint8{0, 7, 16, 21}

func (i Typ) String() string {
	if i < 0 || i >= Typ(len(_Typ_index)-1) {
		return "Typ(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Typ_name[_Typ_index[i]:_Typ_index[i+1]]
}
