// Package multistates defines the states of a Multi UI element.
package multistates

import (
	"github.com/golang/glog"
	"github.com/kward/venue/venuelib"
)

type MultiState int

//go:generate stringer -type=MultiState

const (
	Unknown MultiState = iota
	Pressed
	Released
)

// MultiState returns the state change of a Multi-* control.
func State(v interface{}) MultiState {
	if glog.V(3) {
		glog.Info(venuelib.FnName())
	}
	switch v.(type) {
	case int, int32, float32:
		switch v {
		case 1, float32(1.0):
			return Pressed
		case 0, float32(0.0):
			return Released
		default:
			glog.Errorf("unknown multistate value: %v", v)
		}
	default:
		glog.Errorf("unrecognized multistate type, value: %v", v)
	}
	return Unknown
}
