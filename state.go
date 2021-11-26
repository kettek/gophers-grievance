package main

import "github.com/hajimehoshi/ebiten"

type State interface {
	init() error
	update(screen *ebiten.Image) error
	draw(screen *ebiten.Image)
}
