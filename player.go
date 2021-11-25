package main

import (
	"time"

	"github.com/hajimehoshi/ebiten"
)

// Player controls a gopher.
type Player struct {
	moveDelayer time.Duration
	direction   Direction
	dirs        map[Direction]struct{}
	lives       int
	score       int
}

func (p *Player) update(last time.Time, current time.Time, delta time.Duration) {
	// TODO: Separate direction out into a player type.
	if _, ok := p.dirs[west]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyA) && !ebiten.IsKeyPressed(ebiten.KeyH) {
			delete(p.dirs, west)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyH) {
			p.direction = west
			p.dirs[west] = struct{}{}
		}
	}
	if _, ok := p.dirs[east]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyD) && !ebiten.IsKeyPressed(ebiten.KeyL) {
			delete(p.dirs, east)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyL) {
			p.direction = east
			p.dirs[east] = struct{}{}
		}
	}
	if _, ok := p.dirs[north]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyW) && !ebiten.IsKeyPressed(ebiten.KeyK) {
			delete(p.dirs, north)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyK) {
			p.direction = north
			p.dirs[north] = struct{}{}
		}
	}
	if _, ok := p.dirs[south]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyJ) {
			delete(p.dirs, south)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyJ) {
			p.direction = south
			p.dirs[south] = struct{}{}
		}
	}
}

func (p *Player) reduceLives() {
	p.lives--
}
