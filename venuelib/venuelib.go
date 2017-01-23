package venuelib

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"strings"

	"github.com/howeyc/gopass"
)

func GetPasswd() string {
	fmt.Printf("Password: ")
	p, err := gopass.GetPasswdMasked()
	if err != nil {
		log.Fatal(err)
	}
	return string(p)
}

// ToInt converts a string
func ToInt(s string) (i int) {
	_, err := fmt.Fscanf(bytes.NewBufferString(s), "%d", &i)
	if err != nil {
		return 0
	}
	return
}

// FnName returns the calling function name, e.g. "SomeFunction()". This is
// useful for logging the function name with glog.
func FnName() string {
	pc := make([]uintptr, 10) // At least 1 entry needed.
	runtime.Callers(2, pc)
	name := runtime.FuncForPC(pc[0]).Name()
	return name[strings.LastIndex(name, ".")+1:] + "()"
}
