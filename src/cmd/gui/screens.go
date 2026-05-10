package main

import (
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tucil3/internal/algorithm"
	"tucil3/internal/board"
	"tucil3/internal/heuristic"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	pad    = 20
	btnW   = 200
	btnH   = 40
	inputW = 700
	inputH = 28
)

// helpers

func drawRect(screen *ebiten.Image, x, y, w, h float64, c color.RGBA) {
	ebitenutil.DrawRect(screen, x, y, w, h, c)
}

func drawButton(screen *ebiten.Image, label string, x, y, w, h int, hover bool) {
	bg := color.RGBA{60, 60, 100, 255}
	if hover {
		bg = color.RGBA{90, 90, 160, 255}
	}
	drawRect(screen, float64(x), float64(y), float64(w), float64(h), bg)
	ebitenutil.DebugPrintAt(screen, label, x+10, y+10)
}

func inRect(x, y, w, h int) bool {
	cx, cy := ebiten.CursorPosition()
	return cx >= x && cx <= x+w && cy >= y && cy <= y+h
}

type uiRect struct{ x, y, w, h int }

func (r uiRect) contains() bool {
	return inRect(r.x, r.y, r.w, r.h)
}

func centerX(w int) int { return (screenW - w) / 2 }
func centerY(h int) int { return (screenH - h) / 2 }

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func textWidth(text string) int {
	return len([]rune(text)) * 7
}

func centerTextX(text string) int {
	return centerX(textWidth(text))
}

func formatMilliseconds(d time.Duration) string {
	if d < time.Microsecond {
		return "<0.001 ms"
	}
	return fmt.Sprintf("%.3f ms", float64(d.Nanoseconds())/1_000_000)
}

func fileSelectLayout() (uiRect, uiRect) {
	input := uiRect{x: centerX(inputW), y: 230, w: inputW, h: inputH}
	load := uiRect{x: centerX(btnW), y: 270, w: btnW, h: btnH}
	return input, load
}

func algoButtonsLayout() (uiRect, uiRect, uiRect, uiRect) {
	gap := 20
	w := 150
	totalW := w*4 + gap*3
	startX := centerX(totalW)
	y := 270
	ucs := uiRect{x: startX, y: y, w: w, h: 50}
	gbfs := uiRect{x: startX + w + gap, y: y, w: w, h: 50}
	astar := uiRect{x: startX + (w+gap)*2, y: y, w: w, h: 50}
	idastar := uiRect{x: startX + (w+gap)*3, y: y, w: w, h: 50}
	return ucs, gbfs, astar, idastar
}

func heuristicButtonsLayout() []uiRect {
	gap := 16
	w := 110
	totalW := w*len(heuristic.All) + gap*(len(heuristic.All)-1)
	startX := centerX(totalW)
	y := 270
	rects := make([]uiRect, len(heuristic.All))
	for i := range rects {
		rects[i] = uiRect{x: startX + i*(w+gap), y: y, w: w, h: 50}
	}
	return rects
}

func backButtonLayout() uiRect {
	return uiRect{x: pad, y: screenH - 50, w: 100, h: 30}
}

func resultPlaybackLayout() uiRect {
	return uiRect{x: centerX(btnW) + 180, y: 460, w: btnW, h: btnH}
}

func resultBackLayout() uiRect {
	return uiRect{x: centerX(btnW) - 180, y: 460, w: btnW, h: btnH}
}

func resultSaveLayout() uiRect {
	return uiRect{x: centerX(btnW), y: 460, w: btnW, h: btnH}
}

func playbackLayout() (uiRect, uiRect, uiRect, int, int) {
	buttonsY := screenH - 50
	hudY := screenH - 75
	hintY := screenH - 42
	prev := uiRect{x: pad, y: buttonsY, w: 80, h: 30}
	next := uiRect{x: pad + 90, y: buttonsY, w: 80, h: 30}
	play := uiRect{x: pad + 180, y: buttonsY, w: 100, h: 30}
	return prev, next, play, hudY, hintY
}

