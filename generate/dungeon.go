package generate

// https://dirkkok.wordpress.com/dungeon-generation-article-series/

import (
	"fmt"
	"math"
	"math/rand"
)

type (
	tileType      int
	sideType      int
	moveDirection int

	tile struct {
		tiletype                 tileType
		x, y                     int
		north, south, east, west sideType
	}

	// Map is a map
	Map struct {
		Width, Height int
		data          []rune
	}

	// Room is for holding things
	Room struct {
		X, Y          int
		Width, Height int
	}

	// Dungeon is a dungeon with tiles and rooms and a map
	Dungeon struct {
		Width, Height int
		tiles         []tile
		Rooms         []Room
		M             Map
	}

	// DungeonGenerator generates dungeons
	DungeonGenerator struct {
		d            *Dungeon
		currentTile  *tile
		visitedTiles []tile
	}

	// Opts are options for generating a dungeon
	Opts struct {
		Width, Height     int
		Sparseness        int
		DirectionModifier int
		RoomCount         int
		RoomMin           int
		RoomMax           int
	}

	directionPicker struct {
		directionsRemaining []moveDirection
	}
)

const (
	tileTypeUnvisited = iota
	tileTypeCorridor
	tileTypeRoom
	tileTypeWall
	tileTypeEmpty
)

const (
	sideTypeWall = iota
	sideTypeEmpty
	siteTypeDoor
)

const (
	moveNorth = iota
	moveSouth
	moveWest
	moveEast
)

// RuneAt reeturns the rune at the point x, y
func (m *Map) RuneAt(x, y int) rune {
	if x >= 0 && x < m.Width && y >= 0 && y < m.Height {
		return m.data[x*m.Width+y]
	}
	return 0
}

// IsWalkable returns true if x, y is walkable
func (m *Map) IsWalkable(x, y int) bool {
	r := m.RuneAt(x, y)
	return r == ' ' || r == '.'
}

func newDirectionPicker() *directionPicker {
	dirs := [4]moveDirection{}

	for i, v := range rand.Perm(4) {
		dirs[i] = moveDirection(v)
	}

	return &directionPicker{dirs[:]}
}

func (p *directionPicker) hasNextDirection() bool {
	return len(p.directionsRemaining) > 0
}

func (p *directionPicker) nextDirection() moveDirection {
	nd := p.directionsRemaining[0]
	p.directionsRemaining = p.directionsRemaining[1:]
	return nd
}

func newDungeon(w, h, roomCount int) *Dungeon {
	d := Dungeon{w, h, make([]tile, w*h), make([]Room, roomCount), Map{}}

	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			t := d.tileAt(x, y)
			t.x = x
			t.y = y
		}
	}

	return &d
}

func (t *tile) visited() bool {
	return t.tiletype != tileTypeUnvisited
}

func (t *tile) isDeadEnd() bool {
	walls := 0
	if t.north == sideTypeWall {
		walls++
	}
	if t.south == sideTypeWall {
		walls++
	}
	if t.east == sideTypeWall {
		walls++
	}
	if t.west == sideTypeWall {
		walls++
	}
	return walls == 3
}

func (t *tile) deadEndDirection() moveDirection {
	if t.north == sideTypeEmpty {
		return moveNorth
	}
	if t.south == sideTypeEmpty {
		return moveSouth
	}
	if t.east == sideTypeEmpty {
		return moveEast
	}
	if t.west == sideTypeEmpty {
		return moveWest
	}
	return 0
}

func (d *Dungeon) visitAllTiles(f func(*tile)) {
	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			f(d.tileAt(x, y))
		}
	}
}

func (d *Dungeon) randomXY() (x, y int) {
	x = rand.Intn(d.Width)
	y = rand.Intn(d.Height)
	return
}

func (d *Dungeon) tileAt(x, y int) *tile {
	if x >= 0 && x < d.Width && y >= 0 && y < d.Height {
		return &d.tiles[y*d.Height+x]
	}
	return nil
}

func (d *Dungeon) randomTile() *tile {
	x, y := d.randomXY()
	t := d.tileAt(x, y)
	t.tiletype = tileTypeCorridor
	return t
}

func (d *Dungeon) tilesOfType(tt tileType) []*tile {
	r := make([]*tile, 0)
	for _, t := range d.tiles {
		if t.tiletype == tt {
			r = append(r, &t)
		}
	}

	return r
}

