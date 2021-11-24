package main

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"github.com/kettek/gophers-grievance/resources"
)

type Direction int

const (
	none Direction = iota
	north
	east
	south
	west
)

type GameState struct {
	game         *Game
	field        Field
	direction    Direction
	turnTime     time.Duration
	lastTurnTime time.Time
	dirs         map[Direction]struct{}
	turn         int
	difficulty   int
}

func (s *GameState) update(screen *ebiten.Image) error {
	// TODO: Separate direction out into a player type.
	if _, ok := s.dirs[west]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyA) && !ebiten.IsKeyPressed(ebiten.KeyH) {
			delete(s.dirs, west)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyH) {
			s.direction = west
			s.dirs[west] = struct{}{}
		}
	}
	if _, ok := s.dirs[east]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyD) && !ebiten.IsKeyPressed(ebiten.KeyL) {
			delete(s.dirs, east)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyL) {
			s.direction = east
			s.dirs[east] = struct{}{}
		}
	}
	if _, ok := s.dirs[north]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyW) && !ebiten.IsKeyPressed(ebiten.KeyK) {
			delete(s.dirs, north)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyK) {
			s.direction = north
			s.dirs[north] = struct{}{}
		}
	}
	if _, ok := s.dirs[south]; ok {
		if !ebiten.IsKeyPressed(ebiten.KeyS) && !ebiten.IsKeyPressed(ebiten.KeyJ) {
			delete(s.dirs, south)
		}
	} else {
		if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyJ) {
			s.direction = south
			s.dirs[south] = struct{}{}
		}
	}

	t := time.Now()
	if t.Sub(s.lastTurnTime) >= s.turnTime {
		s.simulate()
		s.turn++

		s.lastTurnTime = t
		s.direction = none
	}

	return nil
}

func (s *GameState) simulate() {
	if s.direction != none {
		if len(s.field.gophers) > 0 {
			s.field.moveObject(&s.field.gophers[0], s.direction)
		}
	}
	trapCount := 0
	for i, _ := range s.field.predators {
		p := &s.field.predators[i]
		if s.field.isTrapped(p) {
			if !p.trapped {
				p.trapped = true
				if p.t == snakeType {
					p.image = resources.SnakeSnoozeImage
				}
			}
			trapCount++
		} else {
			if p.trapped {
				p.trapped = false
				if p.t == snakeType {
					p.image = resources.SnakeImage
				}
			} else {
				if s.turn%(s.difficulty*5) == 1 {
					g := s.field.nearestGopher(p.x, p.y)
					if g != nil {
						s.field.moveTowards(p, g)
					}
					// Only move if the predator move timer is ready.
					// Get nearest gopher and its direction and begin moving towards it with a slight random x/y variance.
				}
			}
		}
	}
	// If all current predators are trapped, vegetize 'em.
	if trapCount == len(s.field.predators) {
		for _, p := range s.field.predators {
			s.field.tiles[p.y][p.x] = Tile{
				image: resources.FoodImage,
				food:  100,
			}
		}
		s.field.predators = make([]Object, 0)
	}
}

func (s *GameState) draw(screen *ebiten.Image) {
	screen.Fill(s.field.background)
	op := &ebiten.DrawImageOptions{}
	var offsetX float64 = 1
	var offsetY float64 = 1 + 332 - 276 // for now

	// Draw our scoreboard.
	for i, g := range s.field.gophers {
		score := fmt.Sprintf("Gopher %d - %d", i, g.score)
		text.Draw(screen, score, resources.BoldFont, 0, 10+i*10, color.White)
	}

	// Draw our borders.
	for y := 0; y < s.field.rows; y++ {
		op.GeoM.Reset()
		op.GeoM.Translate(offsetX, offsetY+float64(y)*tileHeight)
		screen.DrawImage(resources.SolidImage, op)

		op.GeoM.Reset()
		op.GeoM.Translate(offsetX+float64(s.field.columns-1)*tileWidth, offsetY+float64(y)*tileHeight)
		screen.DrawImage(resources.SolidImage, op)
	}
	for x := 1; x < s.field.columns-1; x++ {
		op.GeoM.Reset()
		op.GeoM.Translate(offsetX+float64(x)*tileWidth, offsetY)
		screen.DrawImage(resources.SolidImage, op)

		op.GeoM.Reset()
		op.GeoM.Translate(offsetX+float64(x)*tileWidth, offsetY+float64(s.field.rows-1)*tileHeight)
		screen.DrawImage(resources.SolidImage, op)
	}

	// Draw our map.
	for y, row := range s.field.tiles {
		for x, tile := range row {
			if tile.image == nil {
				continue
			}
			op.GeoM.Reset()
			op.GeoM.Translate(offsetX+float64(x)*tileWidth, offsetY+float64(y)*tileHeight)
			screen.DrawImage(tile.image, op)
		}
	}

	// Draw our gophers.
	for _, gopher := range s.field.gophers {
		op.GeoM.Reset()
		op.GeoM.Translate(offsetX+float64(gopher.x)*tileWidth, offsetY+float64(gopher.y)*tileHeight)
		screen.DrawImage(gopher.image, op)
	}

	// Draw our predators.
	for _, predator := range s.field.predators {
		op.GeoM.Reset()
		op.GeoM.Translate(offsetX+float64(predator.x)*tileWidth, offsetY+float64(predator.y)*tileHeight)
		screen.DrawImage(predator.image, op)
	}
}
