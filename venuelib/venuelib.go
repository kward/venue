package venuelib

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/golang/glog"
	"github.com/howeyc/gopass"
)

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
