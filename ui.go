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

	bankX        = 8   // X position of 1st bank.
	bankDX       = 131 // dX between banks.
	chanDX       = 15  // dX between channels in a bank.
	maxArrowKeys = 17  // Max number of consecutive arrow key presses.
)

// The VenueUI interface is the conventional interface for interacting with a
// Venue UI element.
type UIElement interface {
	// Read reads the current state of a UI element.
	Read(v *Venue) error

	// Set value of UI element.
	Set(v *Venue, val int)

	// Update updates the state of a UI element.
	Update(v *Venue) error
}

type VenuePages map[int]*Page

// Page returns the current page.
func (v *Venue) Page() int {
	return v.currPage
}

// SetPage selects the requested page for interaction.
func (v *Venue) SetPage(page int) error {
	pageNames := map[int]string{
		InputsPage:    "INPUTS",
		OutputsPage:   "OUTPUTS",
		FilingPage:    "FILING",
		SnapshotsPage: "SNAPSHOTS",
		PatchbayPage:  "PATCHBAY",
		PluginsPage:   "PLUG-INS",
		OptionsPage:   "OPTIONS",
	}

	if v.currPage != page {
		log.Printf("Changing to %v page.", pageNames[page])
		if err := v.KeyPress(vnc.KeyF1 + uint32(page)); err != nil {
			log.Println("Page() error:", err)
			return err
		}
		v.currPage = page
	}
	return nil
}

// Input selects the requested input for interaction.
func (v *Venue) SetInput(input int) error {
	if input < 1 || input > numInputs {
		err := fmt.Errorf("Input() invalid input: %v", input)
		log.Println(err)
		return err
	}
	log.Printf("Selecting input #%v.", input)

	v.SetPage(InputsPage)

	if v.currInput == nil {
		v.selectInput(1)
		v.currInput = v.inputs[0]
	}
	if input == v.currInput.ch {
		return nil
	}

	const (
		left  = false
		right = true
	)

	dist := input - v.currInput.ch
	kp := abs(dist)
	if kp <= maxArrowKeys {
		// Move with the keyboard.
		dir := left
		if dist > 0 {
			dir = right
		}
		for i := 0; i < kp; i++ {
			if dir == left {
				v.KeyPress(vnc.KeyLeft)
			} else {
				v.KeyPress(vnc.KeyRight)
			}
		}
	} else {
		if err := v.selectInput(input); err != nil {
			return err
		}
	}

	v.currInput = v.inputs[input-1]
	return nil
}

// selectInput select an input directly.
func (v *Venue) selectInput(input int) error {
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

// SetOutput selects the specified output for interaction.
func (v *Venue) SetOutput(name string) error {
	v.SetPage(OutputsPage)
	vp := v.Pages[OutputsPage]

	// Clear solo.
	log.Println("Clearing solo.")
	e := vp.Elements["solo_clear"]
	e.(*Switch).Update(v)

	// Solo output.
	log.Printf("Soloing %v output.", name)
	solo := name + "solo"
	e = vp.Elements[solo]
	e.(*Switch).Update(v)

	v.currOutput = v.outputs[name]
	return nil
}

// Page holds the UI elements of a VENUE page.
type Page struct {
	Elements map[string]UIElement
}

// NewInputsPage returns a populated Inputs page.
func NewInputsPage() *Page {
	const (
		auxOddX  = 316
		auxPanX  = 473
		aux12Y   = 95
		aux34Y   = 146
		aux56Y   = 197
		aux78Y   = 248
		aux910Y  = 299
		aux1112Y = 350
		aux1314Y = 401
		aux1516Y = 452
	)

	return &Page{
		Elements: map[string]UIElement{
			"gain":       &Encoder{image.Point{167, 279}, EncoderBL, true},
			"delay":      &Encoder{image.Point{168, 387}, EncoderBL, false},
			"hpf":        &Encoder{image.Point{168, 454}, EncoderBL, true},
			"pan":        &Encoder{image.Point{239, 443}, EncoderBC, false},
			"aux1":       &Encoder{image.Point{auxOddX, aux12Y}, EncoderTR, true},
			"aux12pan":   &Encoder{image.Point{auxPanX, aux12Y}, EncoderTL, false},
			"aux3":       &Encoder{image.Point{auxOddX, aux34Y}, EncoderTR, true},
			"aux34pan":   &Encoder{image.Point{auxPanX, aux34Y}, EncoderTL, false},
			"aux5":       &Encoder{image.Point{auxOddX, aux56Y}, EncoderTR, true},
			"aux56pan":   &Encoder{image.Point{auxPanX, aux56Y}, EncoderTL, false},
			"aux7":       &Encoder{image.Point{auxOddX, aux78Y}, EncoderTR, true},
			"aux78pan":   &Encoder{image.Point{auxPanX, aux78Y}, EncoderTL, false},
			"aux9":       &Encoder{image.Point{auxOddX, aux910Y}, EncoderTR, true},
			"aux910pan":  &Encoder{image.Point{auxPanX, aux910Y}, EncoderTL, false},
			"aux11":      &Encoder{image.Point{auxOddX, aux1112Y}, EncoderTR, true},
			"aux1112pan": &Encoder{image.Point{auxPanX, aux1112Y}, EncoderTL, false},
			"aux13":      &Encoder{image.Point{auxOddX, aux1314Y}, EncoderTR, true},
			"aux1314pan": &Encoder{image.Point{auxPanX, aux1314Y}, EncoderTL, false},
			"aux15":      &Encoder{image.Point{auxOddX, aux1516Y}, EncoderTR, true},
			"aux1516pan": &Encoder{image.Point{auxPanX, aux1516Y}, EncoderTL, false},
			"grp1":       &Encoder{image.Point{auxOddX, aux12Y}, EncoderTR, true},
			"grp12pan":   &Encoder{image.Point{auxPanX, aux12Y}, EncoderTL, false},
			"grp3":       &Encoder{image.Point{auxOddX, aux34Y}, EncoderTR, true},
			"grp34pan":   &Encoder{image.Point{auxPanX, aux34Y}, EncoderTL, false},
			"grp5":       &Encoder{image.Point{auxOddX, aux56Y}, EncoderTR, true},
			"grp56pan":   &Encoder{image.Point{auxPanX, aux56Y}, EncoderTL, false},
			"grp7":       &Encoder{image.Point{auxOddX, aux78Y}, EncoderTR, true},
			"grp78pan":   &Encoder{image.Point{auxPanX, aux78Y}, EncoderTL, false},
			"solo_clear": newPushButton(980, 490, mediumSwitch),
		},
	}
}

// NewOutputsPage returns a populated Outputs page.
func NewOutputsPage() *Page {
	const (
		soloY = 573
	)

	var b int // Bank
	elements := map[string]UIElement{
		"solo_clear": newPushButton(980, 490, mediumSwitch),
	}

	// Auxes
	pre, post := "aux", "solo"
	for b = 1; b <= 2; b++ { // bank
		for c := 1; c <= 8; c++ { // channel
			n := fmt.Sprintf("%v%v%v", pre, (b-1)*8+c, post)
			elements[n] = newToggle(bankX+(b-1)*bankDX+(c-1)*chanDX, soloY, tinySwitch, false)
		}
	}
	// Groups
	b, pre, post = 5, "grp", "solo"
	for c := 1; c <= 8; c++ { // channel (only 1 bank)
		n := fmt.Sprintf("%v%v%v", pre, c, post)
		elements[n] = newToggle(bankX+(b-1)*bankDX+(c-1)*chanDX, soloY, tinySwitch, false)
	}

	return &Page{Elements: elements}
}
