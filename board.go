package main

import (
	"strings"
	"testing"
)

// newBoard returns an n×n board initialized with '.' (empty) cells
func newBoard(n int) *board {
	chars := make([][]byte, n)
	for i := range chars {
		chars[i] = make([]byte, n)
		for j := range chars[i] {
			chars[i][j] = '.'
		}
	}
	return &board{n: n, chars: chars}
}

// firstEmpty returns the first empty cell (top-left) or false if none
func (b *board) firstEmpty() (int, int, bool) {
	for y := 0; y < b.n; y++ {
		for x := 0; x < b.n; x++ {
			if b.chars[y][x] == '.' {
				return x, y, true
			}
		}
	}
	return 0, 0, false
}

// canPlace checks whether shape s can be placed at origin (offsetX, offsetY) on '.'
func (b *board) canPlace(s shape, offsetX, offsetY int) bool {
	for _, c := range s.cells {
		x, y := offsetX+c.x, offsetY+c.y
		if x < 0 || x >= b.n || y < 0 || y >= b.n || b.chars[y][x] != '.' {
			return false
		}
	}
	return true
}

// place writes letter for all cells of shape s at origin (offsetX,offsetY)
func (b *board) place(s shape, offsetX, offsetY int, letter byte) {
	for _, c := range s.cells {
		x, y := offsetX+c.x, offsetY+c.y
		b.chars[y][x] = letter
	}
}

// remove restores '.' for all cells of shape s at origin (offsetX,offsetY)
func (b *board) remove(s shape, offsetX, offsetY int) {
	for _, c := range s.cells {
		x, y := offsetX+c.x, offsetY+c.y
		if x >= 0 && x < b.n && y >= 0 && y < b.n {
			b.chars[y][x] = '.'
		}
	}
}

// String returns the board as newline-separated rows
func (b *board) String() string {
	var sb strings.Builder
	for y := 0; y < b.n; y++ {
		sb.Write(b.chars[y])
		if y < b.n-1 {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func TestBoard_String(t *testing.T) {
	b := newBoard(3)

	// expected all dots
	want := "...\n...\n..."
	if got := b.String(); got != want {
		t.Fatalf("String() = %q, want %q", got, want)
	}

	// place some letters
	b.chars[0][0] = 'A'
	b.chars[1][1] = 'B'
	b.chars[2][2] = 'C'

	want = "A..\n.B.\n..C"
	if got := b.String(); got != want {
		t.Fatalf("after modifications, String() = %q, want %q", got, want)
	}

	// ensure no trailing newline at end
	got := b.String()
	if got[len(got)-1] == '\n' {
		t.Fatal("unexpected trailing newline at end of String()")
	}
}
