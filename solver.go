package main

import (
	"fmt"
	"math"
	"os"
)

// solve tries increasing board sizes from the theoretical minimum up to a dynamic upper bound
func solve(shapes []shape) (*board, error) {
	if len(shapes) == 0 {
		return nil, errNoSolution
	}

	// build pieces with letters A, B, C, ...
	pieces := make([]piece, len(shapes))
	for i, sh := range shapes {
		pieces[i] = piece{letter: byte('A' + i), shape: sh}
	}

	// minimum n = ceil(sqrt(4 * tetrominoes))
	minSize := max(int(math.Ceil(math.Sqrt(float64(4*len(shapes))))), 1)

	// dynamic upper bound: enough to pack tetrominoes 4x4 "bins" (loose but safe)
	tetrominoes := len(shapes)
	maxN := max(4*int(math.Ceil(math.Sqrt(float64(tetrominoes)))), minSize)

	for n := minSize; n <= maxN; n++ {
		totalCells := n * n
		coveredCells := 4 * tetrominoes
		emptyCellsAllowed := totalCells - coveredCells

		// skip sizes that cannot fit even theoretically
		if emptyCellsAllowed < 0 {
			continue
		}

		fmt.Fprintf(os.Stderr, "trying n=%d (empty cells allowed: %d)...\n", n, emptyCellsAllowed)

		// initialize board and solver
		board := newBoard(n)
		s := &solver{
			board:    board,
			pieces:   pieces,
			oriCache: make(map[string][]shape),
		}

		if s.placeAllPiecesExact(0) {
			return board, nil
		}
	}

	return nil, errNoSolution
}

// placeAllPiecesExact tries to exactly tile the board with all pieces (no holes)
func (s *solver) placeAllPiecesExact(idx int) bool {
	// all pieces placed
	if idx == len(s.pieces) {
		return true
	}

	// pick first empty cell
	x, y, ok := s.board.firstEmpty()
	if !ok {
		return idx == len(s.pieces) // no '.' left means all must be placed
	}

	// try each remaining piece
	for i := idx; i < len(s.pieces); i++ {
		// swap-in piece i at position idx
		if i != idx {
			s.pieces[idx], s.pieces[i] = s.pieces[i], s.pieces[idx]
		}
		p := s.pieces[idx]

		// try each unique orientation
		orientation := s.orientationsCanonical(p.shape)
		for _, orient := range orientation {
			// anchor each cell of the orientation onto (x,y)
			for _, cell := range orient.cells {
				offsetX, offsetY := x-cell.x, y-cell.y // compute origin offset
				if !s.board.canPlace(orient, offsetX, offsetY) {
					continue // blocked by bounds or occupancy
				}
				s.board.place(orient, offsetX, offsetY, p.letter) // commit placement
				if s.placeAllPiecesExact(idx + 1) {               // recurse to next piece
					return true // propagate success
				}
				s.board.remove(orient, offsetX, offsetY) // backtrack placement
			}
		}

		// restore original order before trying next i
		if i != idx {
			s.pieces[idx], s.pieces[i] = s.pieces[i], s.pieces[idx]
		}
	}

	// no feasible placement for current decision
	return false
}
