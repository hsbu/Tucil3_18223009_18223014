package algorithm

import (
	"time"
	"tucil3/internal/board"
	"tucil3/internal/heuristic"
)

type Algorithm interface {
	Search(b *board.Board, start board.State, h heuristic.Func) Result
}

type Node struct {
	State board.State
	G int
	H int
	F int
	Parent *Node
	Move board.Direction
}

type Result struct {
	Moves []board.Direction
	TotalCost int
	Iterations int
	Duration time.Duration
	Snapshots []board.State
	Found bool
}

func allVisited(s board.State, b *board.Board) bool {
	for i := range b.Checkpoints {
		if !s.HasVisited(i) {
			return false
		}
	}
	return true
}

func reconstructPath(goal *Node) ([]board.Direction, []board.State) {
	var moves []board.Direction
	var states []board.State
	for n := goal; n.Parent != nil; n = n.Parent {
		moves = append([]board.Direction{n.Move}, moves...)
		states = append([]board.State{n.State}, states...)
	}
	return moves, states
}