func fittedTileSize(b *board.Board) int {
	maxW := screenW - pad*2
	maxH := screenH - 130
	size := minInt(tileSize, minInt(maxW/b.Cols, maxH/b.Rows))
	if size < 8 {
		return 8
	}
	return size
}

// Screen

type FileSelectScreen struct {
	path   string
	errMsg string
}

func newFileSelectScreen() *FileSelectScreen { return &FileSelectScreen{} }

func (s *FileSelectScreen) Update() (Screen, error) {
	for _, r := range ebiten.AppendInputChars(nil) {
		s.path += string(r)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len([]rune(s.path)) > 0 {
		r := []rune(s.path)
		s.path = string(r[:len(r)-1])
	}

	load := inpututil.IsKeyJustPressed(ebiten.KeyEnter)
	if !load && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_, loadRect := fileSelectLayout()
		load = loadRect.contains()
	}
	if load {
		trimmed := strings.TrimSpace(s.path)
		if trimmed != "" && !filepath.IsAbs(trimmed) {
			trimmed = filepath.Join("..", "test", trimmed)
		}
		f, err := os.Open(trimmed)
		if err != nil {
			s.errMsg = "Cannot open file: " + err.Error()
			return s, nil
		}
		defer f.Close()
		b, err := board.Parse(f)
		if err != nil {
			s.errMsg = "Invalid board: " + err.Error()
			return s, nil
		}
		return newAlgoSelectScreen(b), nil
	}
	return s, nil
}

func (s *FileSelectScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 30, 255})
	title := "ICE SLIDING PUZZLE SOLVER"
	label := "Input file path (.txt):"
	footer := "Press Enter or click Load to continue"
	ebitenutil.DebugPrintAt(screen, title, centerTextX(title), 80)
	ebitenutil.DebugPrintAt(screen, label, centerTextX(label), 210)
	inputRect, loadRect := fileSelectLayout()
	drawRect(screen, float64(inputRect.x), float64(inputRect.y), float64(inputRect.w), float64(inputRect.h), color.RGBA{35, 35, 60, 255})
	displayPath := s.path
	maxChars := (inputRect.w - 14) / 7
	if len([]rune(displayPath)) > maxChars {
		runes := []rune(displayPath)
		displayPath = "..." + string(runes[len(runes)-maxChars+3:])
	}
	ebitenutil.DebugPrintAt(screen, displayPath+"_", inputRect.x+5, inputRect.y+6)
	drawButton(screen, "Load", loadRect.x, loadRect.y, loadRect.w, loadRect.h, loadRect.contains())
	if s.errMsg != "" {
		ebitenutil.DebugPrintAt(screen, s.errMsg, centerTextX(s.errMsg), 330)
	}
	ebitenutil.DebugPrintAt(screen, footer, centerTextX(footer), 400)
}

type AlgoSelectScreen struct {
	b *board.Board
}

func newAlgoSelectScreen(b *board.Board) *AlgoSelectScreen {
	return &AlgoSelectScreen{b: b}
}

func (s *AlgoSelectScreen) Update() (Screen, error) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return s, nil
	}
	if backButtonLayout().contains() {
		return newFileSelectScreen(), nil
	}
	ucs, gbfs, astar, idastar := algoButtonsLayout()
	if ucs.contains() {
		return newSolvingScreen(s.b, algorithm.UCS{}, nil), nil
	}
	if gbfs.contains() {
		return newHeuristicSelectScreen(s.b, algorithm.GBFS{}, "GBFS"), nil
	}
	if astar.contains() {
		return newHeuristicSelectScreen(s.b, algorithm.AStar{}, "A*"), nil
	}
	if idastar.contains() {
		return newHeuristicSelectScreen(s.b, algorithm.IDAStar{}, "IDA*"), nil
	}
	return s, nil
}

