package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
)

type Field struct {
	name          string
	columns, rows int

	tiles     [][]Tile
	predators []Object
	gophers   []Object
}

type Tile struct {
	image    *ebiten.Image
	blocking bool
	pushable bool
	food     int
}

type ObjectType int

const (
	noType ObjectType = iota
	gopherType
	snakeType
)

type Object struct {
	y, x            int
	image           *ebiten.Image
	ripImage        *ebiten.Image
	t               ObjectType
	trapped         bool
	dead            bool
	failedMovements int
}

func (f *Field) fromMap(m resources.Map, clearRip bool) {
	// Reset
	f.predators = make([]Object, 0)
	f.gophers = make([]Object, 0)

	// Store our dead gophers. FIXME: we're being lazy and temporarily storing them as objects.
	deadGophers := make([]Object, 0)
	if !clearRip {
		for y, row := range f.tiles {
			for x, tile := range row {
				if tile.image == resources.Images["gopher-rip"] {
					deadGophers = append(deadGophers, Object{
						y: y,
						x: x,
					})
				}
			}
		}
	}

	// Size our map
	f.tiles = make([][]Tile, m.Rows)
	for i := 0; i < m.Rows; i++ {
		f.tiles[i] = make([]Tile, m.Columns)
	}
	f.columns = m.Columns
	f.rows = m.Rows
	f.name = m.Name

	if !clearRip {
		for _, o := range deadGophers {
			if f.tiles[o.y][o.x].image == nil {
				f.tiles[o.y][o.x].image = resources.Images["gopher-rip"]
			}
		}
	}

	// Convert our map resource to a live map.
	for y, row := range m.Cells {
		for x, r := range row {
			switch r {
			case '%':
				f.tiles[y][x] = Tile{
					image:    resources.Images["box"],
					pushable: true,
					blocking: true,
				}
			case 'B':
				f.tiles[y][x] = Tile{
					image:    resources.Images["boulder"],
					pushable: true,
					blocking: true,
				}
			case '#':
				f.tiles[y][x] = Tile{
					image:    resources.Images["solid"],
					pushable: false,
					blocking: true,
				}
			case 's':
				f.predators = append(f.predators, Object{
					image: resources.Images["snake"],
					y:     y,
					x:     x,
					t:     snakeType,
				})
			case 'p':
				f.tiles[y][x] = Tile{
					image:    resources.Images["plant"],
					pushable: false,
					blocking: false,
					food:     100,
				}
			case '1', '2', '3', '4':
				// TODO: Limit based on current player count.
				f.gophers = append(f.gophers, Object{
					image:    resources.Images["gopher"],
					ripImage: resources.Images["gopher-rip"],
					y:        y,
					x:        x,
					t:        gopherType,
				})
			}
		}
	}
	// If we have no predators, randomly add one or more to an open space.
	if len(f.predators) == 0 {
		min := 1
		max := 2
		r := rand.Intn(max-min+1) + min
		for i := 0; i < r; i++ {
			f.spawnRandomPredator(snakeType)
		}
	}
}

func (f *Field) spawnRandomPredator(ot ObjectType) {
	freeTiles := f.findOpenTiles()
	if len(freeTiles) == 0 {
		return
	}
	t := rand.Intn(len(freeTiles) - 1 + 1)
	f.predators = append(f.predators, Object{
		image: resources.Images["snake"],
		y:     freeTiles[t].y,
		x:     freeTiles[t].x,
		t:     ot,
	})
}

type Coord struct {
	x, y int
}

func (f *Field) findOpenTiles() (c []Coord) {
	for y := 0; y < f.rows; y++ {
		for x := 0; x < f.columns; x++ {
			found := false
			for _, g := range f.gophers {
				if g.x == x && g.y == y {
					found = true
					break
				}
			}
			for _, p := range f.predators {
				if p.x == x && p.y == y {
					found = true
					break
				}
			}

			if f.tiles[y][x].image == nil && !found {
				c = append(c, Coord{x, y})
			}
		}
	}
	return c
}

