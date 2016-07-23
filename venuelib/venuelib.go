package venuelib

import (
	"bytes"
	"fmt"
	"log"

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
