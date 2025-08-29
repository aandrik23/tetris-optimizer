package main

import "sort"

// candidatesForFixedPiece returns canonical orientation placements for piece i,
// anchored so that one of the orientation's cells sits on (tx,ty). Results are
// sorted deterministically and deduplicated by (orientation-shape, ox, oy).
func (s *solver) candidatesForFixedPiece(pieceIdx, tx, ty int) []cand {
	key := candKey{pieceIdx: pieceIdx, tx: tx, ty: ty, n: s.b.n} // Χρήση struct
	if cached, ok := s.candCache[key]; ok {
		return cached
	}

	sh := s.pieces[pieceIdx].shape
	oris := s.orientationsCanonical(sh)
	out := make([]cand, 0, 32)

	for oi, o := range oris {
		// Sort cells by (y,x) to anchor consistently
		cells := append([]cell(nil), o.cells...)
		sort.Slice(cells, func(i, j int) bool {
			if cells[i].y != cells[j].y {
				return cells[i].y < cells[j].y
			}
			return cells[i].x < cells[j].x
		})

		// Collect positions for THIS orientation
		for _, c := range cells {
			ox := tx - c.x
			oy := ty - c.y
			if ox < 0 || oy < 0 || ox+o.w > s.b.n || oy+o.h > s.b.n {
				continue
			}
			out = append(out, cand{ori: o, ox: ox, oy: oy, oIndex: oi})
		}
	}

	// Sort all candidates by (oy, ox, oIndex)
	sort.Slice(out, func(i, j int) bool {
		if out[i].oy != out[j].oy {
			return out[i].oy < out[j].oy
		}
		if out[i].ox != out[j].ox {
			return out[i].ox < out[j].ox
		}
		return out[i].oIndex < out[j].oIndex
	})

	// Deduplicate - keep only unique (shape, ox, oy) combinations
	type dkey struct {
		k      string
		ox, oy int
	}
	seen := make(map[dkey]struct{}, len(out))
	dst := out[:0]
	for _, x := range out {
		kk := dkey{k: keyFor(x.ori.cells), ox: x.ox, oy: x.oy}
		if _, ok := seen[kk]; ok {
			continue
		}
		seen[kk] = struct{}{}
		dst = append(dst, x)
	}

	res := append([]cand(nil), dst...)
	s.candCache[key] = res
	return res
}
