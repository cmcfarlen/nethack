package generate

import "testing"

func TestNewDungeon(t *testing.T) {
	d := newDungeon(10, 10)

	for x := 0; x < 10; x++ {
		for y := 0; y < 10; y++ {
			tile := d.tileAt(x, y)
			if tile.tiletype != tileTypeUnvisited {
				t.Errorf("tile at %d, %d is not unvisited", x, y)
			}

			if tile.x != x {
				t.Errorf("tile at %d, %d has wrong x (%d)", x, y, tile.x)
			}

			if tile.y != y {
				t.Errorf("tile at %d, %d has wrong y (%d)", x, y, tile.y)
			}
		}
	}
}

func TestTileAt(t *testing.T) {
	d := newDungeon(10, 10)

	x, y := d.randomXY()
	tile := d.tileAt(x, y)

	if tile.x != x {
		t.Errorf("random tile at %d, %d has wrong x (%d)", x, y, tile.x)
	}

	if tile.y != y {
		t.Errorf("random tile at %d, %d has wrong y (%d)", x, y, tile.y)
	}
}

func TestDirectionPicker(t *testing.T) {
	d := newDirectionPicker()

	dir1 := d.nextDirection()
	dir2 := d.nextDirection()
	dir3 := d.nextDirection()
	dir4 := d.nextDirection()

	t.Logf("directions: %d %d %d %d", dir1, dir2, dir3, dir4)

	if d.hasNextDirection() {
		t.Errorf("picker still had directions after picking 4: %d", d.directionsRemaining)
	}

	if dir1 == dir2 {
		t.Errorf("picker picked the same direction: %d %d %d %d", dir1, dir2, dir3, dir4)
	}

}

/*
func TestNewDungeonPrint(t *testing.T) {
	opts := GenerateOpts{
		Width:  10,
		Height: 10,
	}
	d := GenerateDungeon(opts)

	d.PrintDungeon()
}
*/
