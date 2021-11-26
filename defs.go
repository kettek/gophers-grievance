package main

import (
	"image/color"
	"time"
)

const tileWidth float64 = 12
const tileHeight float64 = 12

const winWidth int = 276
const winHeight int = 332
const scale int = 2
const winTitle string = "Gopher's Grievance"

const floatingTextDuration = 2 * time.Second
const floatingTextAlpha uint8 = 128

var gopherColor = color.RGBA{
	90, 218, 255, 255,
}

const maxLives int = 3
