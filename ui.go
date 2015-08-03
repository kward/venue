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

	// Select the UI element.
	Select(v *Venue)

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

	elements := map[string]UIElement{
		"gain":       &Encoder{image.Point{167, 279}, encoderBL, true},
		"delay":      &Encoder{image.Point{168, 387}, encoderBL, false},
		"hpf":        &Encoder{image.Point{168, 454}, encoderBL, true},
		"pan":        &Encoder{image.Point{239, 443}, encoderBC, false},
		"aux1":       &Encoder{image.Point{auxOddX, aux12Y}, encoderTR, true},
		"aux1pan":    &Encoder{image.Point{auxPanX, aux12Y}, encoderTL, false},
		"aux3":       &Encoder{image.Point{auxOddX, aux34Y}, encoderTR, true},
		"aux3pan":    &Encoder{image.Point{auxPanX, aux34Y}, encoderTL, false},
		"aux5":       &Encoder{image.Point{auxOddX, aux56Y}, encoderTR, true},
		"aux5pan":    &Encoder{image.Point{auxPanX, aux56Y}, encoderTL, false},
		"aux7":       &Encoder{image.Point{auxOddX, aux78Y}, encoderTR, true},
		"aux7pan":    &Encoder{image.Point{auxPanX, aux78Y}, encoderTL, false},
		"aux9":       &Encoder{image.Point{auxOddX, aux910Y}, encoderTR, true},
		"aux9pan":    &Encoder{image.Point{auxPanX, aux910Y}, encoderTL, false},
		"aux11":      &Encoder{image.Point{auxOddX, aux1112Y}, encoderTR, true},
		"aux11pan":   &Encoder{image.Point{auxPanX, aux1112Y}, encoderTL, false},
		"aux13":      &Encoder{image.Point{auxOddX, aux1314Y}, encoderTR, true},
		"aux13pan":   &Encoder{image.Point{auxPanX, aux1314Y}, encoderTL, false},
		"aux15":      &Encoder{image.Point{auxOddX, aux1516Y}, encoderTR, true},
		"aux15pan":   &Encoder{image.Point{auxPanX, aux1516Y}, encoderTL, false},
		"grp1":       &Encoder{image.Point{auxOddX, aux12Y}, encoderTR, true},
		"grp1pan":    &Encoder{image.Point{auxPanX, aux12Y}, encoderTL, false},
		"grp3":       &Encoder{image.Point{auxOddX, aux34Y}, encoderTR, true},
		"grp3pan":    &Encoder{image.Point{auxPanX, aux34Y}, encoderTL, false},
		"grp5":       &Encoder{image.Point{auxOddX, aux56Y}, encoderTR, true},
		"grp5pan":    &Encoder{image.Point{auxPanX, aux56Y}, encoderTL, false},
		"grp7":       &Encoder{image.Point{auxOddX, aux78Y}, encoderTR, true},
		"grp7pan":    &Encoder{image.Point{auxPanX, aux78Y}, encoderTL, false},
		"solo_clear": newPushButton(980, 490, mediumSwitch),
	}

	return &Page{Elements: elements}
}

// NewOutputsPage returns a populated Outputs page.
func NewOutputsPage() *Page {
	const (
		meterY = 512
		muteY  = 588
		soloY  = 573
	)

	elements := map[string]UIElement{
		"solo_clear": newPushButton(980, 490, mediumSwitch),
	}

	// Auxes & Groups
	for _, b := range []int{1, 2, 5} { // bank
		var pre string
		switch b {
		case 5:
			pre = "grp"
		default:
			pre = "aux"
		}
		for c := 1; c <= 8; c++ { // channel
			ch, x := (b-1)*8+c, bankX+(b-1)*bankDX+(c-1)*chanDX

			n := fmt.Sprintf("%v%vsolo", pre, ch)
			elements[n] = newToggle(x, soloY, tinySwitch, false)

			n = fmt.Sprintf("%v%vmeter", pre, ch)
			elements[n] = &Meter{
				pos:  image.Point{x, meterY},
				size: smallVMeter,
			}
		}
	}

	return &Page{Elements: elements}
}
