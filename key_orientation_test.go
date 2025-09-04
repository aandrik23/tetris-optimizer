package main

import (
	"cmp"
	"slices"
	"testing"
)

// helper to sort cells for deterministic compare
func sortCells(cs []cell) []cell {
	out := slices.Clone(cs)
	slices.SortFunc(out, func(a, b cell) int {
		if a.y != b.y {
			return cmp.Compare(a.y, b.y)
		}
		return cmp.Compare(a.x, b.x)
	})
	return out
}

func TestKeyFor_TranslationInvariant(t *testing.T) {
	// same shape moved by (+2,+1) should produce same key after normalization + sorting
	a := []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}}
	b := []cell{{2, 1}, {3, 1}, {2, 2}, {3, 2}}
	ka := keyFor(a)
	kb := keyFor(b)
	if ka != kb {
		t.Fatalf("expected same key for translated shapes, got %q vs %q", ka, kb)
	}
}

func TestOrientationsCanonical_O(t *testing.T) {
	// O piece should have exactly 1 unique orientation
	o := shape{cells: []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}}, width: 2, height: 2}
	s := &solver{oriCache: make(map[string][]shape)}
	oris := s.orientationsCanonical(o)
	if len(oris) != 1 {
		t.Fatalf("expected 1 orientation for O, got %d", len(oris))
	}
}

func TestOrientationsCanonical_I(t *testing.T) {
	// I piece should have 2 unique orientations (vertical/horizontal)
	i := shape{cells: []cell{{0, 0}, {0, 1}, {0, 2}, {0, 3}}, width: 1, height: 4}
	s := &solver{oriCache: make(map[string][]shape)}
	oris := s.orientationsCanonical(i)
	if len(oris) != 2 {
		t.Fatalf("expected 2 orientations for I, got %d", len(oris))
	}
}

func TestOrientationsCanonical_L(t *testing.T) {
	// L piece should have 4 unique orientations
	l := shape{cells: []cell{{0, 0}, {0, 1}, {0, 2}, {1, 2}}, width: 2, height: 3}
	s := &solver{oriCache: make(map[string][]shape)}
	oris := s.orientationsCanonical(l)
	if len(oris) != 4 {
		t.Fatalf("expected 4 orientations for L, got %d", len(oris))
	}
}
