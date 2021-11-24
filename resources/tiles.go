package resources

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

//go:embed *
//go:embed maps/*
var f embed.FS

var (
	BoxImage         *ebiten.Image
	BoulderImage     *ebiten.Image
	SolidImage       *ebiten.Image
	GopherImage      *ebiten.Image
	SnakeImage       *ebiten.Image
	SnakeSnoozeImage *ebiten.Image
	FoodImage        *ebiten.Image
)

func loadImages() error {
	data, err := f.ReadFile("box.png")
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

	data, _ = f.ReadFile("solid.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	SolidImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("boulder.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	BoulderImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("gopher.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	GopherImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("snake.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	SnakeImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("snake-snooze.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	SnakeSnoozeImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	data, _ = f.ReadFile("plant.png")
	img, _, _ = image.Decode(bytes.NewReader(data))
	FoodImage, _ = ebiten.NewImageFromImage(img, ebiten.FilterNearest)

	return nil
}