func (f *Field) moveObject(o *Object, dir Direction) moveResult {
	x := 0
	y := 0
	switch dir {
	case west:
		x = -1
	case east:
		x = 1
	case south:
		y = 1
	case north:
		y = -1
	}
	tx := o.x + x
	ty := o.y + y
	if !f.inBounds(tx, ty) {
		return moveBlockedResult{
			x: tx,
			y: ty,
		}
	}
	if f.isBlocked(tx, ty) {
		if f.isPushable(tx, ty) {
			if f.push(tx, ty, x, y) {
				o.x = tx
				o.y = ty
				return movePushResult{
					x: tx,
					y: ty,
				}
			}
		}
		return moveBlockedResult{
			x: tx,
			y: ty,
		}
	}

	if o.t == gopherType {
		if f.tiles[ty][tx].food > 0 {
			o.x = tx
			o.y = ty
			food := f.tiles[ty][tx].food
			f.tiles[ty][tx] = Tile{}
			return moveEatResult{
				score: food,
				x:     tx,
				y:     ty,
			}
		}
	}

	if f.isEmpty(tx, ty) {
		o.x = tx
		o.y = ty
		return moveSuccessResult{
			x: tx,
			y: ty,
		}
	}
	return moveBlockedResult{
		x: tx,
		y: ty,
	}
}

func (f *Field) inBounds(x, y int) bool {
	if x < 0 || x >= f.columns {
		return false
	}
	if y < 0 || y >= f.rows {
		return false
	}
	return true
}

func (f *Field) isEmpty(x, y int) bool {
	for _, p := range f.predators {
		if p.x == x && p.y == y {
			return false
		}
	}
	for _, g := range f.gophers {
		if g.dead {
			continue
		}
		if g.x == x && g.y == y {
			return false
		}
	}
	return !f.tiles[y][x].blocking && !f.tiles[y][x].pushable
}

func (f *Field) isBlocked(x, y int) bool {
	return f.tiles[y][x].blocking
}

func (f *Field) isPushable(x, y int) bool {
	return f.tiles[y][x].pushable
}

func (f *Field) isGopherAt(x, y int) bool {
	for _, o := range f.gophers {
		if o.dead {
			continue
		}
		if o.x != x || o.y != y {
			continue
		}
		return true
	}
	return false
}

func (f *Field) getGopherAt(x, y int) int {
	for i, o := range f.gophers {
		if o.dead {
			continue
		}
		if o.x != x || o.y != y {
			continue
		}
		return i
	}
	return -1
}

func (f *Field) push(x, y int, xDir, yDir int) bool {
	swap := func(x1, y1, x2, y2 int) {
		a := f.tiles[y1][x1]
		b := f.tiles[y2][x2]
		f.tiles[y1][x1] = b
		f.tiles[y2][x2] = a
	}

	if xDir != 0 {
		if f.isEmpty(x+xDir, y) {
			swap(x, y, x+xDir, y)
			return true
		} else if f.isPushable(x+xDir, y) {
			if f.push(x+xDir, y, xDir, yDir) {
				swap(x, y, x+xDir, y)
				return true
			}
		}
	} else if yDir != 0 {
		if f.isEmpty(x, y+yDir) {
			swap(x, y, x, y+yDir)
			return true
		} else if f.isPushable(x, y+yDir) {
			if f.push(x, y+yDir, xDir, yDir) {
				swap(x, y, x, y+yDir)
				return true
			}
		}
	}
	return false
}

