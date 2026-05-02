package board

func Slide(b *Board, s State, dir Direction) (State, int, bool) {
	dRow, dCol := dirDelta(dir)
	row, col := s.Pos.Row, s.Pos.Col
	cost := 0
	newVisited := s.Visited

	for {
		nextRow := row + dRow
		nextCol := col + dCol

		// If out of bounds, game over
		if nextRow < 0 || nextRow >= b.Rows || nextCol < 0 || nextCol >= b.Cols {
			return s, 0, false
		}

		next := b.Grid[nextRow][nextCol]

		// Stop before obstacle
		if next == TileObstacle {
			break
		}

		// If lava, game over
		if next == TileLava {
			return s, 0, false
		}

		// Checkpoint ordering constraint
		if IsCheckpoint(next) {
			idx := CheckpointIndex(next)
			if (newVisited & (1 << uint(idx))) == 0 {
				if idx > 0 && (newVisited&(1<<uint(idx-1))) == 0 {
					return s, 0, false
				}
				newVisited |= 1 << uint(idx)
			}
		}

		row, col = nextRow, nextCol
		cost += b.Cost[row][col]

		// Stop at goal
		if next == TileGoal {
			break
		}
	}

	if row == s.Pos.Row && col == s.Pos.Col {
		return s, 0, false
	}

	return State{Pos: Position{row, col}, Visited: newVisited}, cost, true
}

func dirDelta(dir Direction) (int, int) {
	switch dir {
	case Up:
		return -1, 0
	case Down:
		return 1, 0
	case Left:
		return 0, -1
	case Right:
		return 0, 1
	}
	return 0, 0
}
