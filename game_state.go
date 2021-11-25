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
	game            *Game
	players         []Player
	field           Field
	turnTime        time.Duration
	lastTurnTime    time.Time
	turn            int
	difficulty      int
	currentMap      resources.Map
	backgroundImage *ebiten.Image
}

func (s *GameState) loadMap(m resources.Map) {
	s.currentMap = m
	s.field.fromMap(s.currentMap, true)
	// Mark unused gophers
	for i := range s.field.gophers {
		if i >= len(s.players) {
			s.field.gophers[i].dead = true
		}
	}
	s.backgroundImage.Fill(s.currentMap.Background)
	ebiten.SetWindowTitle(fmt.Sprintf("%s: %s", winTitle, m.Name))
}

func (s *GameState) resetMap() {
	s.field.fromMap(s.currentMap, false)
	// Mark unused gophers
	for i := range s.field.gophers {
		if i >= len(s.players) {
			s.field.gophers[i].dead = true
		}
	}
}

// reset resets the game state to a fresh one.
func (s *GameState) reset() {
	for i := range s.players {
		s.players[i] = Player{
			lives: maxLives,
			dirs:  make(map[Direction]struct{}),
		}
	}
	s.lastTurnTime = time.Now()
	s.loadMap(s.currentMap)
}

func (s *GameState) update(screen *ebiten.Image) error {
	t := time.Now()
	d := t.Sub(s.lastTurnTime)
	for i := range s.players {
		p := &s.players[i]
		p.update(s.lastTurnTime, t, d)
	}
	if d >= s.turnTime {
		s.simulate()
		s.turn++

		s.lastTurnTime = t

		for i := range s.players {
			p := &s.players[i]
			p.direction = none
		}
	}

	return nil
}

func (s *GameState) simulate() {
	for i := range s.players {
		p := &s.players[i]
		if p.direction != none {
			if i < len(s.field.gophers) {
				if s.field.gophers[i].dead {
					continue
				}
				r := s.field.moveObject(&s.field.gophers[i], p.direction)
				switch v := r.(type) {
				case moveEatResult:
					p.score += v.score
				}
			}
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
						r := s.field.moveTowards(p, g, s.turn)
						switch v := r.(type) {
						case moveTouchResult:
							s.players[v.gopher].reduceLives()
							s.field.gophers[v.gopher].dead = true
							s.field.setTile(g.x, g.y, Tile{
								image: resources.GopherRipImage,
							})
						}
					}
				}
			}
		}
	}
	deathCount := 0
	for _, g := range s.field.gophers {
		if g.dead {
			deathCount++
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
		// TODO: Wait for a game end timer to allow gophers to pick up vegetized predators before traveling to next map (or, if the timer ticks, spawn more predators)
		s.loadMap(resources.GetNextMap(s.currentMap))
	} else if deathCount == len(s.field.gophers) { // Prioritize predator death over gopher death.
		playersOut := 0
		for _, p := range s.players {
			if p.lives < 0 {
				playersOut++
			}
		}
		if playersOut == len(s.players) {
			// TODO: Game over
			s.reset()
		} else {
			// TODO: Pause until a player signals ready
			s.resetMap()
		}
	}
}

func (s *GameState) draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	var offsetX float64 = 1
	var offsetY float64 = 1 + 332 - 276 // for now

	// Draw our scoreboard.
	for i, p := range s.players {
		var offsetX float64 = 0
		var offsetY float64 = 1 + float64(i)*10
		for l := 0; l < maxLives; l++ {
			op.GeoM.Reset()
			op.GeoM.Translate(offsetX+float64(l)*tileWidth, offsetY+float64(i)*tileHeight)
			if l >= p.lives {
				screen.DrawImage(resources.GopherRipImage, op)
			} else {
				screen.DrawImage(resources.GopherImage, op)
			}
		}

		score := fmt.Sprintf("Gopher %d - %d", i, p.score)
		text.Draw(screen, score, resources.BoldFont, int(float64(maxLives)*tileWidth), 10+i*10, color.White)
	}

	op.GeoM.Reset()
	op.GeoM.Translate(0, offsetY)
	screen.DrawImage(s.backgroundImage, op)

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
		if !gopher.dead {
			screen.DrawImage(gopher.image, op)
		}
	}

	// Draw our predators.
	for _, predator := range s.field.predators {
		op.GeoM.Reset()
		op.GeoM.Translate(offsetX+float64(predator.x)*tileWidth, offsetY+float64(predator.y)*tileHeight)
		screen.DrawImage(predator.image, op)
	}
}
