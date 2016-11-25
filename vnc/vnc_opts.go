package vnc

type options struct {
	host   string // VENUE VNC host.
	port   uint   // VENUE VNC port
	passwd string // VENUE VNC password.
}

// Host is an option for New() that sets the Venue VNC host.
func Host(v string) func(*options) error {
	return func(o *options) error { return o.setHost(v) }
}

// setHost sets the Venue VNC host.
func (o *options) setHost(v string) error {
	o.host = v
	return nil
}

// Password is an option for New() that sets the Venue VNC password.
func Password(v string) func(*options) error {
	return func(o *options) error { return o.setPassword(v) }
}

// setPassword sets the Venue VNC password.
func (o *options) setPassword(v string) error {
	o.passwd = v
	return nil
}

// Port is an option for New() that sets the Venue VNC port.
func Port(v uint) func(*options) error {
	return func(o *options) error { return o.setPort(v) }
}

// setPort sets the Venue VNC port.
func (o *options) setPort(v uint) error {
	o.port = v
	return nil
}
