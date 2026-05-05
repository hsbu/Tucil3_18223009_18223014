package heuristic

import "tucil3/internal/board"

type Func func(s board.State, b *board.Board) int

var All = []struct {
	Name       string
	Fn         Func
	Admissible bool
}{
	{"H1", H1, true},
	{"H2", H2, true},
	{"H3", H3, false},
}
