package main

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
)

type MenuState struct {
	game    *Game
	buttons []Button
	ui      *UiManager
}

func (s *MenuState) init() error {
	s.buttons = []Button{
		{
			t: "New Game",
			cb: func() bool {
				gameState := &GameState{
					game: s.game,
					ui:   s.ui,
				}
				s.game.setState(gameState)
				gameState.reset()
				gameState.loadMap(resources.GetAnyMap())

				return true
			},
		},
		{
			t: "Exit",
			cb: func() bool {
				os.Exit(0)
				return true
			},
		},
	}
	return nil
}

func (s *MenuState) update(screen *ebiten.Image) error {
	s.ui.checkButtons(s.buttons)
	return nil
}

func (s *MenuState) draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 125, 156, 255})

	// Ignore this shameful hack to get our button sizes before really rendering.
	btnX := 0
	btnY := 0
	padding := resources.ButtonMiddleImage.Bounds().Dy()
	for i := range s.buttons {
		btn := &s.buttons[i]
		btn.draw(screen, btnX, btnY)
		btnY += btn.h + padding
	}
	maxH := btnY

	screen.Fill(color.RGBA{0, 125, 156, 255})
	btnY = winHeight/2 - maxH/2
	btnX = winWidth / 2
	for i := range s.buttons {
		btn := &s.buttons[i]
		btn.draw(screen, btnX-btn.w/2, btnY-btn.h/2)
		btnY += btn.h + padding
	}
}