func (s *AlgoSelectScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 30, 255})
	ebitenutil.DebugPrintAt(screen, "Choose Algorithm", 352, 180)
	ucs, gbfs, astar, idastar := algoButtonsLayout()
	drawButton(screen, "UCS", ucs.x, ucs.y, ucs.w, ucs.h, ucs.contains())
	drawButton(screen, "GBFS", gbfs.x, gbfs.y, gbfs.w, gbfs.h, gbfs.contains())
	drawButton(screen, "A*", astar.x, astar.y, astar.w, astar.h, astar.contains())
	drawButton(screen, "IDA*", idastar.x, idastar.y, idastar.w, idastar.h, idastar.contains())
	ebitenutil.DebugPrintAt(screen, "GBFS, A*, and IDA* require a heuristic", 260, 340)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Board: %dx%d  Checkpoints: %d", s.b.Rows, s.b.Cols, len(s.b.Checkpoints)), 286, 420)
	back := backButtonLayout()
	drawButton(screen, "< Back", back.x, back.y, back.w, back.h, back.contains())
}

type HeuristicSelectScreen struct {
	b        *board.Board
	algo     algorithm.Algorithm
	algoName string
}

func newHeuristicSelectScreen(b *board.Board, algo algorithm.Algorithm, algoName string) *HeuristicSelectScreen {
	return &HeuristicSelectScreen{b: b, algo: algo, algoName: algoName}
}

func (s *HeuristicSelectScreen) Update() (Screen, error) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return s, nil
	}
	if backButtonLayout().contains() {
		return newAlgoSelectScreen(s.b), nil
	}
	for i, rect := range heuristicButtonsLayout() {
		if rect.contains() {
			return newSolvingScreen(s.b, s.algo, heuristic.All[i].Fn), nil
		}
	}
	return s, nil
}

func (s *HeuristicSelectScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 30, 255})
	title := fmt.Sprintf("Choose Heuristic for %s", s.algoName)
	ebitenutil.DebugPrintAt(screen, title, centerTextX(title), 180)
	for i, rect := range heuristicButtonsLayout() {
		h := heuristic.All[i]
		label := h.Name
		if h.Admissible {
			label += " *"
		}
		drawButton(screen, label, rect.x, rect.y, rect.w, rect.h, rect.contains())
	}
	ebitenutil.DebugPrintAt(screen, "* admissible", centerTextX("* admissible"), 340)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Board: %dx%d  Checkpoints: %d", s.b.Rows, s.b.Cols, len(s.b.Checkpoints)), 286, 420)
	back := backButtonLayout()
	drawButton(screen, "< Back", back.x, back.y, back.w, back.h, back.contains())
}

type SolvingScreen struct {
	b    *board.Board
	done chan algorithm.Result
}

func newSolvingScreen(b *board.Board, algo algorithm.Algorithm, h heuristic.Func) *SolvingScreen {
	s := &SolvingScreen{b: b, done: make(chan algorithm.Result, 1)}
	start := board.State{Pos: b.Start}
	go func() { s.done <- algo.Search(b, start, h) }()
	return s
}

func (s *SolvingScreen) Update() (Screen, error) {
	select {
	case r := <-s.done:
		return newResultScreen(s.b, r), nil
	default:
		return s, nil
	}
}

func (s *SolvingScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 30, 255})
	ebitenutil.DebugPrintAt(screen, "Solving...", 355, 280)
	ebitenutil.DebugPrintAt(screen, "Please wait", 345, 300)
}

type ResultScreen struct {
	b          *board.Board
	result     algorithm.Result
	saveStatus string
}

func newResultScreen(b *board.Board, result algorithm.Result) *ResultScreen {
	return &ResultScreen{b: b, result: result}
}

func (s *ResultScreen) Update() (Screen, error) {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if resultBackLayout().contains() {
			return newAlgoSelectScreen(s.b), nil
		}
		if s.result.Found && resultSaveLayout().contains() {
			if err := saveSolutionFile("solution.txt", s.b, s.result); err != nil {
				s.saveStatus = "Save failed: " + err.Error()
			} else {
				s.saveStatus = "Saved to solution.txt"
			}
			return s, nil
		}
		if s.result.Found && resultPlaybackLayout().contains() {
			return newPlaybackScreen(s.b, s.result), nil
		}
	}
	return s, nil
}

