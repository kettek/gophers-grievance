package resources

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

//go:embed maps/*
//go:embed ui/*
//go:embed tiles/*
var f embed.FS

var (
	BoxImage         *ebiten.Image
	BoulderImage     *ebiten.Image
	SolidImage       *ebiten.Image
	GopherImage      *ebiten.Image
	GopherRipImage   *ebiten.Image
	SnakeImage       *ebiten.Image
	SnakeSnoozeImage *ebiten.Image
	FoodImage        *ebiten.Image
	TimeImage        *ebiten.Image
	TimeBorderImage  *ebiten.Image
)

func loadImages() error {
	data, err := f.ReadFile("tiles/box.png")
	if err != nil {
		return err
	}
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}
	BoxImage, err = ebiten.NewImageFromImage(img, ebiten.FilterNearest)
	if err != nil {
		return err
	}

	data, _ = f.ReadFile("tiles/solid.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	SolidImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/boulder.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	BoulderImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/gopher.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	GopherImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/gopher-rip.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	GopherRipImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/snake.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	SnakeImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/snake-snooze.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	SnakeSnoozeImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/plant.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	FoodImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/time.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	TimeImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("tiles/time-border.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	TimeBorderImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	return nil
}
