package main

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
	"github.com/kettek/gophers-grievance/resources"
)

type Field struct {
	columns, rows int
	background    color.RGBA

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
	y, x    int
	image   *ebiten.Image
	t       ObjectType
	trapped bool
	score   int // TODO: separate Object out into an interface.
}

func (f *Field) fromMap(m resources.Map) {
	// Size our map
	f.tiles = make([][]Tile, m.Rows)
	for i := 0; i < m.Rows; i++ {
		f.tiles[i] = make([]Tile, m.Columns)
	}
	f.columns = m.Columns
	f.rows = m.Rows

	// Convert our map resource to a live map.
	for y, row := range m.Cells {
		for x, r := range row {
			switch r {
			case 'B':
				f.tiles[y][x] = Tile{
					image:    resources.BoulderImage,
					pushable: true,
					blocking: true,
				}
			case '#':
				f.tiles[y][x] = Tile{
					image:    resources.SolidImage,
					pushable: false,
					blocking: true,
				}
			case 's':
				f.predators = append(f.predators, Object{
					image: resources.SnakeImage,
					y:     y,
					x:     x,
					t:     snakeType,
				})
			case 'p':
				f.tiles[y][x] = Tile{
					image:    resources.FoodImage,
					pushable: false,
					blocking: false,
					food:     100,
				}
			case '1', '2', '3', '4':
				// TODO: Limit based on current player count.
				f.gophers = append(f.gophers, Object{
					image: resources.GopherImage,
					y:     y,
					x:     x,
					t:     gopherType,
				})
			}
		}
	}
}

func (f *Field) moveObject(o *Object, dir Direction) {
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
		return
	}
	if f.isBlocked(tx, ty) {
		if f.isPushable(tx, ty) {
			if f.push(tx, ty, x, y) {
				o.x = tx
				o.y = ty
			}
		}
		return
	}

	if o.t == gopherType {
		if f.tiles[ty][tx].food > 0 {
			o.score += f.tiles[ty][tx].food
			f.tiles[ty][tx] = Tile{}
		}
	}

	// Check if a predator is there.
	/*var predator *Object = nil
	for _, p := range f.predators {
		if p.x == tx && p.y == ty {
			predator = &p
			break
		}
	}
	if predator != nil {
	}*/

	o.x = tx
	o.y = ty
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
	if f.inBounds(o.x-1, o.y) && !f.isBlocked(o.x-1, o.y) {
		return false
	}
	if f.inBounds(o.x+1, o.y) && !f.isBlocked(o.x+1, o.y) {
		return false
	}
	if f.inBounds(o.x, o.y-1) && !f.isBlocked(o.x, o.y-1) {
		return false
	}
	if f.inBounds(o.x, o.y+1) && !f.isBlocked(o.x, o.y+1) {
		return false
	}
	// diagonals
	if f.inBounds(o.x-1, o.y-1) && !f.isBlocked(o.x-1, o.y-1) {
		return false
	}
	if f.inBounds(o.x+1, o.y-1) && !f.isBlocked(o.x+1, o.y-1) {
		return false
	}
	if f.inBounds(o.x-1, o.y+1) && !f.isBlocked(o.x-1, o.y+1) {
		return false
	}
	if f.inBounds(o.x+1, o.y+1) && !f.isBlocked(o.x+1, o.y+1) {
		return false
	}

	return true
}

func (f *Field) nearestGopher(x, y int) (o *Object) {
	for i, g := range f.gophers {
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

func (f *Field) moveTowards(o *Object, t *Object) {
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
	tx := o.x + dirX
	ty := o.y + dirY
	if !f.inBounds(tx, ty) {
		return
	}
	if f.isBlocked(tx, ty) || !f.isEmpty(tx, ty) {
		return
	}
	o.x = tx
	o.y = ty
}
