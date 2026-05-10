package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"tucil3/internal/algorithm"
	"tucil3/internal/board"
	"tucil3/internal/heuristic"

	"golang.org/x/term"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// ── File input ──
	fmt.Print(">> Masukan file input :\n")
	path, _ := reader.ReadString('\n')
	path = strings.TrimSpace(path)

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error membuka file:", err)
		os.Exit(1)
	}
	defer f.Close()

	b, err := board.Parse(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error membaca board:", err)
		os.Exit(1)
	}

	// ── Algorithm selection ──
	fmt.Print(">> Algoritma apa yang anda pilih? (UCS/GBFS/A*/IDA*) :\n")
	algoStr, _ := reader.ReadString('\n')
	algoStr = strings.TrimSpace(strings.ToUpper(algoStr))

	var algo algorithm.Algorithm

	switch algoStr {
	case "UCS":
		algo = algorithm.UCS{}
	case "GBFS":
		algo = algorithm.GBFS{}
	case "A*", "A", "ASTAR":
		algo = algorithm.AStar{}
	case "IDA*", "IDA", "IDASTAR":
		algo = algorithm.IDAStar{}
	default:
		fmt.Fprintln(os.Stderr, "Algoritma tidak dikenal. Pilihan: UCS, GBFS, A*, IDA*")
		os.Exit(1)
	}

	// ── Heuristic selection ──
	var h heuristic.Func

	switch algo.(type) {
	case algorithm.UCS:
		h = nil
	default:
		fmt.Print(">> Heuristic apa yang anda pilih? (H1/H2/H3/H4/H5) :\n")
		hStr, _ := reader.ReadString('\n')
		hStr = strings.TrimSpace(strings.ToUpper(hStr))

		switch hStr {
		case "H1":
			h = heuristic.H1
		case "H2", "":
			h = heuristic.H2
		case "H3":
			h = heuristic.H3
		case "H4":
			h = heuristic.H4
		case "H5":
			h = heuristic.H5
		default:
			fmt.Println("Heuristic tidak dikenal, menggunakan H2.")
			h = heuristic.H2
		}
	}

	// ── Run search ──
	start := board.State{Pos: b.Start}
	result := algo.Search(b, start, h)

	if !result.Found {
		fmt.Println("Tidak ditemukan solusi.")
		return
	}

	// ── Output solution ──
	var movesStr strings.Builder
	for _, m := range result.Moves {
		movesStr.WriteString(m.String())
	}
	fmt.Printf("\nSolusi Yang Ditemukan : %s\n", movesStr.String())
	fmt.Printf("Cost dari Solusi      : %d\n", result.TotalCost)

	// ── Execution stats ──
	fmt.Printf("\n>> Waktu eksekusi: %s\n", formatMilliseconds(result.Duration))
	fmt.Printf(">> Banyak iterasi yang dilakukan: %d iterasi\n", result.Iterations)

	// ── Playback prompt ──
	fmt.Print("\n>> Apakah Anda ingin melakukan playback? (Ya/Tidak) :\n")
	playbackAns, _ := reader.ReadString('\n')
	playbackAns = strings.TrimSpace(strings.ToLower(playbackAns))

	if playbackAns == "ya" || playbackAns == "y" {
		allStates := make([]board.State, 0, len(result.Snapshots)+1)
		allStates = append(allStates, start)
		allStates = append(allStates, result.Snapshots...)
		runPlayback(reader, b, result, allStates)
	}

	// ── Save solution prompt ──
	fmt.Print("\n>> Apakah Anda ingin menyimpan solusi? (Ya/Tidak) :\n")
	saveAns, _ := reader.ReadString('\n')
	saveAns = strings.TrimSpace(strings.ToLower(saveAns))

	if saveAns == "ya" || saveAns == "y" {
		fmt.Print(">> Masukan path file output :\n")
		outPath, _ := reader.ReadString('\n')
		outPath = strings.TrimSpace(outPath)
		if outPath == "" {
			outPath = "solusi.txt"
		}

		outFile, err := os.Create(outPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error menyimpan file:", err)
		} else {
			defer outFile.Close()
			fmt.Fprintf(outFile, "Solusi : %s\n", movesStr.String())
			fmt.Fprintf(outFile, "Cost   : %d\n", result.TotalCost)
			fmt.Fprintf(outFile, "Waktu  : %s\n", formatMilliseconds(result.Duration))
			fmt.Fprintf(outFile, "Iterasi: %d\n", result.Iterations)
			fmt.Fprintln(outFile)
			fmt.Fprintln(outFile, "Initial")
			writeBoard(outFile, b, start)
			for i, move := range result.Moves {
				fmt.Fprintf(outFile, "Step %d : %v\n", i+1, move)
				writeBoard(outFile, b, result.Snapshots[i])
			}
			fmt.Printf(">> Solusi disimpan pada %s\n", outPath)
		}
	}
}

