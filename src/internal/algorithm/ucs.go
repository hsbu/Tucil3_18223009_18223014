package algorithm

import (
	"time"
	"tucil3/internal/board"
	"tucil3/internal/heuristic"
)

type UCS struct{}

func (UCS) Search(b *board.Board, start board.State, _ heuristic.Func) Result {
	t0 := time.Now()
	pq := NewPriorityQueue()
	best := map[board.State]int{}
	iterations := 0

	pq.Push(&Node{State: start, G: 0, F: 0})

	for pq.Len() > 0 {
		cur := pq.Pop()
		iterations++

		if prev, seen := best[cur.State]; seen && prev <= cur.G {
			continue
		}
		best[cur.State] = cur.G

		if cur.State.Pos == b.Goal && allVisited(cur.State, b) {
			moves, states := reconstructPath(cur)
			return Result{
				Moves: moves,
				TotalCost: cur.G,
				Iterations: iterations,
				Duration: time.Since(t0),
				Snapshots: states,
				Found: true,
			}
		}

		for _, dir := range []board.Direction{board.Up, board.Down, board.Left, board.Right} {
			ns, cost, valid := board.Slide(b, cur.State, dir)
			if !valid {
				continue
			}
			g := cur.G + cost
			if prev, seen := best[ns]; seen && prev <= g {
				continue
			}
			pq.Push(&Node{State: ns, G: g, F: g, Parent: cur, Move: dir})
		}
	}

	return Result{Iterations: iterations, Duration: time.Since(t0)}
}
