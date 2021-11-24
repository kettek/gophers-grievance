package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
)

type LoadState struct {
	game *Game
}

func (s *LoadState) update(screen *ebiten.Image) error {
	fmt.Println("resources load")
	if err := resources.Load(); err != nil {
		panic(err)
	}
	s.game.setState(&MenuState{
		game: s.game,
	})
	return nil
}

func (s *LoadState) draw(screen *ebiten.Image) {}
