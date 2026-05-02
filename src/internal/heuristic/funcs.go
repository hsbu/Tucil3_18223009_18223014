package heuristic

import (
	"tucil3/internal/board"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func manhattan(a, b board.Position) int {
	return abs(a.Row-b.Row) + abs(a.Col-b.Col)
}

func nextTarget(s board.State, b *board.Board) board.Position {
	for i, cp := range b.Checkpoints {
		if !s.HasVisited(i) {
			return cp
		}
	}
	return b.Goal
}

// Heuristics

// Manhattan distance from current position to the next unvisited checkpoint, or to the goal if all checkpoints have been visited.
var H1 Func = func(s board.State, b *board.Board) int {
	return manhattan(s.Pos, nextTarget(s, b))
}
