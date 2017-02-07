package venuelib

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/golang/glog"
	"github.com/howeyc/gopass"
	"github.com/kward/venue/codes"
)

// venueError defines the status of a Venue call.
type venueError struct {
	code codes.Code
	desc string
}

func (e *venueError) Error() string {
	return fmt.Sprintf("venue error: %s: %s", e.code, e.desc)
}

// Code returns the error code for `err` if it was produced by Venue.
// Otherwise, it returns codes.Unknown.
func Code(err error) codes.Code {
	if err == nil {
		return codes.OK
	}
	if e, ok := err.(*venueError); ok {
		return e.code
	}
	return codes.Unknown
}

// ErrorDesc returns the error description of `err` if it was produced by Venue.
// Otherwise, it returns err.Error(), or an empty string when `err` is nil.
func ErrorDesc(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(*venueError); ok {
		return e.desc
	}
	return err.Error()
}

// Errorf returns an error containing an error code and a description.
// Errorf returns nil if `c` is OK.
func Errorf(c codes.Code, format string, a ...interface{}) error {
	if c == codes.OK {
		return nil
	}
	return &venueError{
		code: c,
		desc: fmt.Sprintf(format, a...),
	}
}

// GetPasswd requests the user for a masked password.
func GetPasswd() string {
	fmt.Printf("Password: ")
	p, err := gopass.GetPasswdMasked()
	if err != nil {
		glog.Fatal(err)
	}
	return string(p)
}

// ToInt converts a string to an int.
func ToInt(s string) int {
	var i int
	_, err := fmt.Fscanf(bytes.NewBufferString(s), "%d", &i)
	if err != nil {
		return 0
	}
	return i
}

// FnName returns the calling function name, e.g. "SomeFunction()". This is
// useful for logging the function name with glog.
func FnName() string {
	pc := make([]uintptr, 10) // At least 1 entry needed.
	runtime.Callers(2, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	return name[strings.LastIndex(name, ".")+1:] + "()"
}
