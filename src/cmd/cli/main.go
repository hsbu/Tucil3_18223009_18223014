package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"tucil3/internal/algorithm"
	"tucil3/internal/board"
	"tucil3/internal/heuristic"
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

	// ── Step visualization ──
	fmt.Println("Initial")
	printBoard(b, start)

	for i, move := range result.Moves {
		fmt.Printf("\nStep %d : %v\n", i+1, move)
		printBoard(b, result.Snapshots[i])
	}

	// ── Execution stats ──
	fmt.Printf("\n>> Waktu eksekusi: %s\n", formatDuration(result.Duration))
	fmt.Printf(">> Banyak iterasi yang dilakukan: %d iterasi\n", result.Iterations)

	// ── Playback prompt ──
	fmt.Print("\n>> Apakah Anda ingin melakukan playback? (Ya/Tidak) :\n")
	playbackAns, _ := reader.ReadString('\n')
	playbackAns = strings.TrimSpace(strings.ToLower(playbackAns))

	if playbackAns == "ya" || playbackAns == "y" {
		allStates := make([]board.State, 0, len(result.Snapshots)+1)
		allStates = append(allStates, start)
		allStates = append(allStates, result.Snapshots...)

		running := true
		for running {
			fmt.Print(">> Pada step berapa anda ingin melakukan playback :\n")
			stepStr, _ := reader.ReadString('\n')
			stepStr = strings.TrimSpace(stepStr)
			if stepStr == "" {
				break
			}
			step, err := strconv.Atoi(stepStr)
			if err != nil || step < 0 || step > len(result.Moves) {
				fmt.Println("Step tidak valid. Masukan angka antara 0 dan", len(result.Moves))
				continue
			}

			st := allStates[step]
			if step == 0 {
				fmt.Println("Initial")
			} else {
				fmt.Printf("Step %d : %v\n", step, result.Moves[step-1])
			}
			printBoard(b, st)

			fmt.Print("\n>> Masukan step selanjutnya (atau ketik 'exit' untuk keluar) :\n")
			nextStr, _ := reader.ReadString('\n')
			nextStr = strings.TrimSpace(strings.ToLower(nextStr))
			if nextStr == "exit" || nextStr == "" {
				running = false
			} else {
				if n, err := strconv.Atoi(nextStr); err == nil && n >= 0 && n <= len(result.Moves) {
					fmt.Printf("Step %d : %v\n", n, result.Moves[n-1])
					printBoard(b, allStates[n])
				}
			}
		}
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
			fmt.Fprintf(outFile, "Waktu  : %s\n", formatDuration(result.Duration))
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

func formatDuration(d time.Duration) string {
	if d <= 0 {
		return "<1 us"
	}
	if d < time.Microsecond {
		return "<1 us"
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.3f us", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.3f ms", float64(d.Nanoseconds())/1_000_000)
	}
	return fmt.Sprintf("%.3f s", d.Seconds())
}

func printBoard(b *board.Board, s board.State) {
	writeBoard(os.Stdout, b, s)
}

func writeBoard(w *os.File, b *board.Board, s board.State) {
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
