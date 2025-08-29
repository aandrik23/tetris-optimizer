package main

import "strings"

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

// canPlace checks whether shape s can be placed at origin (ox,oy) on '.'
func (b *board) canPlace(s shape, ox, oy int) bool {
	for _, c := range s.cells {
		x, y := ox+c.x, oy+c.y
		if x < 0 || x >= b.n || y < 0 || y >= b.n || b.chars[y][x] != '.' {
			return false
		}
	}
	return true
}

// place writes letter for all cells of shape s at origin (ox,oy)
func (b *board) place(s shape, ox, oy int, letter byte) {
	for _, c := range s.cells {
		x, y := ox+c.x, oy+c.y
		b.chars[y][x] = letter
	}
}

// remove restores '.' for all cells of shape s at origin (ox,oy)
func (b *board) remove(s shape, ox, oy int) {
	for _, c := range s.cells {
		x, y := ox+c.x, oy+c.y
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
