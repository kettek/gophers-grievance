package ui

import (
	"github.com/hajimehoshi/ebiten"
)

type Image struct {
	Widget
	img *ebiten.Image
}

func MakeImage(img *ebiten.Image, cb func(WidgetI) bool) WidgetI {
	return &Image{
		Widget: Widget{
			cb: cb,
		},
		img: img,
	}
}

func (w *Image) Draw(screen *ebiten.Image, x, y int) (int, int) {
	w.x = x
	w.y = y
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(w.img, op)

	w.w = w.img.Bounds().Dx()
	w.h = w.img.Bounds().Dy()
	return w.w, w.h
}
