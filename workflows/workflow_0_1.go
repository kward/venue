package workflow

func init() {
	const (
		version = "0.1"
		page    = "soundcheck"
	)

	for _, layout := range []string{"th"} {
		var (
			// The dx and dy vars are always based on a vertical orientation.
			dxInSelect, dyInSelect   int
			dxOutSelect, dyOutSelect int
			dxOutLevel, dyOutLevel   int
			orientation              int
		)

		switch layout {
		case "pv": // phone/vertical
			dxInSelect, dyInSelect = 8, 4
			dxOutSelect, dyOutSelect = 6, 1
			dxOutLevel, dyOutLevel = 4, 4
			orientation = vertical
		case "th": // tablet/horizontal
			dxInSelect, dyInSelect = 12, 4
			dxOutSelect, dyOutSelect = 12, 1
			dxOutLevel, dyOutLevel = 12, 4
			orientation = horizontal
		}

		// Input select.
		for x := 1; x < dxInSelect; x++ {
			for y := 1; y < dyInSelect; y++ {
				fx := []Control{inSelectVertical, inSelectHorizontal}
				Register(Workflow{
					version: version,
					layout:  layout,
					page:    page,
					control: "input",
					verb:    "select",
					fx:      fx[orientation]})
			}
		}
		// Output select.
		for x := 1; x < dxOutSelect; x++ {
			for y := 1; y < dyOutSelect; y++ {
				fx := []Control{outSelectVertical, outSelectHorizontal}
				Register(Workflow{
					version: version,
					layout:  layout,
					page:    page,
					control: "output",
					verb:    "select",
					fx:      fx[orientation]})
			}
		}
		// Output level.
		for x := 1; x < dxOutLevel; x++ {
			for y := 1; y < dyOutLevel; y++ {
				fx := []Control{outLevelVertical, outLevelHorizontal}
				Register(Workflow{
					version: version,
					layout:  layout,
					page:    page,
					control: "output",
					verb:    "level",
					fx:      fx[orientation]})
			}
		}
	}
}

func inSelectVertical(x, y, dx, dy, b int) {
}

func inSelectHorizontal(x, y, dx, dy, b int) {
	return
}

func outSelectVertical(x, y, dx, dy, b int) {
	return
}

func outSelectHorizontal(x, y, _, dy, _ int) {
	// TODO(kward:20161120) This doesn't do anything yet.
	x, y = multiRotate(x, y, dy)
}

func outLevelVertical(x, y, dx, dy, b int) {
	return
}

func outLevelHorizontal(x, y, dx, dy, b int) {
	return
}
