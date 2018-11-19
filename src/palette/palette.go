package palette

func GreyPalette() ([]Thread, error) {
	return palette("../../palette/black-white-grey.csv")
}

func DMCPalette() []Thread {
	return GetDMCColors()
}

func palette(path string) ([]Thread, error) {
	return GetDMCColors(), nil
}
