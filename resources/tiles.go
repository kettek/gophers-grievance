package resources

import (
	"bytes"
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten"
)

//go:embed maps/*
//go:embed ui/*
//go:embed tiles/*
var f embed.FS

var Images = map[string]*ebiten.Image{
	"box":          nil,
	"solid":        nil,
	"boulder":      nil,
	"gopher":       nil,
	"gopher-rip":   nil,
	"snake":        nil,
	"snake-snooze": nil,
	"plant":        nil,
	"time":         nil,
	"time-border":  nil,
}

var Runes = map[rune]string{
	'%': "box",
	'#': "solid",
	'B': "boulder",
	'1': "gopher",
	'2': "gopher",
	'3': "gopher",
	'4': "gopher",
	'r': "gopher-rip",
	's': "snake",
	'p': "plant",
}

func loadDefaultImages() error {
	for k := range Images {
		data, err := f.ReadFile(fmt.Sprintf("tiles/%s.png", k))
		if err != nil {
			return err
		}
		img, _, err := image.Decode(bytes.NewReader(data))
		if err != nil {
			return err
		}
		ebimg, err := ebiten.NewImageFromImage(img, ebiten.FilterNearest)
		if err != nil {
			return err
		}
		Images[k] = ebimg
	}
	return nil
}
