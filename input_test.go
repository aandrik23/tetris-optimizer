package main

import (
	"os"
	"path/filepath"
	"testing"
)

// helper to write a temp file with content
func writeTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}

// helper to count shape cells
func cellCount(s shape) int {
	return len(s.cells)
}

func TestReadShapes_TwoShapes_WithBlankSeparator(t *testing.T) {
	// two 4x4 blocks (O-tetromino style), separated by a blank line, no trailing newline
	content := "" +
		"....\n" +
		".##.\n" +
		".##.\n" +
		"....\n" +
		"\n" +
		"....\n" +
		"##..\n" +
		"##..\n" +
		"...."

	path := writeTempFile(t, "shapes.txt", content)

	shs, err := readShapes(path)
	if err != nil {
		t.Fatalf("readShapes error: %v", err)
	}
	if len(shs) != 2 {
		t.Fatalf("expected 2 shapes, got %d", len(shs))
	}
	// expect each to have 4 cells (tetromino)
	if cellCount(shs[0]) != 4 || cellCount(shs[1]) != 4 {
		t.Fatalf("expected 4 cells per shape, got %d and %d", cellCount(shs[0]), cellCount(shs[1]))
	}
}

func TestReadShapes_SingleShape_NoTrailingNewline(t *testing.T) {
	// single shape without trailing newline
	content := "" +
		"....\n" +
		".##.\n" +
		".##.\n" +
		"...." // no \n at EOF

	path := writeTempFile(t, "one.txt", content)

	shs, err := readShapes(path)
	if err != nil {
		t.Fatalf("readShapes error: %v", err)
	}
	if len(shs) != 1 {
		t.Fatalf("expected 1 shape, got %d", len(shs))
	}
	if cellCount(shs[0]) != 4 {
		t.Fatalf("expected 4 cells, got %d", cellCount(shs[0]))
	}
}

func TestReadShapes_EmptyFile_Error(t *testing.T) {
	path := writeTempFile(t, "empty.txt", "")

	shs, err := readShapes(path)
	if err == nil {
		t.Fatalf("expected error for empty file, got nil with %d shapes", len(shs))
	}
	if err != errInvalidFormat {
		t.Fatalf("expected errInvalidFormat, got %v", err)
	}
}

func TestReadShapes_ExtraBlankLines_AreIgnored(t *testing.T) {
	// leading, multiple mid, and trailing blank lines
	content := "" +
		"\n\n" +
		"....\n" +
		".##.\n" +
		".##.\n" +
		"....\n" +
		"\n\n\n" +
		"....\n" +
		"..##\n" +
		"..##\n" +
		"....\n" +
		"\n"

	path := writeTempFile(t, "blanks.txt", content)

	shs, err := readShapes(path)
	if err != nil {
		t.Fatalf("readShapes error: %v", err)
	}
	if len(shs) != 2 {
		t.Fatalf("expected 2 shapes, got %d", len(shs))
	}
}

func TestAddShape_EmptyLines_NoOp(t *testing.T) {
	var shapes []shape
	// passing empty lines should not append a shape
	out, err := addShape(nil, shapes)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Fatalf("expected 0 shapes, got %d", len(out))
	}
}

func TestAddShape_ValidLines_Appends(t *testing.T) {
	lines := []string{
		"....",
		".##.",
		".##.",
		"....",
	}
	var shapes []shape

	out, err := addShape(lines, shapes)
	if err != nil {
		t.Fatalf("addShape error: %v", err)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 shape, got %d", len(out))
	}
	if cellCount(out[0]) != 4 {
		t.Fatalf("expected 4 cells, got %d", cellCount(out[0]))
	}
}

func TestAddShape_InvalidLines_Error(t *testing.T) {
	// malformed lines (width mismatch or invalid chars)
	badCases := [][]string{
		{"##", "#"},           // jagged rows
		{"AB..", "...."},      // invalid chars
		{"...", "..", "...."}, // inconsistent widths
	}

	for i, lines := range badCases {
		var shapes []shape
		if _, err := addShape(lines, shapes); err == nil {
			t.Fatalf("case %d: expected error for invalid lines, got nil", i)
		}
	}
}