func (d *Dungeon) moveTile(t *tile, md moveDirection) *tile {
	n := d.tileInDirection(t, md)
	if n != nil && n.tiletype == tileTypeUnvisited {
		return n
	}

	return nil
}

func (d *Dungeon) tileInDirection(t *tile, md moveDirection) *tile {
	x, y := move(t.x, t.y, md)

	if x >= 0 && x < d.Width && y >= 0 && y < d.Height {
		return d.tileAt(x, y)
	}
	return nil
}

func randomDirection() moveDirection {
	return moveDirection(rand.Intn(4))
}

func move(x, y int, dir moveDirection) (int, int) {
	switch dir {
	case moveNorth:
		return x, y - 1
	case moveSouth:
		return x, y + 1
	case moveWest:
		return x - 1, y
	case moveEast:
		return x + 1, y
	}
	return x, y
}

func pickValidDirection(d *Dungeon, t *tile) (moveDirection, bool) {
	for _, v := range rand.Perm(4) {
		md := moveDirection(v)
		if d.moveTile(t, md) != nil {
			return md, true
		}
	}

	return moveNorth, false
}

func (d *Dungeon) updateMap() {
	mw := d.Width*2 + 1
	mh := d.Height*2 + 1
	m := make([]rune, mw*mh)
	ttom := func(x, y int) (int, int) {
		return x*2 + 1, y*2 + 1
	}
	setm := func(x, y int, c rune) {
		m[x*mw+y] = c
	}

	for i := 0; i < len(m); i++ {
		m[i] = '#'
	}

	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			t := d.tileAt(x, y)
			mx, my := ttom(x, y)

			if t.tiletype == tileTypeEmpty {
				continue
			}

			floorRune := ' '
			if t.tiletype == tileTypeCorridor {
				floorRune = ' '
			}
			if t.tiletype == tileTypeRoom {
				floorRune = '.'
			}

			setm(mx, my, floorRune)

			if t.west == sideTypeWall {
				setm(mx-1, my, '|')
				setm(mx-1, my-1, '|')
				setm(mx-1, my+1, '|')
			} else {
				setm(mx-1, my, floorRune)
			}

			if t.east == sideTypeWall {
				setm(mx+1, my, '|')
				setm(mx+1, my-1, '|')
				setm(mx+1, my+1, '|')
			} else {
				setm(mx+1, my, floorRune)
			}

			if t.north == sideTypeWall {
				setm(mx, my-1, '-')
				setm(mx-1, my-1, '-')
				setm(mx+1, my-1, '-')
			} else {
				setm(mx, my-1, floorRune)
			}

			if t.south == sideTypeWall {
				setm(mx, my+1, '-')
				setm(mx+1, my+1, '-')
				setm(mx-1, my+1, '-')
			} else {
				setm(mx, my+1, floorRune)
			}

			if t.south == sideTypeEmpty &&
				t.east == sideTypeEmpty {
				setm(mx+1, my+1, floorRune)
			}

		}
	}

	d.M.data = m
	d.M.Width = mw
	d.M.Height = mh
}

func (m *Map) PrintMap() {
	for x := 0; x < m.Width; x++ {
		for y := 0; y < m.Height; y++ {
			fmt.Print(string(m.data[y*m.Width+x]))
		}
		fmt.Println("")
	}
}

func (d *Dungeon) setSideType(from *tile, to *tile, dir moveDirection, tt sideType) {
	switch dir {
	case moveNorth:
		from.north = tt
		to.south = tt
	case moveSouth:
		from.south = tt
		to.north = tt
	case moveWest:
		from.west = tt
		to.east = tt
	case moveEast:
		from.east = tt
		to.west = tt
	}
}

func (d *Dungeon) createCorridor(from *tile, to *tile, dir moveDirection) {
	from.tiletype = tileTypeCorridor
	to.tiletype = tileTypeCorridor

	d.setSideType(from, to, dir, sideTypeEmpty)
}

func (d *Dungeon) deadEndTiles() []*tile {
	det := make([]*tile, 0)

	for x := 0; x < d.Width; x++ {
		for y := 0; y < d.Height; y++ {
			t := d.tileAt(x, y)

			if t.isDeadEnd() {
				det = append(det, t)
			}
		}
	}
	return det
}

