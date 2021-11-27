package main

import (
	"fmt"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
	"github.com/kettek/gophers-grievance/ui"
)

type EditorState struct {
	game             *Game
	buttons          []ui.WidgetI
	tileWidgets      []ui.WidgetI
	ui               *ui.Manager
	selectedRune     rune
	cursorX, cursorY int
	currentMap       resources.Map
}

func (s *EditorState) init() error {
	s.selectedRune = '#'
	s.buttons = []ui.WidgetI{
		ui.MakeButton("Main Menu", func(w ui.WidgetI) bool {
			s.game.setState(&MenuState{
				game: s.game,
				ui:   s.ui,
			})
			return true
		}),
		ui.MakeButton("New", func(w ui.WidgetI) bool {
			return false
		}),
		ui.MakeButton("Save", func(w ui.WidgetI) bool {
			return false
		}),
		ui.MakeButton("Load", func(w ui.WidgetI) bool {
			return false
		}),
	}

	for _, r := range resources.TileRuneList {
		str, ok := resources.Runes[r]
		if !ok {
			continue
		}
		img, ok := resources.Images[str]
		if !ok {
			continue
		}
		w := ui.MakeImage(img, func(w ui.WidgetI) bool {
			s.selectedRune = rune(w.Id()[0])
			return false
		})
		w.SetId(fmt.Sprintf("%c", r))
		s.tileWidgets = append(s.tileWidgets, w)
	}

	// Create base map.
	s.currentMap = resources.Map{
		Name:       "Gopher's Paradise",
		Background: color.RGBA{0, 125, 156, 255},
		Columns:    23,
		Rows:       23,
	}
	s.currentMap.Cells = make([][]rune, s.currentMap.Rows)
	for i := 0; i < s.currentMap.Rows; i++ {
		if i == 0 || i == s.currentMap.Rows-1 {
			s.currentMap.Cells[i] = []rune(strings.Repeat("#", s.currentMap.Columns))
		} else {
			s.currentMap.Cells[i] = []rune(strings.Repeat(" ", s.currentMap.Columns))
			s.currentMap.Cells[i][0] = '#'
			s.currentMap.Cells[i][s.currentMap.Columns-1] = '#'
		}
	}

	return nil
}

func (s *EditorState) update(screen *ebiten.Image) error {
	s.ui.CheckWidgets(append(s.tileWidgets, s.buttons...))
	if s.cursorX >= 0 && s.cursorY >= 0 && s.cursorX < s.currentMap.Columns && s.cursorY < s.currentMap.Rows {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.currentMap.Cells[s.cursorY][s.cursorX] = s.selectedRune
		} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonRight) {
			s.currentMap.Cells[s.cursorY][s.cursorX] = ' '
		}
	}
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
	btnY += s.buttons[0].Height() + padding

	// Draw Tile/Object buttons
	btnX = padding
	for _, t := range s.tileWidgets {
		t.Draw(screen, btnX, btnY)
		btnX += t.Width() + padding
	}

	var offsetY float64 = 332 - 276 // for now

	// Draw map
	op := &ebiten.DrawImageOptions{}
	for y, row := range s.currentMap.Cells {
		for x, r := range row {
			s, ok := resources.Runes[r]
			if !ok {
				continue
			}
			img, ok := resources.Images[s]
			if !ok {
				continue
			}
			op.GeoM.Reset()
			op.GeoM.Translate(float64(offsetX)+float64(x)*tileWidth, float64(offsetY)+float64(y)*tileHeight)
			screen.DrawImage(img, op)
		}
	}

	// Draw cursor
	if str, ok := resources.Runes[s.selectedRune]; ok {
		if img, ok := resources.Images[str]; ok {
			cx, cy := ebiten.CursorPosition()
			cx -= offsetX
			cy -= int(offsetY)
			cx = int(float64(cx) / tileWidth)
			cy = int(float64(cy) / tileHeight)
			s.cursorX = cx
			s.cursorY = cy
			if cx >= 0 && cy >= 0 && cx < s.currentMap.Columns && cy < s.currentMap.Rows {
				op.GeoM.Reset()
				op.GeoM.Translate(float64(offsetX)+float64(cx)*tileWidth, offsetY+float64(cy)*tileHeight)
				screen.DrawImage(img, op)
			}
		}
	}
}
