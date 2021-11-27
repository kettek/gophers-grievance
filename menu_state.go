package main

import (
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
	"github.com/kettek/gophers-grievance/ui"
)

type MenuState struct {
	game    *Game
	buttons []ui.WidgetI
	ui      *ui.Manager
}

func (s *MenuState) init() error {
	// Set up UI
	s.buttons = []ui.WidgetI{
		ui.MakeButton("New Game", func(w ui.WidgetI) bool {
			gameState := &GameState{
				game: s.game,
				ui:   s.ui,
			}
			s.game.setState(gameState)
			gameState.reset()
			gameState.loadMap(resources.GetAnyMap())

			return true
		}),
		ui.MakeButton("Map Editor", func(w ui.WidgetI) bool {
			editorState := &EditorState{
				game: s.game,
				ui:   s.ui,
			}
			s.game.setState(editorState)
			return true
		}),
		ui.MakeButton("Exit", func(w ui.WidgetI) bool {
			os.Exit(0)
			return true
		}),
	}

	return nil
}

func (s *MenuState) update(screen *ebiten.Image) error {
	s.ui.CheckWidgets(s.buttons)
	return nil
}

func (s *MenuState) draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 125, 156, 255})

	// Ignore this shameful hack to get our button sizes before really rendering.
	btnX := 0
	btnY := 0
	padding := resources.ButtonMiddleImage.Bounds().Dy()
	for _, btn := range s.buttons {
		btn.Draw(screen, btnX, btnY)
		btnY += btn.Height() + padding
	}
	maxH := btnY

	screen.Fill(color.RGBA{0, 125, 156, 255})
	btnY = winHeight/2 - maxH/2
	btnX = winWidth / 2
	for _, btn := range s.buttons {
		btn.Draw(screen, btnX-btn.Width()/2, btnY-btn.Height()/2)
		btnY += btn.Height() + padding
	}
}
