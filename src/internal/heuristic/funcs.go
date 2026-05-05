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

func remainingChain(s board.State, b *board.Board) []board.Position {
	chain := []board.Position{s.Pos}
	for i, cp := range b.Checkpoints {
		if !s.HasVisited(i) {
			chain = append(chain, cp)
		}
	}
	chain = append(chain, b.Goal)
	return chain
}

func chainManhattan(s board.State, b *board.Board) int {
	chain := remainingChain(s, b)
	sum := 0
	for i := 1; i < len(chain); i++ {
		sum += manhattan(chain[i-1], chain[i])
	}
	return sum
}

// Heuristics

// H1: Manhattan dari posisi saat ini ke checkpoint berikutnya atau goal
var H1 Func = func(s board.State, b *board.Board) int {
	return manhattan(s.Pos, nextTarget(s, b))
}

func unvisitedCount(s board.State, b *board.Board) int {
	count := 0
	for i := range b.Checkpoints {
		if !s.HasVisited(i) {
			count++
		}
	}
	return count
}

// H2: Jumlah Manhattan dari posisi saat ini ke setiap checkpoint yang belum dikunjungi, lalu ke goal
var H2 Func = func(s board.State, b *board.Board) int {
	return chainManhattan(s, b)
}

// H3: H2 + jumlah checkpoint yang belum dikunjungi sebagai penalti
var H3 Func = func(s board.State, b *board.Board) int {
	return H2(s, b) + unvisitedCount(s, b)
}
