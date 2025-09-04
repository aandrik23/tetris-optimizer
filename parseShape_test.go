package main

import "testing"

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
