package board

type Position struct{ Row, Col int }

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "U"
	case Down:
		return "D"
	case Left:
		return "L"
	case Right:
		return "R"
	}
	return "?"
}

type Tile byte

const (
	TilePath Tile = '*'
	TileObstacle Tile = 'X'
	TileLava Tile = 'L'
	TileStart Tile = 'Z'
	TileGoal Tile = 'O'
)

func IsCheckpoint(t Tile) bool { return t >= '0' && t <= '9'}
func CheckpointIndex(t Tile) int {return int(t - '0')}

type State struct {
	Pos Position
	Visited uint16
}

func (s State) HasVisited(i int) bool { return s.Visited&(1<<uint(i)) != 0 }
func (s State) WithVisited(i int) State {
	s.Visited |= 1 << uint(i)
	return s
}
