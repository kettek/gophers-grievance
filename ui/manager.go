package ui

import "github.com/hajimehoshi/ebiten"

type Manager struct {
	mouseState map[int]bool
}

func MakeManager() *Manager {
	return &Manager{
		mouseState: make(map[int]bool),
	}
}

func (ui *Manager) CheckWidgets(widgets []WidgetI) bool {
	// Check for UI interactions.
	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && ui.mouseState[int(ebiten.MouseButtonLeft)] {
		x, y := ebiten.CursorPosition()
		for _, w := range widgets {
			if w.Hit(x, y) && w.Held() {
				if w.Cb() {
					return true
				}
			}
			w.SetHeld(false)
		}
		ui.mouseState[int(ebiten.MouseButtonLeft)] = false
	} else if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !ui.mouseState[int(ebiten.MouseButtonLeft)] {
		ui.mouseState[int(ebiten.MouseButtonLeft)] = true
		x, y := ebiten.CursorPosition()
		for _, w := range widgets {
			if w.Hit(x, y) {
				w.SetHeld(true)
			}
		}
	}
	return false
}