func (f *Field) isTrapped(o *Object) bool {
	if f.inBounds(o.x-1, o.y) && !f.isBlocked(o.x-1, o.y) && f.isEmpty(o.x-1, o.y) {
		return false
	}
	if f.inBounds(o.x+1, o.y) && !f.isBlocked(o.x+1, o.y) && f.isEmpty(o.x+1, o.y) {
		return false
	}
	if f.inBounds(o.x, o.y-1) && !f.isBlocked(o.x, o.y-1) && f.isEmpty(o.x, o.y-1) {
		return false
	}
	if f.inBounds(o.x, o.y+1) && !f.isBlocked(o.x, o.y+1) && f.isEmpty(o.x, o.y+1) {
		return false
	}
	// diagonals
	if f.inBounds(o.x-1, o.y-1) && !f.isBlocked(o.x-1, o.y-1) && f.isEmpty(o.x-1, o.y-1) {
		return false
	}
	if f.inBounds(o.x+1, o.y-1) && !f.isBlocked(o.x+1, o.y-1) && f.isEmpty(o.x+1, o.y-1) {
		return false
	}
	if f.inBounds(o.x-1, o.y+1) && !f.isBlocked(o.x-1, o.y+1) && f.isEmpty(o.x-1, o.y+1) {
		return false
	}
	if f.inBounds(o.x+1, o.y+1) && !f.isBlocked(o.x+1, o.y+1) && f.isEmpty(o.x+1, o.y+1) {
		return false
	}

	return true
}

func (f *Field) nearestGopher(x, y int) (o *Object) {
	for i, g := range f.gophers {
		if g.dead {
			continue
		}
		if o == nil {
			o = &f.gophers[i]
		} else {
			if math.Sqrt(math.Pow(math.Abs(float64(g.x-x)), 2)+math.Pow(math.Abs(float64(g.y-y)), 2)) < math.Sqrt(math.Pow(math.Abs(float64(o.x-x)), 2)+math.Pow(math.Abs(float64(o.y-y)), 2)) {
				o = &f.gophers[i]
			}
		}
	}
	return o
}

func (f *Field) moveTowards(o *Object, t *Object, turn int) moveResult {
	dirX := 0
	dirY := 0
	if t.x < o.x {
		dirX = -1
	} else if t.x > o.x {
		dirX = 1
	}
	if t.y < o.y {
		dirY = -1
	} else if t.y > o.y {
		dirY = 1
	}

	if o.failedMovements > 0 {
		min := 0
		max := 1
		r := rand.Intn(max-min+1) + min

		if r == 0 {
			dirX = 1
		} else if r == 1 {
			dirX = -1
		}
		r = rand.Intn(max-min+1) + min
		if r == 0 {
			dirY = -1
		} else if r == 1 {
			dirY = 1
		}
	}

	tx := o.x + dirX
	ty := o.y + dirY
	if !f.inBounds(tx, ty) || f.isBlocked(tx, ty) || !f.isEmpty(tx, ty) {
		if f.isGopherAt(tx, ty) {
			return moveTouchResult{
				gopher: f.getGopherAt(tx, ty),
				x:      tx,
				y:      ty,
			}
		}
		tx = o.x
		if !f.inBounds(tx, ty) || f.isBlocked(tx, ty) || !f.isEmpty(tx, ty) {
			if f.isGopherAt(tx, ty) {
				return moveTouchResult{
					gopher: f.getGopherAt(tx, ty),
					x:      tx,
					y:      ty,
				}
			}
			tx = o.x + dirX
			ty = o.y
			if !f.inBounds(tx, ty) || f.isBlocked(tx, ty) || !f.isEmpty(tx, ty) {
				if f.isGopherAt(tx, ty) {
					return moveTouchResult{
						gopher: f.getGopherAt(tx, ty),
						x:      tx,
						y:      ty,
					}
				}
				o.failedMovements++
				return moveBlockedResult{}
			}
		}
	}
	o.x = tx
	o.y = ty
	o.failedMovements = 0
	return moveSuccessResult{
		x: tx,
		y: ty,
	}
}

func (f *Field) setTile(x, y int, t Tile) {
	if !f.inBounds(x, y) {
		return
	}
	f.tiles[y][x] = t
}
