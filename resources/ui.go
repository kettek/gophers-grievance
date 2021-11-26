package resources

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten"
)

var (
	ButtonLeftImage        *ebiten.Image
	ButtonRightImage       *ebiten.Image
	ButtonMiddleImage      *ebiten.Image
	ButtonLeftPressImage   *ebiten.Image
	ButtonRightPressImage  *ebiten.Image
	ButtonMiddlePressImage *ebiten.Image
)

func loadUi() error {
	data, _ := f.ReadFile("ui/button-left.png")
	img, _, _ := image.Decode(bytes.NewReader(data))
	ButtonLeftImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("ui/button-right.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	ButtonRightImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("ui/button-middle.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	ButtonMiddleImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("ui/button-left-press.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	ButtonLeftPressImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("ui/button-right-press.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	ButtonRightPressImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("ui/button-middle-press.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	ButtonMiddlePressImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	return nil
}
