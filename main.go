package main

import "github.com/hajimehoshi/ebiten"

func main() {
	ebiten.SetWindowSize(winWidth*scale, winHeight*scale)
	ebiten.SetWindowTitle(winTitle)

	game := &Game{}

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
