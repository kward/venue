package venue

import "time"

type options struct {
	inputs  uint
	refresh time.Duration // VNC Framebuffer refresh period
}

// Inputs is an option for New() that sets the number of inputs.
func Inputs(v uint) func(*options) error {
	return func(o *options) error { return o.setInputs(v) }
}

// setInputs sets the number of inputs.
func (o *options) setInputs(v uint) error {
	o.inputs = v
	return nil
}

// Refresh is an option for New() that sets the VNC framebuffer refresh period.
func Refresh(v time.Duration) func(*options) error {
	return func(o *options) error { return o.setRefresh(v) }
}

// setRefresh sets the VNC framebuffer refresh period.
func (o *options) setRefresh(v time.Duration) error {
	o.refresh = v
	return nil
}
