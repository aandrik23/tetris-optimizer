package main

import "testing"

func mustParse(t *testing.T, lines []string) shape {
	t.Helper()
	sh, err := parseShape(lines)
	if err != nil {
		t.Fatalf("parseShape failed: %v", err)
	}
	return sh
}

func TestPlaceAllPiecesExact_SingleOIn2x2(t *testing.T) {
	// one O should exactly fill 2x2
	o := mustParse(t, []string{
		"##..",
		"##..",
		"....",
		"....",
	})
	b := newBoard(2)
	s := &solver{
		board:    b,
		pieces:   []piece{{letter: 'A', shape: o}},
		oriCache: make(map[string][]shape),
	}
	if !s.placeAllPiecesExact(0) {
		t.Fatal("expected exact tiling for single O in 2x2")
	}
	// ensure no '.' remains
	for y := 0; y < b.n; y++ {
		for x := 0; x < b.n; x++ {
			if b.chars[y][x] == '.' {
				t.Fatalf("expected full cover, found '.' at (%d,%d)", x, y)
			}
		}
	}
}

func TestPlaceAllPiecesExact_IILZ_FailsOn4x4(t *testing.T) {
	// I (vertical), I (horizontal), L, Z cannot exactly tile a 4x4
	Iv := mustParse(t, []string{
		"...#",
		"...#",
		"...#",
		"...#",
	})
	Ih := mustParse(t, []string{
		"....",
		"....",
		"....",
		"####",
	})
	L := mustParse(t, []string{
		".###",
		"...#",
		"....",
		"....",
	})
	Z := mustParse(t, []string{
		"....",
		"..##",
		".##.",
		"....",
	})

	b := newBoard(4)
	s := &solver{
		board:    b,
		pieces:   []piece{{'A', L}, {'B', Z}, {'C', Iv}, {'D', Ih}},
		oriCache: make(map[string][]shape),
	}

	// exact cover should be impossible
	if s.placeAllPiecesExact(0) {
		t.Fatal("expected no exact tiling for {I,I,L,Z} on 4x4")
	}
}

func TestSolve_Exact4OOn4x4(t *testing.T) {
	// four O pieces should solve a 4x4 exactly
	o := mustParse(t, []string{
		"##..",
		"##..",
		"....",
		"....",
	})

	shapes := []shape{o, o, o, o}
	board, err := solve(shapes) // your exact solver version that tries n with exact area
	if err != nil {
		t.Fatalf("expected solution, got err: %v", err)
	}
	// basic sanity: all cells should be letters
	for y := 0; y < board.n; y++ {
		for x := 0; x < board.n; x++ {
			if board.chars[y][x] == '.' {
				t.Fatalf("expected no '.', found at (%d,%d)", x, y)
			}
		}
	}
}

func TestPlaceAllPiecesExact_OOOI_FailsOn4x4(t *testing.T) {
	// three O plus one I cannot exactly tile a 4x4
	o := mustParse(t, []string{
		"##..",
		"##..",
		"....",
		"....",
	})
	I := mustParse(t, []string{
		"####",
		"....",
		"....",
		"....",
	})

	b := newBoard(4)
	s := &solver{
		board:    b,
		pieces:   []piece{{'A', o}, {'B', o}, {'C', o}, {'D', I}},
		oriCache: make(map[string][]shape),
	}
	// exact cover must be impossible
	if s.placeAllPiecesExact(0) {
		t.Fatal("expected no exact tiling for {O,O,O,I} on 4x4")
	}
}
