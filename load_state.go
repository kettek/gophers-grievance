package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
	"github.com/kettek/gophers-grievance/ui"
)

type LoadState struct {
	game *Game
	ui   *ui.Manager
}

func (s *LoadState) init() error {
	s.ui = ui.MakeManager()

	return nil
}

func (s *LoadState) update(screen *ebiten.Image) error {
	fmt.Println("resources load")
	if err := resources.Load(); err != nil {
		panic(err)
	}
	gameState := &GameState{
		game: s.game,
		ui:   s.ui,
	}
	s.game.setState(gameState)
	gameState.reset()
	gameState.loadMap(resources.GetAnyMap())

	return nil
}

func (s *LoadState) draw(screen *ebiten.Image) {}
