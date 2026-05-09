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

func minTileCost(b *board.Board) int {
	minCost := 0
	found := false
	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			tile := b.Grid[row][col]
			if tile == board.TileObstacle || tile == board.TileLava {
				continue
			}
			cost := b.Cost[row][col]
			if !found || cost < minCost {
				minCost = cost
				found = true
			}
		}
	}
	return minCost
}

func slideLengthEstimate(b *board.Board, pos board.Position, dir board.Direction) (int, bool) {
	dRow, dCol := dirDelta(dir)
	row, col := pos.Row, pos.Col
	length := 0

	for {
		nextRow := row + dRow
		nextCol := col + dCol
		if nextRow < 0 || nextRow >= b.Rows || nextCol < 0 || nextCol >= b.Cols {
			return 0, false
		}

		next := b.Grid[nextRow][nextCol]
		if next == board.TileObstacle {
			break
		}
		if next == board.TileLava {
			return 0, false
		}

		row, col = nextRow, nextCol
		length++
		if next == board.TileGoal {
			break
		}
	}

	return length, length > 0
}

func dirDelta(dir board.Direction) (int, int) {
	switch dir {
	case board.Up:
		return -1, 0
	case board.Down:
		return 1, 0
	case board.Left:
		return 0, -1
	case board.Right:
		return 0, 1
	}
	return 0, 0
}

func averageSlideLength(b *board.Board) (total, count int) {
	for row := 0; row < b.Rows; row++ {
		for col := 0; col < b.Cols; col++ {
			tile := b.Grid[row][col]
			if tile == board.TileObstacle || tile == board.TileLava {
				continue
			}
			pos := board.Position{Row: row, Col: col}
			for _, dir := range []board.Direction{board.Up, board.Down, board.Left, board.Right} {
				if length, ok := slideLengthEstimate(b, pos, dir); ok {
					total += length
					count++
				}
			}
		}
	}
	return total, count
}

func estimatedSlideCost(s board.State, b *board.Board) int {
	distance := H2(s, b)
	if distance == 0 {
		return 0
	}
	total, count := averageSlideLength(b)
	if total == 0 || count == 0 {
		return 0
	}
	return (distance*count + total - 1) / total
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

// H4: H2 dikalikan biaya tile minimum sebagai lower bound biaya
var H4 Func = func(s board.State, b *board.Board) int {
	return H2(s, b) * minTileCost(b)
}

// H5: H4 + estimasi biaya slide berdasarkan H2 / rata-rata panjang slide
var H5 Func = func(s board.State, b *board.Board) int {
	return H4(s, b) + estimatedSlideCost(s, b)
}
