package main

import (
	"bufio"
	"os"
	"strings"
)

// readShapes reads the input file and parses a list of shapes
func readShapes(filename string) ([]shape, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	shapes := make([]shape, 0)
	current := make([]string, 0, 4)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			if len(current) > 0 {
				sh, err := parseShape(current)
				if err != nil {
					return nil, err
				}
				shapes = append(shapes, sh)
				current = current[:0]
			}
			continue
		}
		current = append(current, line)
	}

	if len(current) > 0 {
		sh, err := parseShape(current)
		if err != nil {
			return nil, err
		}
		shapes = append(shapes, sh)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(shapes) == 0 {
		return nil, errInvalidFormat
	}
	return shapes, nil
}
