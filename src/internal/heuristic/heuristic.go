package heuristic

import "tucil3/internal/board"

type Func func(s board.State, b *board.Board) int

var All = []struct {
	Name string
	Fn         Func
	Admissible bool
}{
	{"H1", H1, true},
	// Add another heuristics here
}
