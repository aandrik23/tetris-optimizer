package main

import (
	"errors"
)

var errNoSolution = errors.New("no solution")
var errInvalidFormat = errors.New("ERROR")

// shape represents a tetromino shape with its cells and dimensions
// Coordinates are always normalized to start from (0,0)
type shape struct {
	cells []cell
	w, h  int
}

// cell represents a single coordinate in a shape
type cell struct {
	x, y int
}

// piece represents a tetromino with its identifying letter
type piece struct {
	letter byte
	shape  shape
}

// board represents the playing field
type board struct {
	n     int
	chars [][]byte
}

// solver encapsulates state and lightweight caches for performance
// (orientation and candidate caches avoid recomputation during DFS)
type solver struct {
	b      *board
	pieces []piece

	// cache: canonical key of shape -> orientations
	oriCache map[string][]shape

	// cache: (pieceIdx, tx, ty) -> candidates for fixed cell
	candCache map[candKey][]cand
}

// cand represents a single candidate placement for a given piece
type cand struct {
	ori    shape
	ox, oy int
	oIndex int // canonical orientation index for deterministic order
}

type candKey struct {
	pieceIdx int
	tx, ty   int
	n        int
}
