package algorithm

import (
	"math"
	"time"

	"tucil3/internal/board"
	"tucil3/internal/heuristic"
)

type IDAStar struct{}

func (IDAStar) Search(b *board.Board, start board.State, h heuristic.Func) Result {
	t0 := time.Now()
	iterations := 0

	hVal := 0
	if h != nil {
		hVal = h(start, b)
	}
	root := &Node{State: start, G: 0, H: hVal, F: hVal}
	bound := root.F

	for {
		path := map[board.State]bool{start: true}
		nextBound, found := idaSearch(b, root, bound, h, path, &iterations)
		if found != nil {
			moves, states := reconstructPath(found)
			return Result{
				Moves:      moves,
				TotalCost:  found.G,
				Iterations: iterations,
				Duration:   time.Since(t0),
				Snapshots:  states,
				Found:      true,
			}
		}
		if nextBound == math.MaxInt {
			return Result{Iterations: iterations, Duration: time.Since(t0)}
		}
		bound = nextBound
	}
}

func idaSearch(
	b *board.Board,
	cur *Node,
	bound int,
	h heuristic.Func,
	path map[board.State]bool,
	iterations *int,
) (int, *Node) {
	(*iterations)++
	if cur.F > bound {
		return cur.F, nil
	}
	if cur.State.Pos == b.Goal && allVisited(cur.State, b) {
		return cur.F, cur
	}

	minNextBound := math.MaxInt
	for _, dir := range []board.Direction{board.Up, board.Down, board.Left, board.Right} {
		ns, cost, valid := board.Slide(b, cur.State, dir)
		if !valid || path[ns] {
			continue
		}

		g := cur.G + cost
		hVal := 0
		if h != nil {
			hVal = h(ns, b)
		}
		next := &Node{State: ns, G: g, H: hVal, F: g + hVal, Parent: cur, Move: dir}
		path[ns] = true
		nextBound, found := idaSearch(b, next, bound, h, path, iterations)
		delete(path, ns)

		if found != nil {
			return nextBound, found
		}
		if nextBound < minNextBound {
			minNextBound = nextBound
		}
	}

	return minNextBound, nil
}
