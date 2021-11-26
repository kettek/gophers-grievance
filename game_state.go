package main

import (
	"fmt"
	"image/color"
	"math/rand"
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

type FloatingText struct {
	remainingTime time.Duration
	text          string
	x, y          int
	color         color.RGBA
}

type GameState struct {
	game            *Game
	players         []Player
	field           Field
	turnTime        time.Duration
	lastTurnTime    time.Time
	currentTurnTime time.Time
	turn            int
	difficulty      int
	currentMap      resources.Map
	backgroundImage *ebiten.Image
	buttonAreaImage *ebiten.Image
	mapDone         bool
	mapExitTime     time.Time
	buttons         []Button
	mouseState      map[int]bool
	ui              *UiManager
	floatingText    []FloatingText
}

func (s *GameState) init() error {
	// TODO: Allow passing a seed.
	rand.Seed(time.Now().UnixNano())
	// Set our default settings.
	s.turnTime = 50 * time.Millisecond
	s.difficulty = 5
	s.players = []Player{
		{
			dirs:  make(map[Direction]struct{}),
			lives: maxLives,
		},
	}

	// Setup UI.
	backgroundImage, err := ebiten.NewImage(276, 276, ebiten.FilterLinear)
	if err != nil {
		return err
	}
	s.backgroundImage = backgroundImage

	buttonAreaImage, err := ebiten.NewImage(276, resources.ButtonLeftImage.Bounds().Dy()+2, ebiten.FilterLinear)
	if err != nil {
		return err
	}
	s.buttonAreaImage = buttonAreaImage
	s.buttonAreaImage.Fill(color.RGBA{
		180, 180, 180, 255,
	})

	s.mouseState = make(map[int]bool)
	s.buttons = []Button{
		{
			t: "Exit to Menu",
			cb: func() bool {
				s.game.setState(&MenuState{
					game: s.game,
					ui:   s.ui,
				})
				return true
			},
		},
		{
			t: "New Game",
			cb: func() bool {
				s.reset()
				return false
			},
		},
	}

	return nil
}

func (s *GameState) loadMap(m resources.Map) {
	s.currentMap = m
	s.mapDone = false
	s.turn = 0
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
	s.mapDone = false
	// Mark unused gophers
	for i := range s.field.gophers {
		if i >= len(s.players) {
			s.field.gophers[i].dead = true
		}
	}
	s.turn = 0
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
	s.currentTurnTime = s.lastTurnTime
	s.loadMap(resources.GetAnyMap())
}

func (s *GameState) update(screen *ebiten.Image) error {
	s.ui.checkButtons(s.buttons)

	// Otherwise, let's simulate.
	s.currentTurnTime = time.Now()
	d := s.currentTurnTime.Sub(s.lastTurnTime)

	// Update floating text.
	for i := len(s.floatingText) - 1; i >= 0; i-- {
		s.floatingText[i].remainingTime -= d
		if s.floatingText[i].remainingTime < 0 {
			s.floatingText = append(s.floatingText[:i], s.floatingText[i+1:]...)
		}
	}

	for i := range s.players {
		p := &s.players[i]
		p.update(s.lastTurnTime, s.currentTurnTime, d)
	}
	if d >= s.turnTime {
		s.simulate()
		s.turn++

		s.lastTurnTime = s.currentTurnTime

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
					s.floatingText = append(s.floatingText, FloatingText{
						text:          fmt.Sprintf("%d", v.score),
						x:             v.x * int(tileWidth),
						y:             v.y * int(tileHeight),
						color:         color.RGBA{gopherColor.R, gopherColor.G, gopherColor.B, 128},
						remainingTime: floatingTextDuration,
					})
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
	if s.mapDone {
		if s.currentTurnTime.After(s.mapExitTime) {
			s.loadMap(resources.GetNextMap(s.currentMap))
			s.mapDone = false
		}
	} else if trapCount == len(s.field.predators) { // If all current predators are trapped, vegetize 'em.
		for _, p := range s.field.predators {
			s.field.tiles[p.y][p.x] = Tile{
				image: resources.FoodImage,
				food:  100,
			}
			// Reward 500 to each alive gopher for each vegitization.
			for i, g := range s.field.gophers {
				if !g.dead {
					s.players[i].score += 500
					s.floatingText = append(s.floatingText, FloatingText{
						text:          fmt.Sprintf("%d", 500),
						x:             p.x * int(tileWidth),
						y:             p.y * int(tileHeight),
						color:         color.RGBA{gopherColor.R, gopherColor.G, gopherColor.B, 128},
						remainingTime: floatingTextDuration,
					})
				}
			}
		}
		s.field.predators = make([]Object, 0)
		s.mapDone = true
		s.mapExitTime = s.currentTurnTime.Add(5 * time.Second)
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
	screen.Fill(color.RGBA{0, 125, 156, 255})
	// Draw our clock.
	clockX := float64(winWidth)/2 - float64(resources.TimeImage.Bounds().Max.X)/2
	clockY := float64(332 - 276 - resources.TimeImage.Bounds().Max.Y - 1)
	op.GeoM.Translate(-float64(resources.TimeImage.Bounds().Max.X)/2, -float64(resources.TimeImage.Bounds().Max.Y)/2)
	op.GeoM.Rotate(float64(s.turn) / 2000.0 * 6.28)
	op.GeoM.Translate(float64(resources.TimeImage.Bounds().Max.X)/2, float64(resources.TimeImage.Bounds().Max.Y)/2)
	op.GeoM.Translate(clockX, clockY)
	screen.DrawImage(resources.TimeImage, op)
	op.GeoM.Reset()
	op.GeoM.Translate(clockX, clockY)
	screen.DrawImage(resources.TimeBorderImage, op)

	op.GeoM.Reset()
	screen.DrawImage(s.buttonAreaImage, op)
	btnX := 1
	for i := range s.buttons {
		btn := &s.buttons[i]
		btn.draw(screen, btnX, 1)
		btnX += btn.w + 2
	}

	topOffset := resources.ButtonMiddleImage.Bounds().Dy() + 2

	// Draw our next map timer
	if s.mapDone {
		op.GeoM.Reset()
		op.GeoM.Translate(clockX, clockY)
		timeLeft := fmt.Sprintf("%d", int(s.mapExitTime.Sub(s.currentTurnTime).Seconds()))
		r := text.BoundString(resources.BoldFont, timeLeft)
		text.Draw(screen, timeLeft, resources.BoldFont, int(clockX+float64(resources.TimeImage.Bounds().Max.X/2)-float64(r.Max.X/2)), int(clockY+float64(resources.TimeImage.Bounds().Max.Y/2)+float64(r.Max.Y/2)), color.RGBA{255, 0, 255, 255})
	}

	// Draw our scoreboard.
	for i, p := range s.players {
		var offsetX float64 = 0
		var offsetY float64 = float64(topOffset) + 1 + float64(i)*10
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
		text.Draw(screen, score, resources.BoldFont, int(float64(maxLives)*tileWidth), topOffset+12+i*12, color.White)
	}

	var offsetX float64 = 0
	var offsetY float64 = 332 - 276 // for now

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

	// Draw our silly score text.
	for _, t := range s.floatingText {
		r := text.BoundString(resources.NormalFont, t.text)
		centerX := r.Dx() / 2
		c := t.color
		p := float64(t.remainingTime) / float64(floatingTextDuration)
		c.A = uint8(p * float64(floatingTextAlpha))
		text.Draw(screen, t.text, resources.NormalFont, int(offsetX)+t.x-centerX, int(offsetY)-4+t.y-int(1-p*6), c)
	}
}
