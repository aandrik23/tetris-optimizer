package main

import "testing"

func TestBoard_CanPlace_Place_Remove(t *testing.T) {
	b := newBoard(4)
	o := shape{cells: []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, width: 2, height: 2}

	// can place at (1,1)
	if !b.canPlace(o, 1, 1) {
		t.Fatal("expected canPlace at (1,1)")
	}
	// place and verify letters
	b.place(o, 1, 1, 'A')
	if b.chars[1][1] != 'A' || b.chars[2][2] != 'A' {
		t.Fatal("expected 'A' cells after place")
	}
	// cannot place overlapping
	if b.canPlace(o, 1, 1) {
		t.Fatal("expected cannot place overlapping")
	}
	// remove and verify cleanup
	b.remove(o, 1, 1)
	if b.chars[1][1] != '.' || b.chars[2][2] != '.' {
		t.Fatal("expected '.' after remove")
	}
}

func TestBoard_FirstEmpty(t *testing.T) {
	b := newBoard(3)
	// fill (0,0)
	b.chars[0][0] = 'X'
	x, y, ok := b.firstEmpty()
	if !ok || x != 1 || y != 0 {
		t.Fatalf("expected first empty at (1,0), got (%d,%d), ok=%v", x, y, ok)
	}
}
