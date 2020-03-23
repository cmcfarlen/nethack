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
	m.V = randomDirection()
	m.P = clampZeroToV(add(m.P, m.V), w.currentFloor.Width, w.currentFloor.Height)
}

/*
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

*/

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

	opts := generate.GenerateOpts{
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

	height, width := w.currentFloor.M.Height, w.currentFloor.M.Width

	px := width / 2
	py := height / 2

	m := Monster{gc.Char('M'), Vector{}, Vector{}}

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
			w.currentFloor = generate.GenerateDungeon(opts)
		}
	}
}
