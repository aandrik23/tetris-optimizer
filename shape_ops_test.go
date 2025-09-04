package main

import (
	"testing"
)

func TestParseShape_ValidO(t *testing.T) {
	// 2x2 O tetromino
	block := []string{
		"##..",
		"##..",
		"....",
		"....",
	}
	sh, err := parseShape(block)
	if err != nil {
		t.Fatalf("expected valid O piece, got error: %v", err)
	}
	if len(sh.cells) != 4 {
		t.Fatalf("expected 4 cells, got %d", len(sh.cells))
	}
	// expect normalized origin at (0,0)
	want := map[[2]int]bool{{0, 0}: true, {1, 0}: true, {0, 1}: true, {1, 1}: true}
	for _, c := range sh.cells {
		if !want[[2]int{c.x, c.y}] {
			t.Fatalf("unexpected cell in O: (%d,%d)", c.x, c.y)
		}
	}
	if sh.width != 2 || sh.height != 2 {
		t.Fatalf("expected width=2 height=2, got %d x %d", sh.width, sh.height)
	}
}

func TestParseShape_InvalidChar(t *testing.T) {
	// invalid character 'X'
	block := []string{
		"##..",
		"#X..",
		"#...",
		"....",
	}
	if _, err := parseShape(block); err == nil {
		t.Fatal("expected error for invalid character, got nil")
	}
}

func TestParseShape_NotFourRows(t *testing.T) {
	// only 3 rows
	block := []string{
		"####",
		"....",
		"....",
	}
	if _, err := parseShape(block); err == nil {
		t.Fatal("expected error for rows != 4, got nil")
	}
}

func TestParseShape_NotFourHashes(t *testing.T) {
	// only 3 hashes
	block := []string{
		"###.",
		"....",
		"....",
		"....",
	}
	if _, err := parseShape(block); err == nil {
		t.Fatal("expected error for != 4 filled cells, got nil")
	}
}

func TestParseShape_Disconnected(t *testing.T) {
	// 2+2 disconnected
	block := []string{
		"##..",
		"..##",
		"....",
		"....",
	}
	if _, err := parseShape(block); err == nil {
		t.Fatal("expected error for disconnected cells, got nil")
	}
}

func TestKeyFor_Basic(t *testing.T) {
	a := []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}} // O piece
	b := []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}}

	ka := keyFor(a)
	kb := keyFor(b)
	if ka != kb {
		t.Fatalf("expected same key for identical shapes, got %q vs %q", ka, kb)
	}

	c := []cell{{0, 0}, {1, 0}, {2, 0}, {3, 0}} // I piece
	kc := keyFor(c)
	if ka == kc {
		t.Fatalf("expected different keys for different shapes, got %q vs %q", ka, kc)
	}
}

func TestOrientationsCanonical_O(t *testing.T) {
	// O piece should have exactly 1 unique orientation
	o := shape{cells: []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, width: 2, height: 2}
	s := &solver{oriCache: make(map[string][]shape)}
	oris := s.orientationsCanonical(o)
	if len(oris) != 1 {
		t.Fatalf("expected 1 orientation for O, got %d", len(oris))
	}
}

func TestOrientationsCanonical_I(t *testing.T) {
	// I piece should have 2 unique orientations (vertical/horizontal)
	i := shape{cells: []cell{{0, 0}, {0, 1}, {0, 2}, {0, 3}}, width: 1, height: 4}
	s := &solver{oriCache: make(map[string][]shape)}
	oris := s.orientationsCanonical(i)
	if len(oris) != 2 {
		t.Fatalf("expected 2 orientations for I, got %d", len(oris))
	}
}

func TestOrientationsCanonical_L(t *testing.T) {
	// L piece should have 4 unique orientations
	l := shape{cells: []cell{{0, 0}, {0, 1}, {0, 2}, {1, 2}}, width: 2, height: 3}
	s := &solver{oriCache: make(map[string][]shape)}
	oris := s.orientationsCanonical(l)
	if len(oris) != 4 {
		t.Fatalf("expected 4 orientations for L, got %d", len(oris))
	}
}

func TestRotate90_Rectangle(t *testing.T) {
	// shape 2x3: cells at (0,0),(1,0),(0,1),(1,2) just for diversity
	s := shape{
		cells:  []cell{{0, 0}, {1, 0}, {0, 1}, {1, 2}},
		width:  2,
		height: 3,
	}
	r := rotate90(s)

	// check swapped dimensions
	if r.width != s.height || r.height != s.width {
		t.Fatalf("expected dims %dx%d, got %dx%d", s.height, s.width, r.width, r.height)
	}

	// verify each (x,y) -> (h-1-y, x)
	want := map[[2]int]bool{
		{3 - 1 - 0, 0}: true, // (0,0)->(2,0)
		{3 - 1 - 0, 1}: true, // (1,0)->(2,1)
		{3 - 1 - 1, 0}: true, // (0,1)->(1,0)
		{3 - 1 - 2, 1}: true, // (1,2)->(0,1)
	}
	for _, c := range r.cells {
		if !want[[2]int{c.x, c.y}] {
			t.Fatalf("unexpected rotated cell: (%d,%d)", c.x, c.y)
		}
	}
}

func TestRotate90_FourTimesReturnsOriginal(t *testing.T) {
	s := shape{
		cells:  []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
		width:  2,
		height: 2,
	}
	r := s
	for i := 0; i < 4; i++ {
		r = rotate90(r)
	}
	// compare sets
	seen := map[[2]int]bool{}
	for _, c := range s.cells {
		seen[[2]int{c.x, c.y}] = true
	}
	for _, c := range r.cells {
		if !seen[[2]int{c.x, c.y}] {
			t.Fatalf("after 4 rotations expected original cells, missing (%d,%d)", c.x, c.y)
		}
	}
	if r.width != s.width || r.height != s.height {
		t.Fatalf("expected same dims after 4x rotate, got %dx%d vs %dx%d", r.width, r.height, s.width, s.height)
	}
}