func (s *ResultScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 30, 255})
	if !s.result.Found {
		ebitenutil.DebugPrintAt(screen, "No solution found.", 310, 280)
		return
	}
	var sb strings.Builder
	for _, m := range s.result.Moves {
		sb.WriteString(m.String())
	}
	moves := sb.String()
	ebitenutil.DebugPrintAt(screen, "SOLUTION FOUND", 320, 100)
	ebitenutil.DebugPrintAt(screen, "Moves : "+moves, 50, 160)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cost : %d", s.result.TotalCost), 50, 200)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Iterations : %d", s.result.Iterations), 50, 240)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Time : %s", formatMilliseconds(s.result.Duration)), 50, 280)
	backRect := resultBackLayout()
	saveRect := resultSaveLayout()
	playRect := resultPlaybackLayout()
	drawButton(screen, "< Rechoose", backRect.x, backRect.y, backRect.w, backRect.h, backRect.contains())
	drawButton(screen, "Save TXT", saveRect.x, saveRect.y, saveRect.w, saveRect.h, saveRect.contains())
	drawButton(screen, "Playback >>", playRect.x, playRect.y, playRect.w, playRect.h, playRect.contains())
	if s.saveStatus != "" {
		ebitenutil.DebugPrintAt(screen, s.saveStatus, centerTextX(s.saveStatus), 520)
	}
}

func saveSolutionFile(path string, b *board.Board, result algorithm.Result) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintf(f, "Solusi : %s\n", movesString(result.Moves))
	fmt.Fprintf(f, "Cost   : %d\n", result.TotalCost)
	fmt.Fprintf(f, "Waktu  : %s\n", formatMilliseconds(result.Duration))
	fmt.Fprintf(f, "Iterasi: %d\n", result.Iterations)
	fmt.Fprintln(f)
	fmt.Fprintln(f, "Initial")
	writeBoardText(f, b, board.State{Pos: b.Start})
	for i, move := range result.Moves {
		fmt.Fprintf(f, "Step %d : %v\n", i+1, move)
		writeBoardText(f, b, result.Snapshots[i])
	}
	return nil
}

func movesString(moves []board.Direction) string {
	var sb strings.Builder
	for _, m := range moves {
		sb.WriteString(m.String())
	}
	return sb.String()
}

func writeBoardText(w io.Writer, b *board.Board, s board.State) {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if s.Pos == (board.Position{Row: r, Col: c}) {
				fmt.Fprint(w, "Z")
				continue
			}
			t := b.Grid[r][c]
			if board.IsCheckpoint(t) && s.HasVisited(board.CheckpointIndex(t)) {
				fmt.Fprint(w, "*")
			} else {
				fmt.Fprintf(w, "%c", t)
			}
		}
		fmt.Fprintln(w)
	}
}

const tileSize = 36

type PlaybackScreen struct {
	b              *board.Board
	result         algorithm.Result
	allStates      []board.State
	step           int
	playing        bool
	tick           int
	stepsPerSecond float64
}

func newPlaybackScreen(b *board.Board, result algorithm.Result) *PlaybackScreen {
	initial := board.State{Pos: b.Start}
	all := make([]board.State, 0, len(result.Snapshots)+1)
	all = append(all, initial)
	all = append(all, result.Snapshots...)
	return &PlaybackScreen{
		b: b, result: result,
		allStates: all, stepsPerSecond: 3,
	}
}

