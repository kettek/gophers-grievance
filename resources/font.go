package resources

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var (
	NormalFont font.Face
	BoldFont   font.Face
)

func loadFonts() error {
	data, err := f.ReadFile("ui/OpenSansPX.ttf")
	if err != nil {
		return err
	}
	tt, err := opentype.Parse(data)
	if err != nil {
		return err
	}

	const dpi = 72
	NormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	data, err = f.ReadFile("ui/OpenSansPXBold.ttf")
	if err != nil {
		return err
	}
	tt, err = opentype.Parse(data)
	if err != nil {
		return err
	}

	BoldFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    16,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	return nil
}
