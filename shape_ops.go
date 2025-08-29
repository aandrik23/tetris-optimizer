package main

import (
	"sort"
	"strings"
)

// parseShape parses a 4x4 grid into a normalized shape
// (top-left cell becomes (0,0); width/height are adjusted accordingly).
func parseShape(grid []string) (shape, error) {
	cells := make([]cell, 0, 4)
	minX, minY := 3, 3
	maxX, maxY := 0, 0
	count := 0

	for y, row := range grid {
		if len(row) != 4 {
			return shape{}, errInvalidFormat
		}
		for x, ch := range row {
			if ch == '#' {
				count++
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
				cells = append(cells, cell{x, y})
			} else if ch != '.' {
				return shape{}, errInvalidFormat
			}
		}
	}

	// After collecting cells and checking count==4
	if count != 4 {
		return shape{}, errInvalidFormat
	}

	// Connectivity check: BFS/DFS over the 4 cells
	visited := make(map[cell]bool)
	queue := []cell{cells[0]}
	visited[cells[0]] = true

	dirs := []cell{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for i := 0; i < len(queue); i++ {
		cur := queue[i]
		for _, d := range dirs {
			nb := cell{cur.x + d.x, cur.y + d.y}
			for _, c := range cells {
				if c == nb && !visited[c] {
					visited[c] = true
					queue = append(queue, c)
				}
			}
		}
	}

	if len(visited) != 4 {
		return shape{}, errInvalidFormat
	}

	// Normalize to (0,0)
	for i := range cells {
		cells[i].x -= minX
		cells[i].y -= minY
	}
	w := maxX - minX + 1
	h := maxY - minY + 1

	return shape{cells: cells, w: w, h: h}, nil
}

// rotate90 rotates a shape 90 degrees clockwise.
func rotate90(s shape) shape {
	newCells := make([]cell, len(s.cells))
	for i, c := range s.cells {
		newCells[i] = cell{s.h - 1 - c.y, c.x}
	}
	return shape{cells: newCells, w: s.h, h: s.w}
}

// flipHorizontal flips a shape horizontally.
func flipHorizontal(s shape) shape {
	newCells := make([]cell, len(s.cells))
	for i, c := range s.cells {
		newCells[i] = cell{s.w - 1 - c.x, c.y}
	}
	return shape{cells: newCells, w: s.w, h: s.h}
}

// keyFor creates a unique key for a set of cells (for canonical identity).
func keyFor(cells []cell) string {
	cellsCopy := make([]cell, len(cells))
	copy(cellsCopy, cells)

	sort.Slice(cellsCopy, func(i, j int) bool {
		if cellsCopy[i].y != cellsCopy[j].y {
			return cellsCopy[i].y < cellsCopy[j].y
		}
		return cellsCopy[i].x < cellsCopy[j].x
	})

	var sb strings.Builder
	for _, c := range cellsCopy {
		sb.WriteByte(byte(c.x))
		sb.WriteByte(byte(c.y))
	}
	return sb.String()
}

// orientationsCanonical returns unique orientations using rotations only (no flips)
func (s *solver) orientationsCanonical(sh shape) []shape {
	k := keyFor(sh.cells)
	if oris, ok := s.oriCache[k]; ok {
		return oris
	}

	seen := make(map[string]struct{})
	orientations := make([]shape, 0, 4)

	push := func(shp shape) {
		cellsCopy := append([]cell(nil), shp.cells...)
		sort.Slice(cellsCopy, func(i, j int) bool {
			if cellsCopy[i].y != cellsCopy[j].y {
				return cellsCopy[i].y < cellsCopy[j].y
			}
			return cellsCopy[i].x < cellsCopy[j].x
		})
		var sb strings.Builder
		for _, c := range cellsCopy {
			sb.WriteByte(byte(c.x))
			sb.WriteByte(byte(c.y))
		}
		keyStr := sb.String()
		if _, ok := seen[keyStr]; !ok {
			seen[keyStr] = struct{}{}
			orientations = append(orientations, shp)
		}
	}

	base := sh
	for r := 0; r < 4; r++ {
		push(base) // push current rotation
		base = rotate90(base)
	}

	s.oriCache[k] = orientations
	return orientations
}
