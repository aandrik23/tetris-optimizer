package main

import (
	"cmp"
	"slices"
	"strconv"
	"strings"
)

// parseShape parses a 4x4 grid into a normalized shape (0,0 top-left) and checks validity
func parseShape(block []string) (shape, error) {
	// require exactly 4 rows
	if len(block) != 4 {
		return shape{}, errInvalidFormat
	}

	cells := make([]cell, 0, 4)
	minX, minY := 3, 3
	maxX, maxY := 0, 0
	hashes := 0

	// scan the 4x4 block
	for y := range 4 {
		row := block[y]
		// each row must be exactly 4 chars
		if len(row) != 4 {
			return shape{}, errInvalidFormat
		}
		for x := range 4 {
			char := row[x]
			switch char {
			case '#':
				hashes++
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
			case '.':
				// ok: empty
			default:
				return shape{}, errInvalidFormat // invalid character
			}
		}
	}

	// must have exactly 4 filled cells
	if hashes != 4 {
		return shape{}, errInvalidFormat
	}

	// connectivity check on the 4 cells (4-neighborhood)
	present := make(map[cell]bool, 4)
	for _, cell := range cells {
		present[cell] = true // mark each filled cell
	}

	visited := make(map[cell]bool, 4)
	stack := []cell{cells[0]}
	visited[cells[0]] = true // mark start as visited

	dirs := [4]cell{
		{1, 0},  // right
		{-1, 0}, // left
		{0, 1},  // down
		{0, -1}, // up
	}

	for len(stack) > 0 {
		current := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		for _, dir := range dirs {
			neighbor := cell{current.x + dir.x, current.y + dir.y}

			if !present[neighbor] {
				continue // neighbor is not a filled cell
			}

			if visited[neighbor] {
				continue // already visited
			}

			visited[neighbor] = true
			stack = append(stack, neighbor)
		}
	}

	if len(visited) != 4 {
		return shape{}, errInvalidFormat // not a single connected tetromino
	}

	// normalize cells so that top-left becomes (0,0)
	for i := range cells {
		cells[i].x -= minX
		cells[i].y -= minY
	}

	// compute bounding box dimensions
	w := maxX - minX + 1
	h := maxY - minY + 1

	return shape{cells: cells, width: w, height: h}, nil
}

// orientationsCanonical returns unique orientations using rotations only (no flips).
func (s *solver) orientationsCanonical(sh shape) []shape {
	// use cache to avoid recomputing the 4 rotations for identical base shapes
	key := keyFor(sh.cells)
	if oris, ok := s.oriCache[key]; ok {
		return oris // return from cache if present
	}

	// seen holds keys of already added orientations
	seen := make(map[string]bool)
	orientations := make([]shape, 0, 4)

	push := func(shp shape) {
		keyStr := keyFor(shp.cells)
		if !seen[keyStr] {
			seen[keyStr] = true
			orientations = append(orientations, shp)
		}
	}

	base := sh
	for range 4 {
		push(base)            // push current rotation
		base = rotate90(base) // rotate for next iteration
	}

	s.oriCache[key] = orientations
	return orientations
}

// keyFor creates a unique key for a set of cells (for canonical identity).
func keyFor(cells []cell) string {
	cellsCopy := make([]cell, len(cells))
	copy(cellsCopy, cells)

	// sort by y first, then x
	slices.SortFunc(cellsCopy, func(a, b cell) int {
		if a.y != b.y {
			return cmp.Compare(a.y, b.y)
		}
		return cmp.Compare(a.x, b.x)
	})

	var sb strings.Builder
	for _, c := range cellsCopy {
		// use "x,y;" instead of raw bytes for safety
		sb.WriteString(strconv.Itoa(c.x))
		sb.WriteByte(',')
		sb.WriteString(strconv.Itoa(c.y))
		sb.WriteByte(';')
	}
	return sb.String()
}

// rotate90 rotates a shape 90 degrees clockwise.
// Matrix rotation 90° formula (adapted for grid indices):
// (x, y) → (height - 1 - y, x)
func rotate90(s shape) shape {
	newCells := make([]cell, len(s.cells))
	for i, c := range s.cells {
		newCells[i] = cell{s.height - 1 - c.y, c.x}
	}
	return shape{cells: newCells, width: s.height, height: s.width}
}
