package main

import (
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
	"github.com/kettek/gophers-grievance/ui"
)

type EditorState struct {
	game       *Game
	buttons    []ui.WidgetI
	ui         *ui.Manager
	currentMap resources.Map
}

func (s *EditorState) init() error {
	s.buttons = []ui.WidgetI{
		ui.MakeButton("Main Menu", func() bool {
			s.game.setState(&MenuState{
				game: s.game,
				ui:   s.ui,
			})
			return true
		}),
		ui.MakeButton("New", func() bool {
			return false
		}),
		ui.MakeButton("Save", func() bool {
			return false
		}),
		ui.MakeButton("Load", func() bool {
			return false
		}),
	}

	// Create base map.
	s.currentMap = resources.Map{
		Name:       "Gopher's Paradise",
		Background: color.RGBA{0, 125, 156, 255},
		Columns:    23,
		Rows:       23,
	}
	s.currentMap.Cells = make([]string, s.currentMap.Rows)
	for i := 0; i < s.currentMap.Rows; i++ {
		s.currentMap.Cells[i] = strings.Repeat("#", s.currentMap.Columns)
	}

	return nil
}

func (s *EditorState) update(screen *ebiten.Image) error {
	s.ui.CheckWidgets(s.buttons)
	return nil
}

func (s *EditorState) draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 125, 156, 255})
	// Draw Top Buttons
	offsetX := 0
	padding := 2
	btnX := padding
	btnY := padding
	for _, btn := range s.buttons {
		btn.Draw(screen, btnX, btnY)
		btnX += btn.Width() + padding
	}

	var offsetY float64 = 332 - 276 // for now

	// Draw Tile/Object buttons

	// Draw map
	op := &ebiten.DrawImageOptions{}
	for y, row := range s.currentMap.Cells {
		for x, r := range row {
			s, ok := resources.Runes[r]
			if !ok {
				continue
			}
			op.GeoM.Reset()
			op.GeoM.Translate(float64(offsetX)+float64(x)*tileWidth, float64(offsetY)+float64(y)*tileHeight)
			screen.DrawImage(resources.Images[s], op)
		}
	}
}
