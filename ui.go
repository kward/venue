package venue

import (
	"fmt"
	"image"
	"log"
	"math"
	"time"

	vnc "github.com/kward/go-vnc"
)

const (
	InputsPage = iota
	OutputsPage
	FilingPage
	SnapshotsPage
	PatchbayPage
	PluginsPage
	OptionsPage
)
const (
	bankX  = 8   // X position of 1st bank.
	bankDX = 131 // dX between banks.
	chanDX = 15  // dX between channels in a bank.
)

// The VenueUI interface is the conventional interface for interacting with a
// Venue UI element.
type UIElement interface {
	// Read reads the current state of a UI element.
	Read(*Venue) error

	// Update updates the state of a UI element.
	Update(*Venue) error
}

// Page selects the requested page.
func (v *Venue) Page(page int) error {
	log.Println("page:", page)
	if err := v.KeyPress(vnc.KeyF1 + uint32(page)); err != nil {
		log.Println("Page() error:", err)
		return err
	}
	time.Sleep(uiSettle)
	return nil
}

func (v *Venue) Input(input int) error {
	if input < 1 || input > numInputs {
		err := fmt.Errorf("invalid input: %v", input)
		log.Println(err)
		return err
	}

	log.Println("input:", input)
	digit, _ := math.Modf(float64(input) / 10)
	if err := v.KeyPress(vnc.Key0 + uint32(digit)); err != nil {
		log.Println("Input() error on 1st key press:", err)
		return err
	}
	digit = math.Mod(float64(input), 10.0)
	if err := v.KeyPress(vnc.Key0 + uint32(digit)); err != nil {
		log.Println("Input() error on 2nd key press:", err)
		return err
	}
	// TODO(kward): This may slow a user down. maybe only delay once the "search"
	// has completed?
	time.Sleep(1750 * time.Millisecond)

	return nil
}

func (v *Venue) Output(output int) error {
	return nil
}

type Page struct {
	Elements map[string]UIElement
}

func NewInputsPage() *Page {
	const (
		auxOddX = 316
		//aux_even_x :=
		auxPanX = 473
	)

	return &Page{
		Elements: map[string]UIElement{
			"gain":       &Encoder{image.Point{167, 279}, EncoderBL, true},
			"delay":      &Encoder{image.Point{168, 387}, EncoderBL, false},
			"hpf":        &Encoder{image.Point{168, 454}, EncoderBL, true},
			"pan":        &Encoder{image.Point{239, 443}, EncoderBC, false},
			"aux1":       &Encoder{image.Point{auxOddX, 95}, EncoderTR, true},
			"aux12pan":   &Encoder{image.Point{auxPanX, 95}, EncoderTL, false},
			"aux3":       &Encoder{image.Point{auxOddX, 146}, EncoderTR, true},
			"aux34pan":   &Encoder{image.Point{auxPanX, 146}, EncoderTL, false},
			"aux5":       &Encoder{image.Point{auxOddX, 197}, EncoderTR, true},
			"aux56pan":   &Encoder{image.Point{auxPanX, 197}, EncoderTL, false},
			"aux7":       &Encoder{image.Point{auxOddX, 248}, EncoderTR, true},
			"aux78pan":   &Encoder{image.Point{auxPanX, 248}, EncoderTL, false},
			"aux9":       &Encoder{image.Point{auxOddX, 299}, EncoderTR, true},
			"aux910pan":  &Encoder{image.Point{auxPanX, 299}, EncoderTL, false},
			"aux11":      &Encoder{image.Point{auxOddX, 350}, EncoderTR, true},
			"aux1112pan": &Encoder{image.Point{auxPanX, 350}, EncoderTL, false},
			"aux13":      &Encoder{image.Point{auxOddX, 401}, EncoderTR, true},
			"aux1314pan": &Encoder{image.Point{auxPanX, 401}, EncoderTL, false},
			"aux15":      &Encoder{image.Point{auxOddX, 452}, EncoderTR, true},
			"aux1516pan": &Encoder{image.Point{auxPanX, 452}, EncoderTL, false},
		},
	}
}

func NewOutputsPage() *Page {
	const (
		soloY = 573
	)

	var b int // Bank
	elements := map[string]UIElement{}

	// Auxes
	pre, post := "aux", "solo"
	for b = 1; b <= 2; b++ { // bank
		for c := 1; c <= 8; c++ { // channel
			n := fmt.Sprintf("%v%v%v", pre, (b-1)*8+c, post)
			elements[n] = newSwitch(bankX+(b-1)*bankDX+(c-1)*chanDX, soloY, tinySwitch, false, false)
		}
	}
	// Groups
	b, pre, post = 5, "grp", "solo"
	for c := 1; c <= 8; c++ { // channel (only 1 bank)
		n := fmt.Sprintf("%v%v%v", pre, c, post)
		elements[n] = newSwitch(bankX+(b-1)*bankDX+(c-1)*chanDX, soloY, tinySwitch, false, false)
	}

	return &Page{Elements: elements}
}