func (s *PlaybackScreen) Update() (Screen, error) {
	// Keyboard navigation
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) && s.step < len(s.result.Moves) {
		s.step++
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) && s.step > 0 {
		s.step--
	}
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		s.togglePlayback()
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) && s.stepsPerSecond < 12 {
		s.stepsPerSecond += 0.5
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) && s.stepsPerSecond > 0.5 {
		s.stepsPerSecond -= 0.5
	}

	if s.playing {
		s.tick++
		framesPerStep := int(60.0 / s.stepsPerSecond)
		if framesPerStep < 1 {
			framesPerStep = 1
		}
		if s.tick >= framesPerStep {
			s.tick = 0
			if s.step < len(s.result.Moves) {
				s.step++
			} else {
				s.playing = false
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		prevRect, nextRect, playRect, _, _ := playbackLayout()
		if backButtonLayout().contains() {
			return newResultScreen(s.b, s.result), nil
		}
		if prevRect.contains() && s.step > 0 {
			s.step--
		}
		if nextRect.contains() && s.step < len(s.result.Moves) {
			s.step++
		}
		if playRect.contains() {
			s.togglePlayback()
		}
	}

	return s, nil
}

func (s *PlaybackScreen) togglePlayback() {
	if !s.playing && s.step >= len(s.result.Moves) {
		s.step = 0
		s.tick = 0
	}
	s.playing = !s.playing
}

func (s *PlaybackScreen) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{15, 15, 30, 255})

	cur := s.allStates[s.step]
	tile := fittedTileSize(s.b)
	gridW := s.b.Cols * tile
	gridH := s.b.Rows * tile
	offX := maxInt(pad, centerX(gridW))
	offY := maxInt(pad, centerY(gridH)-40)

	tileColors := map[board.Tile]color.RGBA{
		board.TilePath:     {80, 160, 200, 255},
		board.TileObstacle: {200, 200, 200, 255},
		board.TileLava:     {220, 70, 30, 255},
		board.TileGoal:     {50, 200, 80, 255},
		board.TileStart:    {80, 160, 200, 255},
	}

	for r := 0; r < s.b.Rows; r++ {
		for c := 0; c < s.b.Cols; c++ {
			x := float64(offX + c*tile)
			y := float64(offY + r*tile)
			t := s.b.Grid[r][c]

			var col color.RGBA
			if bc, ok := tileColors[t]; ok {
				col = bc
			} else if board.IsCheckpoint(t) {
				idx := board.CheckpointIndex(t)
				if cur.HasVisited(idx) {
					col = color.RGBA{100, 100, 100, 255}
				} else {
					col = color.RGBA{220, 190, 50, 255}
				}
			}
			drawRect(screen, x+1, y+1, float64(tile-2), float64(tile-2), col)

			if tile >= 18 && board.IsCheckpoint(t) {
				ebitenutil.DebugPrintAt(screen, string(t), int(x)+tile/3, int(y)+tile/3)
			}
			if tile >= 18 && t == board.TileGoal {
				ebitenutil.DebugPrintAt(screen, "O", int(x)+tile/3, int(y)+tile/3)
			}
		}
	}

	// Draw
	ax := float64(offX + cur.Pos.Col*tile)
	ay := float64(offY + cur.Pos.Row*tile)
	inset := maxInt(2, tile/7)
	drawRect(screen, ax+float64(inset), ay+float64(inset), float64(tile-inset*2), float64(tile-inset*2), color.RGBA{30, 80, 220, 255})
	if tile >= 18 {
		ebitenutil.DebugPrintAt(screen, "Z", int(ax)+tile/3, int(ay)+tile/3)
	}

	// HUD
	stepLabel := "Initial"
	if s.step > 0 {
		stepLabel = fmt.Sprintf("Step %d/%d  Move: %v", s.step, len(s.result.Moves), s.result.Moves[s.step-1])
	}
	prevRect, nextRect, playRect, hudY, hintY := playbackLayout()
	ebitenutil.DebugPrintAt(screen, stepLabel, pad, hudY)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Cost: %d  Speed: %.1f steps/s", s.result.TotalCost, s.stepsPerSecond), 400, hudY)

	// Buttons
	drawButton(screen, "< Prev", prevRect.x, prevRect.y, prevRect.w, prevRect.h, prevRect.contains())
	drawButton(screen, "Next >", nextRect.x, nextRect.y, nextRect.w, nextRect.h, nextRect.contains())
	playLabel := "Play"
	if s.playing {
		playLabel = "Pause"
	}
	drawButton(screen, playLabel, playRect.x, playRect.y, playRect.w, playRect.h, playRect.contains())
	back := backButtonLayout()
	drawButton(screen, "< Back", back.x, back.y, back.w, back.h, back.contains())
	ebitenutil.DebugPrintAt(screen, "[← →] step  [Space] play  [+/-] speed", 320, hintY)
}
