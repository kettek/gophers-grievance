package main

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/kettek/gophers-grievance/resources"
)

type Button struct {
	t    string
	x, y int
	w, h int
	held bool
	cb   func() bool
}

func (b *Button) draw(screen *ebiten.Image, x, y int) (int, int) {
	b.x = x
	b.y = y
	op := &ebiten.DrawImageOptions{}

	r := text.BoundString(resources.NormalFont, b.t)

	fx := x

	var leftImage *ebiten.Image
	var middleImage *ebiten.Image
	var rightImage *ebiten.Image
	if !b.held {
		leftImage = resources.ButtonLeftImage
		middleImage = resources.ButtonMiddleImage
		rightImage = resources.ButtonRightImage
	} else {
		leftImage = resources.ButtonLeftPressImage
		middleImage = resources.ButtonMiddlePressImage
		rightImage = resources.ButtonRightPressImage
	}

	// Draw button.
	op.GeoM.Translate(float64(fx), float64(y))
	screen.DrawImage(leftImage, op)

	fx += leftImage.Bounds().Max.X

	canFinish := func() bool {
		return int(fx-x)+rightImage.Bounds().Dx() > r.Dx()+4
	}

	if canFinish() {
		op.GeoM.Reset()
		op.GeoM.Translate(float64(fx), float64(y))
		screen.DrawImage(rightImage, op)
		fx += middleImage.Bounds().Max.X
	} else {
		for (fx - x) < r.Dx()+4 {
			op.GeoM.Reset()
			op.GeoM.Translate(float64(fx), float64(y))
			if canFinish() {
				screen.DrawImage(rightImage, op)
				fx += rightImage.Bounds().Max.X
			} else {
				screen.DrawImage(middleImage, op)
				fx += middleImage.Bounds().Max.X
			}
		}
	}

	w := fx - x
	h := middleImage.Bounds().Dy()

	textX := x
	textX += w / 2
	textX -= r.Dx()/2 + r.Min.X

	textY := y + 1
	textY += h/2 + 2
	textY += r.Min.Y / 2
	textY += r.Dy() / 2

	// Draw text.
	text.Draw(screen, b.t, resources.NormalFont, textX, textY, color.Black)
	b.w = w
	b.h = h
	return w, h
}

func (b *Button) hit(x, y int) bool {
	return x >= b.x && x <= b.x+b.w && y >= b.y && y <= b.y+b.h
}

type UiManager struct {
	mouseState map[int]bool
	//buttons    []Button
}

func (ui *UiManager) checkButtons(buttons []Button) bool {
	// Check for UI interactions.
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ui.mouseState[int(ebiten.MouseButtonLeft)] {
		x, y := ebiten.CursorPosition()
		for i := range buttons {
			btn := &buttons[i]
			if btn.hit(x, y) && btn.held {
				if btn.cb() {
					return true
				}
			}
			btn.held = false
		}
		ui.mouseState[int(ebiten.MouseButtonLeft)] = false
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !ui.mouseState[int(ebiten.MouseButtonLeft)] {
		ui.mouseState[int(ebiten.MouseButtonLeft)] = true
		x, y := ebiten.CursorPosition()
		for i := range buttons {
			btn := &buttons[i]
			if btn.hit(x, y) {
				btn.held = true
			}
		}
	}
	return false
}
