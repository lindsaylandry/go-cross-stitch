package palette

func GreyPalette() ([]Thread, error) {
	return palette("../../palette/black-white-grey.csv")
}

func DMCPalette() ([]Thread, error) {
	return palette("../../palette/dmc-floss.csv")
}

func palette(path string) ([]Thread, error) {
	return dmcColors, nil
}
