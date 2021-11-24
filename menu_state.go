package main

import (
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
)

type MenuState struct {
	game *Game
}

func (s *MenuState) update(screen *ebiten.Image) error {
	// Just bump to game for now.
	f := Field{
		background: color.RGBA{128, 128, 128, 255},
	}
	f.fromMap(resources.GetAnyMap())
	s.game.setState(&GameState{
		game:         s.game,
		dirs:         make(map[Direction]struct{}),
		field:        f,
		turnTime:     50 * time.Millisecond,
		lastTurnTime: time.Now(),
	})
	return nil
}

func (s *MenuState) draw(screen *ebiten.Image) {}