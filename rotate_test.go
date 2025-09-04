package main

import "testing"

func TestRotate90_Rectangle(t *testing.T) {
	// shape 2x3: cells at (0,0),(1,0),(0,1),(1,2) just for diversity
	s := shape{
		cells:  []cell{{0, 0}, {1, 0}, {0, 1}, {1, 2}},
		width:  2,
		height: 3,
	}
	r := rotate90(s)

	// check swapped dimensions
	if r.width != s.height || r.height != s.width {
		t.Fatalf("expected dims %dx%d, got %dx%d", s.height, s.width, r.width, r.height)
	}

	// verify each (x,y) -> (h-1-y, x)
	want := map[[2]int]bool{
		{3 - 1 - 0, 0}: true, // (0,0)->(2,0)
		{3 - 1 - 0, 1}: true, // (1,0)->(2,1)
		{3 - 1 - 1, 0}: true, // (0,1)->(1,0)
		{3 - 1 - 2, 1}: true, // (1,2)->(0,1)
	}
	for _, c := range r.cells {
		if !want[[2]int{c.x, c.y}] {
			t.Fatalf("unexpected rotated cell: (%d,%d)", c.x, c.y)
		}
	}
}

func TestRotate90_FourTimesReturnsOriginal(t *testing.T) {
	s := shape{
		cells:  []cell{{0, 0}, {1, 0}, {0, 1}, {1, 1}},
		width:  2,
		height: 2,
	}
	r := s
	for i := 0; i < 4; i++ {
		r = rotate90(r)
	}
	// compare sets
	seen := map[[2]int]bool{}
	for _, c := range s.cells {
		seen[[2]int{c.x, c.y}] = true
	}
	for _, c := range r.cells {
		if !seen[[2]int{c.x, c.y}] {
			t.Fatalf("after 4 rotations expected original cells, missing (%d,%d)", c.x, c.y)
		}
	}
	if r.width != s.width || r.height != s.height {
		t.Fatalf("expected same dims after 4x rotate, got %dx%d vs %dx%d", r.width, r.height, s.width, s.height)
	}
}
