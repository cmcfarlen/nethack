package main

import (
	gc "github.com/rthornton128/goncurses"
	"math/rand"
	"time"

	"github.com/cmcfarlen/nethack/generate"
)

type Floor struct {
	Width   int
	Height  int
	terrain [][]gc.Char
}

func make2DArrayOfChar(w, h int) [][]gc.Char {
	a := make([][]gc.Char, w)
	for i := range a {
		a[i] = make([]gc.Char, h)
	}
	return a
}

func floorGen(w, h int) *Floor {
	terrain := make2DArrayOfChar(w, h)

	for x := 0; x < w; x += 1 {
		for y := 0; y < h; y += 1 {
			terrain[x][y] = gc.Char('.')
		}
	}

	start := Vector{rand.Intn(w-4) + 2, rand.Intn(h-4) + 2}

	terrain[start.x][start.y] = gc.ACS_DEGREE

	origin := add(start, Vector{-2, -2})
	for rooms := 0; rooms < 1; rooms += 1 {
		rw := rand.Intn(w / 5)
		rh := rand.Intn(h / 5)

		for x := origin.x; x < origin.x+rw; x += 1 {
			for y := origin.y; y < origin.y+rh; y += 1 {
				terrain[x][y] = gc.ACS_CKBOARD
			}
		}
	}

	/*
		walls := 0
		for walls < 10 {
			walls += 1
			v := Vector{rand.Intn(w), rand.Intn(h)}
			d := randomDirection()
			s := rand.Intn(10) + 5

			for i := 0; i < s; i += 1 {
				if v.x < 0 || v.y < 0 || v.x >= w || v.y >= h {
					break
				}
				terrain[v.x][v.y] = gc.ACS_CKBOARD
				v = add(v, d)
			}
		}
	*/

	return &Floor{w, h, terrain}
}

type World struct {
	currentFloor *Floor
}

func decAndClampToZero(x int) int {
	x -= 1
	if x < 0 {
		x = 0
	}
	return x
}

type Vector struct {
	x, y int
}

func add(a, b Vector) Vector {
	return Vector{a.x + b.x, a.y + b.y}
}

func clamp(x, from, to int) int {
	if x < from {
		x = from
	} else if x > to {
		x = to
	}
	return x
}

func clampZeroToV(v Vector, x, y int) Vector {
	return Vector{clamp(v.x, 0, x), clamp(v.y, 0, y)}
}

type Monster struct {
	Symbol gc.Char
	V      Vector
	P      Vector
}

func randomDirection() Vector {
	v := Vector{}
	for v.x == 0 && v.y == 0 {
		v.x = rand.Intn(3) - 1
		v.y = rand.Intn(3) - 1
	}
	return v
}

/*
  clamp a value
*/
func incAndClamp(x int, to int) int {
	x += 1
	if x > to {
		x = to
	}
	return x
}

func UpdateMonster(w *World, m *Monster) {
	m.V = randomDirection()
	m.P = clampZeroToV(add(m.P, m.V), w.currentFloor.Width, w.currentFloor.Height)
}

func main() {

	rand.Seed(time.Now().Unix())
	opts := generate.GenerateOpts{
		Width:             20,
		Height:            20,
		Sparseness:        70,
		DirectionModifier: 20,
		RoomCount:         4,
		RoomMin:           3,
		RoomMax:           7,
	}
	d := generate.GenerateDungeon(opts)

	d.PrintDungeon()
}

/*
func main() {
	// Initialize gc. It's essential End() is called to ensure the
	// terminal isn't altered after the program ends
	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal("init", err)
	}
	defer gc.End()

	rand.Seed(0) // same numbers please

	gc.Raw(true)
	gc.Echo(false)
	gc.Cursor(0)

	stdscr.Keypad(true)

	stdscr.Print("Press a key...")
	stdscr.Refresh()

	height, width := stdscr.MaxYX()
	height /= 2

	px := width / 2
	py := height / 2

	w := World{}
	m := Monster{gc.Char('M'), Vector{}, Vector{}}

	w.currentFloor = floorGen(width, height)

	keepGoing := true
	for keepGoing {

		stdscr.Erase()

		for x := 0; x < width; x = x + 1 {
			for y := 0; y < height; y = y + 1 {
				stdscr.MoveAddChar(y, x, w.currentFloor.terrain[x][y])
			}
		}

		stdscr.MoveAddChar(py, px, gc.ACS_DIAMOND)
		UpdateMonster(&w, &m)

		stdscr.MoveAddChar(m.P.y, m.P.x, m.Symbol)

		stdscr.MovePrintf(height, 0, "Fooo %v", m)

		stdscr.Refresh()

		ch := stdscr.GetChar()

		switch ch {
		case gc.KEY_ESC:
			keepGoing = false
		case 'k':
			py = decAndClampToZero(py)
		case 'j':
			py = incAndClamp(py, height)
		case 'h':
			px = decAndClampToZero(px)
		case 'l':
			px = incAndClamp(px, width)
		case 'g':
			w.currentFloor = floorGen(width, height)
		}
	}
}
*/