func formatMilliseconds(d time.Duration) string {
	if d < time.Microsecond {
		return "<0.001 ms"
	}
	return fmt.Sprintf("%.3f ms", float64(d.Nanoseconds())/1_000_000)
}

func runPlayback(reader *bufio.Reader, b *board.Board, result algorithm.Result, states []board.State) {
	step := 0
	drawPlaybackStep(b, result, states, step)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		runLinePlayback(reader, b, result, states)
		return
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 8)
	for {
		n, err := os.Stdin.Read(buf)
		if err != nil {
			return
		}
		if n == 1 && (buf[0] == 'q' || buf[0] == 'Q') {
			fmt.Print("\033[2J\033[H")
			return
		}
		if n == 1 && (buf[0] == 'd' || buf[0] == 'D' || buf[0] == 'n' || buf[0] == 'N') {
			if step < len(states)-1 {
				step++
				drawPlaybackStep(b, result, states, step)
			}
			continue
		}
		if n == 1 && (buf[0] == 'a' || buf[0] == 'A' || buf[0] == 'p' || buf[0] == 'P') {
			if step > 0 {
				step--
				drawPlaybackStep(b, result, states, step)
			}
			continue
		}
		if n == 1 && (buf[0] == 'j' || buf[0] == 'J') {
			term.Restore(int(os.Stdin.Fd()), oldState)
			step = promptPlaybackStep(reader, len(states)-1, step)
			oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				drawPlaybackStep(b, result, states, step)
				runLinePlayback(reader, b, result, states)
				return
			}
			drawPlaybackStep(b, result, states, step)
			continue
		}
		if n >= 3 && buf[0] == 27 && buf[1] == '[' {
			switch buf[2] {
			case 'C':
				if step < len(states)-1 {
					step++
					drawPlaybackStep(b, result, states, step)
				}
			case 'D':
				if step > 0 {
					step--
					drawPlaybackStep(b, result, states, step)
				}
			}
			continue
		}
		if n >= 2 && (buf[0] == 0 || buf[0] == 224) {
			switch buf[1] {
			case 77:
				if step < len(states)-1 {
					step++
					drawPlaybackStep(b, result, states, step)
				}
			case 75:
				if step > 0 {
					step--
					drawPlaybackStep(b, result, states, step)
				}
			}
			continue
		}
		if n >= 1 && buf[0] == 27 {
			term.Restore(int(os.Stdin.Fd()), oldState)
			step = promptPlaybackStep(reader, len(states)-1, step)
			oldState, err = term.MakeRaw(int(os.Stdin.Fd()))
			if err != nil {
				drawPlaybackStep(b, result, states, step)
				runLinePlayback(reader, b, result, states)
				return
			}
			drawPlaybackStep(b, result, states, step)
		}
	}
}

func runLinePlayback(reader *bufio.Reader, b *board.Board, result algorithm.Result, states []board.State) {
	step := 0
	for {
		drawPlaybackStep(b, result, states, step)
		fmt.Print(">> Masukan step, n untuk next, p untuk prev, atau q untuk keluar: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(strings.ToLower(text))
		switch text {
		case "q", "quit", "exit":
			return
		case "n", "":
			if step < len(states)-1 {
				step++
			}
		case "p":
			if step > 0 {
				step--
			}
		default:
			step = parsePlaybackStep(text, len(states)-1, step)
		}
	}
}

func promptPlaybackStep(reader *bufio.Reader, maxStep, current int) int {
	fmt.Printf("\n>> Lompat ke step berapa? (0-%d): ", maxStep)
	text, _ := reader.ReadString('\n')
	return parsePlaybackStep(strings.TrimSpace(text), maxStep, current)
}

func parsePlaybackStep(text string, maxStep, current int) int {
	step, err := strconv.Atoi(text)
	if err != nil || step < 0 || step > maxStep {
		return current
	}
	return step
}

func drawPlaybackStep(b *board.Board, result algorithm.Result, states []board.State, step int) {
	fmt.Print("\033[H\033[2J")
	if step == 0 {
		fmt.Printf("Initial | Step 0/%d\n", len(states)-1)
	} else {
		fmt.Printf("Step %d/%d | Move: %v\n", step, len(states)-1, result.Moves[step-1])
	}
	fmt.Println("Arrow kiri/kanan atau A/D: mundur/maju | ESC/J: lompat step | q: keluar")
	printBoard(b, states[step])
	fmt.Print("\033[J")
}

func printBoard(b *board.Board, s board.State) {
	writeBoard(os.Stdout, b, s)
}

func writeBoard(w io.Writer, b *board.Board, s board.State) {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if s.Pos == (board.Position{Row: r, Col: c}) {
				fmt.Fprint(w, "Z")
			} else {
				t := b.Grid[r][c]
				if board.IsCheckpoint(t) && s.HasVisited(board.CheckpointIndex(t)) {
					fmt.Fprint(w, "*")
				} else {
					fmt.Fprintf(w, "%c", t)
				}
			}
		}
		fmt.Fprintln(w)
	}
}