func randomBetween(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func (d *Dungeon) adjacentTiles(t *tile) []*tile {
	adj := make([]*tile, 0)

	if a := d.tileAt(t.x-1, t.y-1); a != nil {
		adj = append(adj, a)
	}
	if a := d.tileAt(t.x-1, t.y); a != nil {
		adj = append(adj, a)
	}
	if a := d.tileAt(t.x-1, t.y+1); a != nil {
		adj = append(adj, a)
	}
	if a := d.tileAt(t.x+1, t.y-1); a != nil {
		adj = append(adj, a)
	}
	if a := d.tileAt(t.x+1, t.y); a != nil {
		adj = append(adj, a)
	}
	if a := d.tileAt(t.x+1, t.y+1); a != nil {
		adj = append(adj, a)
	}
	if a := d.tileAt(t.x, t.y-1); a != nil {
		adj = append(adj, a)
	}
	if a := d.tileAt(t.x, t.y+1); a != nil {
		adj = append(adj, a)
	}

	return adj
}

func (d *Dungeon) scoreRoom(startx, starty int, rw, rh int) int {

	if (startx+rw) >= d.Width ||
		(starty+rh) >= d.Height {
		return math.MaxInt32
	}

	if t := d.tileAt(startx, starty); t.tiletype != tileTypeCorridor {
		return math.MaxInt32
	}

	score := 0

	for x := startx; x < (startx + rw); x++ {
		for y := starty; y < (starty + rh); y++ {
			t := d.tileAt(x, y)

			if t.tiletype == tileTypeCorridor {
				score += 3
			}

			if t.tiletype == tileTypeRoom {
				score += 100
			}

			for _, a := range d.adjacentTiles(t) {
				if a.tiletype == tileTypeCorridor {
					score += 1
				}
			}
		}
	}

	return score
}

func GenerateDungeon(opts Opts) *Dungeon {
	d := newDungeon(opts.Width, opts.Height, opts.RoomCount)
	tileCount := opts.Width * opts.Height
	visited := make(map[*tile]bool)

	current := d.randomTile()
	visited[current] = true
	dp := newDirectionPicker()
	dir := dp.nextDirection()

	for len(visited) < tileCount {
		if rand.Intn(100) < opts.DirectionModifier {
			if dp.hasNextDirection() {
				dir = dp.nextDirection()
			} else {
				dp = newDirectionPicker()
				dir = dp.nextDirection()
			}
		}

		for d.moveTile(current, dir) == nil {
			if dp.hasNextDirection() {
				dir = dp.nextDirection()
			} else {
				dp = newDirectionPicker()
				dir = dp.nextDirection()

				keys := make([]*tile, len(visited))
				i := 0
				for k := range visited {
					keys[i] = k
					i++
				}
				current = keys[rand.Intn(len(visited))]
			}
		}

		next := d.moveTile(current, dir)
		d.createCorridor(current, next, dir)
		current = next
		visited[current] = true
	}

	removeCount := int(math.Ceil(float64(d.Width) * float64(d.Height) * float64(opts.Sparseness) / 100.0))

	for removeCount > 0 {
		det := d.deadEndTiles()
		if len(det) > removeCount {
			det = det[:removeCount]
		}

		for _, t := range det {
			dir := t.deadEndDirection()
			t.tiletype = tileTypeEmpty

			ot := d.tileInDirection(t, dir)

			d.setSideType(t, ot, dir, sideTypeWall)
		}
		removeCount -= len(det)
	}

	for r := 0; r < opts.RoomCount; r++ {
		roomWidth := randomBetween(opts.RoomMin, opts.RoomMax)
		roomHeight := randomBetween(opts.RoomMin, opts.RoomMax)

		bestScore := math.MaxInt32
		bestx := 0
		besty := 0

		for x := 0; x < d.Width; x++ {
			for y := 0; y < d.Height; y++ {
				score := d.scoreRoom(x, y, roomWidth, roomHeight)
				if score < bestScore {
					bestScore = score
					bestx = x
					besty = y
				}
			}
		}

		if bestScore < math.MaxInt32 {
			for x := bestx; x < bestx+roomWidth; x++ {
				for y := besty; y < besty+roomHeight; y++ {
					t := d.tileAt(x, y)

					t.tiletype = tileTypeRoom
					if x != bestx {
						t.west = sideTypeEmpty
					}
					if y != besty {
						t.north = sideTypeEmpty
					}
					if x < bestx+roomWidth-1 {
						t.east = sideTypeEmpty
					}
					if y < besty+roomHeight-1 {
						t.south = sideTypeEmpty
					}
				}
			}
		}

		d.Rooms[r].X = bestx
		d.Rooms[r].Y = besty
		d.Rooms[r].Width = roomWidth
		d.Rooms[r].Height = roomHeight
	}

	d.updateMap()

	return d
}
