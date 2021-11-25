package main

import (
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
)

type MenuState struct {
	game *Game
}

func (s *MenuState) update(screen *ebiten.Image) error {
	// Just bump to game for now.
	p := Player{
		dirs:  make(map[Direction]struct{}),
		lives: maxLives,
	}
	backgroundImage, err := ebiten.NewImage(276, 276, ebiten.FilterLinear)
	if err != nil {
		return err
	}
	// TODO: Move to a game state based map load sort of deal.
	gameState := &GameState{
		game:            s.game,
		turnTime:        50 * time.Millisecond,
		difficulty:      5,
		players:         []Player{p},
		backgroundImage: backgroundImage,
	}
	gameState.reset()
	gameState.loadMap(resources.GetAnyMap())
	s.game.setState(gameState)
	return nil
}

func (s *MenuState) draw(screen *ebiten.Image) {}
