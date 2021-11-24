package main

import "github.com/hajimehoshi/ebiten"

func main() {
	ebiten.SetWindowSize(winWidth*2, winHeight*2)
	ebiten.SetWindowTitle(winTitle)

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
