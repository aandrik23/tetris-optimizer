package main

import "testing"

func TestNewBoard(t *testing.T) {
	b := newBoard(3)

	if b.n != 3 {
		t.Fatalf("expected size 3, got %d", b.n)
	}
	if len(b.chars) != 3 || len(b.chars[0]) != 3 {
		t.Fatalf("expected 3x3 chars, got %dx%d", len(b.chars), len(b.chars[0]))
	}
	for y := range 3 {
		for x := range 3 {
			if b.chars[y][x] != '.' {
				t.Fatalf("expected '.' at (%d,%d), got %q", x, y, b.chars[y][x])
			}
		}
	}

	// Check initial String()
	want := "...\n...\n..."
	if got := b.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}
}

func TestBoard_CanPlace_Place_Remove(t *testing.T) {
	b := newBoard(4)
	o := shape{cells: []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, width: 2, height: 2}

	t.Run("can place in-bounds", func(t *testing.T) {
		// can place at (1,1)
		if !b.canPlace(o, 1, 1) {
			t.Fatal("expected canPlace at (1,1)")
		}
	})

	t.Run("place fills all cells", func(t *testing.T) {
		// place and verify all four cells
		b.place(o, 1, 1, 'A')
		want := [][2]int{{1, 1}, {2, 1}, {1, 2}, {2, 2}}
		for _, p := range want {
			if b.chars[p[1]][p[0]] != 'A' {
				t.Fatalf("expected 'A' at (%d,%d)", p[0], p[1])
			}
		}
	})

	t.Run("cannot place overlapping", func(t *testing.T) {
		// overlapping with the previous 'A's must be false
		if b.canPlace(o, 1, 1) {
			t.Fatal("expected cannot place overlapping")
		}
	})

	t.Run("out of bounds cannot place", func(t *testing.T) {
		// top-left OK, bottom-right out-of-bounds should be false
		if b.canPlace(o, 3, 3) {
			t.Fatal("expected canPlace to be false out-of-bounds")
		}
		if b.canPlace(o, -1, 0) {
			t.Fatal("expected canPlace to be false for negative origin")
		}
	})

	t.Run("remove cleans exactly those cells", func(t *testing.T) {
		// remove and verify only those four cells become '.'
		b.remove(o, 1, 1)
		want := [][2]int{{1, 1}, {2, 1}, {1, 2}, {2, 2}}
		for _, p := range want {
			if b.chars[p[1]][p[0]] != '.' {
				t.Fatalf("expected '.' at (%d,%d) after remove", p[0], p[1])
			}
		}
	})

	t.Run("remove idempotent", func(t *testing.T) {
		// removing again should not panic and board remains '.'
		b.remove(o, 1, 1)
		for y := 0; y < b.n; y++ {
			for x := 0; x < b.n; x++ {
				if b.chars[y][x] != '.' {
					t.Fatalf("expected '.' at (%d,%d)", x, y)
				}
			}
		}
	})

	t.Run("place then String matches", func(t *testing.T) {
		// place at (0,0) and check board rendering
		b.place(o, 0, 0, 'B')
		got := b.String()
		want := "" +
			"BB..\n" +
			"BB..\n" +
			"....\n" +
			"...."
		if got != want {
			t.Fatalf("String() after place =\n%s\nwant:\n%s", got, want)
		}
	})
}

func TestBoard_FirstEmpty(t *testing.T) {
	b := newBoard(3)

	// first empty in a fresh board is (0,0)
	if x, y, ok := b.firstEmpty(); !ok || x != 0 || y != 0 {
		t.Fatalf("expected (0,0,true), got (%d,%d,%v)", x, y, ok)
	}

	// fill (0,0) and check moves to (1,0)
	b.chars[0][0] = 'X'
	x, y, ok := b.firstEmpty()
	if !ok || x != 1 || y != 0 {
		t.Fatalf("expected first empty at (1,0), got (%d,%d), ok=%v", x, y, ok)
	}

	// fill the rest and expect no empty
	for yy := 0; yy < 3; yy++ {
		for xx := 0; xx < 3; xx++ {
			b.chars[yy][xx] = 'X'
		}
	}
	if x, y, ok := b.firstEmpty(); ok {
		t.Fatalf("expected no empty cell, got (%d,%d)", x, y)
	}
}
