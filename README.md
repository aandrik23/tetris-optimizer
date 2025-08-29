# Tetris Optimizer

A Go program that reads a list of tetrominoes from a text file and assembles them into the **smallest possible square**.  
Each tetromino is identified by a capital Latin letter (`A`, `B`, `C`, …) in the solution.

If it is impossible to form a complete square, the solver leaves empty spaces (`.`) between tetrominoes, as specified in the subject.

---

## Features

- Parses tetrominoes from a file in the given format.
- Places all tetrominoes using backtracking with pruning to find the smallest square.
- Supports **rotations (0°, 90°, 180°, 270°)** of tetrominoes.  
  Reflections (flips) are **not** allowed, matching the expected behavior of the examples.
- Validates tetrominoes: must be 4 connected cells (`#`) in a 4×4 grid.
- Prints the final board with:
  - `A..Z` for tetrominoes,
  - `.` for empty spaces.
- On invalid input format → prints `ERROR`.

---

## Requirements

- Go 1.20+
- Only standard library packages are used (no external dependencies).

---

## Build & Run

Compile:

```bash
go build -o tetris-optimizer
```

Run with a tetromino file:

```bash
./tetris-optimizer examples/sample.txt
```

Or directly:

```bash
go run . examples/sample.txt
```

---

## Input Format

- Each tetromino is described by **4 lines of 4 characters** (`#` for filled cell, `.` for empty).  
- Tetrominoes are separated by a blank line.

Example file:

```
#...
#...
#...
#...

....
....
..##
..##
```

---

## Output

The solver prints the assembled board with the smallest square size.

Example:

Input (`sample.txt`):

```
...#
...#
...#
...#

....
....
....
####

.###
...#
....
....

....
..##
.##.
....
```

Output:

```
ABBBB.
ACCCEE
AFFCEE
A.FFGG
HHHDDG
.HDD.G
```

---

## Error Handling

If the file is missing, incorrectly formatted, or a tetromino is invalid (e.g., not exactly 4 connected `#` cells), the program prints:

```
ERROR
```

---

## Project Structure

```
.
├── board.go        # Board representation and basic operations
├── candidates.go   # Candidate generation for placements
├── shape_ops.go    # Shape parsing, normalization, rotations
├── solver.go       # Backtracking DFS with pruning
├── types.go        # Data structures and shared types
├── main.go         # Entry point (reads file, runs solver)
└── examples/       # Example input files
```

---

## Testing

It is recommended to create `.txt` input files in `examples/` and run them to validate behavior:

```bash
go run . examples/goodexample01.txt
go run . examples/badexample.txt
```

---

## License

This project is developed for learning purposes.  