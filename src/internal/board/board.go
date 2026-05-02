package board

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Board struct {
	Grid [][]Tile
	Cost [][]int
	Rows, Cols int
	Start, Goal Position
	Checkpoints []Position
}

func Parse(r io.Reader) (*Board, error) {
	scanner := bufio.NewScanner(r)

	if !scanner.Scan() {
		return nil, fmt.Errorf("missing dimensions line")
	}
	parts := strings.Fields(scanner.Text())
	if len(parts) != 2 {
		return nil, fmt.Errorf("expected 2 dimensions, got %d", len(parts))
	}
	rows, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid rows: %w", err)
	}
	cols, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid cols: %w", err)
	}

	b := &Board{
		Grid: make([][]Tile, rows),
		Cost: make([][]int, rows),
		Rows: rows,
		Cols: cols,
	}

	cpMap := map[int]Position{}

	for i := 0; i < rows; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing grid row %d", i)
		}
		line := scanner.Text()
		if len(line) != cols {
			return nil, fmt.Errorf("row %d: expected %d cols, got %d", i, cols, len(line))
		}
		b.Grid[i] = make([]Tile, cols)
		for j, ch := range []byte(line) {
			t := Tile(ch)
			b.Grid[i][j] = t
			switch {
			case t == TileStart:
				b.Start = Position{i, j}
			case t == TileGoal:
				b.Goal = Position{i, j}
			case IsCheckpoint(t):
				cpMap[CheckpointIndex(t)] = Position{i, j}
			}
		}
	}

	if len(cpMap) > 0 {
		maxIdx := 0
		for idx := range cpMap {
			if idx > maxIdx {
				maxIdx = idx
			}
		}
		b.Checkpoints = make([]Position, maxIdx+1)
		for idx, pos := range cpMap {
			b.Checkpoints[idx] = pos
		}
	}

	for i := 0; i < rows; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing cost row %d", i)
		}
		fields := strings.Fields(scanner.Text())
		if len(fields) != cols {
			return nil, fmt.Errorf("cost row %d: expected %d values, got %d", i, cols, len(fields))
		}
		b.Cost[i] = make([]int, cols)
		for j, f := range fields {
			c, err := strconv.Atoi(f)
			if err != nil {
				return nil, fmt.Errorf("invalid cost at (%d,%d): %w", i, j, err)
			}
			b.Cost[i][j] = c
		}
	}

	return b, nil
}
