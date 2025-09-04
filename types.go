package main

import (
	"errors"
)

var errNoSolution = errors.New("no solution")
var errInvalidFormat = errors.New("ERROR")

// shape represents a tetromino shape with its cells and dimensions
// Coordinates are always normalized to start from (0,0)
type shape struct {
	cells         []cell
	width, height int
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
	board  *board
	pieces []piece

	// cache: canonical key of shape -> orientations
	oriCache map[string][]shape
}
