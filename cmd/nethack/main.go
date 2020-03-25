package main

import (
	gc "github.com/rthornton128/goncurses"
	"log"
	"math/rand"

	"github.com/cmcfarlen/nethack/generate"
)

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

type World struct {
	currentFloor *generate.Dungeon
}

func UpdateMonster(w *World, m *Monster) {
	if rand.Float32() < 0.4 {
		m.V = randomDirection()
	}
	if rand.Float32() < 0.8 {
		m.P.x, m.P.y = moveIfValid(m.P.x, m.P.y, m.V.x, m.V.y, &w.currentFloor.M)
	}
}

func moveIfValid(x, y, dx, dy int, m *generate.Map) (int, int) {
	nx, ny := x+dx, y+dy

	if m.IsWalkable(nx, ny) {
		return nx, ny
	}
	return x, y
}

func randomWalkablePoint(m *generate.Map) (int, int) {
	x, y := rand.Intn(m.Width), rand.Intn(m.Height)
	for !m.IsWalkable(x, y) {
		x, y = rand.Intn(m.Width), rand.Intn(m.Height)
	}
	return x, y
}

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

	opts := generate.Opts{
		Width:             20,
		Height:            20,
		Sparseness:        70,
		DirectionModifier: 20,
		RoomCount:         4,
		RoomMin:           3,
		RoomMax:           7,
	}

	w := &World{}
	w.currentFloor = generate.GenerateDungeon(opts)

	px, py := randomWalkablePoint(&w.currentFloor.M)

	m := Monster{gc.Char('M'), Vector{}, Vector{}}
	m.P.x, m.P.y = randomWalkablePoint(&w.currentFloor.M)

	keepGoing := true
	for keepGoing {

		stdscr.Erase()

		d := w.currentFloor

		for x := 0; x < d.M.Width; x = x + 1 {
			for y := 0; y < d.M.Height; y = y + 1 {
				stdscr.MoveAddChar(y, x, gc.Char(d.M.RuneAt(x, y)))
			}
		}

		stdscr.MoveAddChar(py, px, gc.ACS_DIAMOND)
		UpdateMonster(w, &m)

		stdscr.MoveAddChar(m.P.y, m.P.x, m.Symbol)

		stdscr.Refresh()

		ch := stdscr.GetChar()

		switch ch {
		case gc.KEY_ESC:
			keepGoing = false
		case 'k':
			px, py = moveIfValid(px, py, 0, -1, &d.M)
		case 'j':
			px, py = moveIfValid(px, py, 0, 1, &d.M)
		case 'h':
			px, py = moveIfValid(px, py, -1, 0, &d.M)
		case 'l':
			px, py = moveIfValid(px, py, 1, 0, &d.M)
		case 'g':
			w.currentFloor = generate.GenerateDungeon(opts)
		}
	}
}
