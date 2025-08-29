package main

import (
	"fmt"
	"math"
	"os"
)

// solve is the main entry point building up the smallest feasible board and
// running DFS with MRV + caching.
func solve(shapes []shape) (*board, error) {
	if len(shapes) == 0 {
		return nil, errNoSolution
	}

	pieces := make([]piece, len(shapes))
	for i, sh := range shapes {
		pieces[i] = piece{letter: byte('A' + i), shape: sh}
	}

	// Calculate minimum board size (ceil(sqrt(4 * piece_count)))
	minSize := int(math.Ceil(math.Sqrt(float64(4 * len(shapes)))))
	if minSize < 1 {
		minSize = 1
	}

	// Try board sizes from minimum to maximum (8x8)
	for n := minSize; n <= 8; n++ {
		totalCells := n * n
		coveredCells := 4 * len(shapes)
		emptyCellsAllowed := totalCells - coveredCells

		// Skip sizes that can't possibly fit all pieces
		if emptyCellsAllowed < 0 {
			continue
		}

		fmt.Fprintf(os.Stderr, "trying n=%d (empty cells allowed: %d)...\n", n, emptyCellsAllowed)

		// Initialize board and solver
		b := newBoard(n)
		s := &solver{
			b:         b,
			pieces:    pieces,
			oriCache:  make(map[string][]shape),
			candCache: make(map[candKey][]cand),
		}

		// Attempt to solve with current board size
		if s.dfsPlaceMRVFirstHole(0, emptyCellsAllowed) {
			fmt.Fprintf(os.Stderr, "n=%d solved\n", n)
			return b, nil
		}
	}

	return nil, errNoSolution
}

// dfsPlaceMRVFirstHole tries to place pieces using MRV on the current first
// empty cell. It also supports leaving exactly `emptyCellsAllowed` holes.
func (s *solver) dfsPlaceMRVFirstHole(idx int, emptyCellsAllowed int) bool {
	if idx == len(s.pieces) {
		// verify exact number of empty cells left
		emptyCount := 0
		for y := 0; y < s.b.n; y++ {
			for x := 0; x < s.b.n; x++ {
				if s.b.chars[y][x] == '.' {
					emptyCount++
				}
			}
		}
		return emptyCount == emptyCellsAllowed
	}

	// Pick first empty cell (top-left)
	tx, ty, ok := s.b.firstEmpty()
	if !ok {
		return idx == len(s.pieces) && emptyCellsAllowed == 0
	}

	// Try ALL pieces A..Z for this cell, before considering a hole
	for i := idx; i < len(s.pieces); i++ {
		cands := s.candidatesForFixedPiece(i, tx, ty)
		if len(cands) == 0 {
			continue
		}

		if i != idx {
			s.pieces[idx], s.pieces[i] = s.pieces[i], s.pieces[idx]
		}
		letter := s.pieces[idx].letter

		for _, c := range cands {
			if !s.b.canPlace(c.ori, c.ox, c.oy) {
				continue
			}
			s.b.place(c.ori, c.ox, c.oy, letter)
			if s.checkHolesDivisibleBy4() && s.dfsPlaceMRVFirstHole(idx+1, emptyCellsAllowed) {
				return true
			}
			s.b.remove(c.ori, c.ox, c.oy)
		}

		if i != idx {
			s.pieces[idx], s.pieces[i] = s.pieces[i], s.pieces[idx]
		}
	}

	// Use a hole only if nothing fits this cell
	if emptyCellsAllowed > 0 {
		orig := s.b.chars[ty][tx]
		s.b.chars[ty][tx] = ' '
		ok := s.dfsPlaceMRVFirstHole(idx, emptyCellsAllowed-1)
		s.b.chars[ty][tx] = orig
		return ok
	}
	return false
}

// checkHolesDivisibleBy4 returns false if any connected component of '.' cells
// has size not divisible by 4. It ignores ' ' holes by design.
func (s *solver) checkHolesDivisibleBy4() bool {
	n := s.b.n
	seen := make([]bool, n*n)
	idx := func(x, y int) int { return y*n + x }

	// Reusable queue slices to avoid allocations
	qx := make([]int, 0, n*n)
	qy := make([]int, 0, n*n)

	for y := 0; y < n; y++ {
		for x := 0; x < n; x++ {
			if s.b.chars[y][x] != '.' || seen[idx(x, y)] {
				continue
			}

			// BFS over '.' cells
			qx = qx[:0]
			qy = qy[:0]
			qx = append(qx, x)
			qy = append(qy, y)
			seen[idx(x, y)] = true

			size := 0
			for head := 0; head < len(qx); head++ {
				cx, cy := qx[head], qy[head]
				size++
				// 4-neighborhood
				if cx+1 < n && s.b.chars[cy][cx+1] == '.' && !seen[idx(cx+1, cy)] {
					seen[idx(cx+1, cy)] = true
					qx = append(qx, cx+1)
					qy = append(qy, cy)
				}
				if cx-1 >= 0 && s.b.chars[cy][cx-1] == '.' && !seen[idx(cx-1, cy)] {
					seen[idx(cx-1, cy)] = true
					qx = append(qx, cx-1)
					qy = append(qy, cy)
				}
				if cy+1 < n && s.b.chars[cy+1][cx] == '.' && !seen[idx(cx, cy+1)] {
					seen[idx(cx, cy+1)] = true
					qx = append(qx, cx)
					qy = append(qy, cy+1)
				}
				if cy-1 >= 0 && s.b.chars[cy-1][cx] == '.' && !seen[idx(cx, cy-1)] {
					seen[idx(cx, cy-1)] = true
					qx = append(qx, cx)
					qy = append(qy, cy-1)
				}
			}

			if size%4 != 0 {
				return false // prune: impossible to fill this island with tetrominoes
			}
		}
	}
	return true
}
