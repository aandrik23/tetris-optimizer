package main

import (
	"bufio"
	"io"
	"os"
	"strings"
)

// readShapes reads tetromino shapes from a file.
// Shapes are separated by blank lines.
func readShapes(filename string) ([]shape, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	shapes := make([]shape, 0, 16) // small prealloc
	current := make([]string, 0, 4)

	// flush appends the current buffer as a shape and resets it
	flush := func() error {
		var err error
		shapes, err = addShape(current, shapes)
		current = current[:0]
		return err
	}

	for {
		line, rerr := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		eof := rerr == io.EOF

		// treat all non-EOF errors as fatal
		if rerr != nil && !eof {
			return nil, rerr
		}

		// add last line if file ends without trailing newline
		if eof && line != "" {
			current = append(current, line)
		}

		// blank line or EOF terminates a shape
		if line == "" || eof {
			if err := flush(); err != nil {
				return nil, err
			}
			if eof {
				break
			}
			continue
		}

		// accumulate non-blank line
		current = append(current, line)
	}

	if len(shapes) == 0 {
		return nil, errInvalidFormat
	}
	return shapes, nil
}

// addShape parses a set of lines into a shape and appends it.
func addShape(lines []string, shapes []shape) ([]shape, error) {
	if len(lines) == 0 {
		return shapes, nil
	}
	sh, err := parseShape(lines)
	if err != nil {
		return shapes, err
	}
	return append(shapes, sh), nil
}
