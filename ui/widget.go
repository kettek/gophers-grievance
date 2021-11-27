package ui

import (
	"github.com/hajimehoshi/ebiten"
)

type WidgetI interface {
	SetX(int)
	X() int
	SetY(int)
	Y() int
	Width() int
	Height() int
	SetHeld(bool)
	Held() bool
	SetCb(func(WidgetI) bool)
	Cb() bool
	SetId(string)
	Id() string
	//
	Draw(screen *ebiten.Image, x, y int) (int, int)
	Hit(x, y int) bool
}

type Widget struct {
	x, y int
	w, h int
	held bool
	id   string
	cb   func(WidgetI) bool
}

func (w *Widget) SetX(x int) {
	w.x = x
}

func (w *Widget) X() int {
	return w.x
}

func (w *Widget) SetY(y int) {
	w.y = y
}

func (w *Widget) Y() int {
	return w.y
}

func (w *Widget) Width() int {
	return w.w
}

func (w *Widget) Height() int {
	return w.h
}

func (w *Widget) SetHeld(b bool) {
	w.held = b
}

func (w *Widget) Held() bool {
	return w.held
}

func (w *Widget) SetCb(cb func(WidgetI) bool) {
	w.cb = cb
}

func (w *Widget) Cb() bool {
	if w.cb == nil {
		return false
	}
	return w.cb(w)
}

func (w *Widget) Draw(screen *ebiten.Image, x, y int) (int, int) {
	return 0, 0
}

func (w *Widget) Id() string {
	return w.id
}

func (w *Widget) SetId(id string) {
	w.id = id
}

func (w *Widget) Hit(x, y int) bool {
	return x >= w.x && x <= w.x+w.w && y >= w.y && y <= w.y+w.h
}
