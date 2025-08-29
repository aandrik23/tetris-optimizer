package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage: go run . <filename>")
		os.Exit(1)
	}

	filename := os.Args[1]
	shapes, err := readShapes(filename)
	if err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}

	board, err := solve(shapes)
	if err != nil {
		fmt.Println("ERROR")
		os.Exit(1)
	}
	fmt.Println(board)
}
