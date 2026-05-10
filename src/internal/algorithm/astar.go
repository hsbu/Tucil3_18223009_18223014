package algorithm

import (
	"time"
	"tucil3/internal/board"
	"tucil3/internal/heuristic"
)

type AStar struct{}

func (AStar) Search(b *board.Board, start board.State, h heuristic.Func) Result {
	t0 := time.Now()
	pq := NewPriorityQueue()
	best := map[board.State]int{}
	iterations := 0

	hVal := 0
	if h != nil {
		hVal = h(start, b)
	}
	pq.Push(&Node{State: start, G: 0, H: hVal, F: hVal})

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
				Moves:      moves,
				TotalCost:  cur.G,
				Iterations: iterations,
				Duration:   time.Since(t0),
				Snapshots:  states,
				Found:      true,
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
			hVal := 0
			if h != nil {
				hVal = h(ns, b)
			}
			pq.Push(&Node{State: ns, G: g, H: hVal, F: g + hVal, Parent: cur, Move: dir})
		}
	}

	return Result{Iterations: iterations, Duration: time.Since(t0)}
}
