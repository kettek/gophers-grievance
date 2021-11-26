package main

import (
	"github.com/hajimehoshi/ebiten"
)

type Game struct {
	state State
}

func (g *Game) Update(screen *ebiten.Image) error {
	if g.state == nil {
		g.setState(&LoadState{
			game: g,
		})
	}
	return g.state.update(screen)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.state.draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return winWidth, winHeight
}

func (g *Game) setState(state State) {
	state.init()
	g.state = state
}
