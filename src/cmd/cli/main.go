package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tucil3/internal/algorithm"
	"tucil3/internal/board"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(">> Masukan file input: ")
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

	fmt.Print(">> Algoritma apa yang anda pilih?: ")
	algoStr, _ := reader.ReadString('\n')
	algoStr = strings.TrimSpace(strings.ToUpper(algoStr))

	var algo algorithm.Algorithm
	switch algoStr {
	case "UCS", "":
		algo = algorithm.UCS{}
	default:
		fmt.Fprintln(os.Stderr, "Algoritma tidak dikenal (hanya UCS tersedia saat ini):", algoStr)
		os.Exit(1)
	}

	start := board.State{Pos: b.Start}
	result := algo.Search(b, start, nil)

	if !result.Found {
		fmt.Println("Tidak ditemukan solusi.")
		return
	}

	fmt.Print("\nSolusi Yang Ditemukan: ")
	for _, m := range result.Moves {
		fmt.Print(m)
	}
	fmt.Printf("\nCost dari Solusi  : %d\n", result.TotalCost)
	fmt.Printf("Banyak iterasi    : %d\n", result.Iterations)
	fmt.Printf("Waktu eksekusi    : %d ms\n", result.Duration.Milliseconds())

	fmt.Println("\n--- Initial ---")
	printBoard(b, start)
	for i, move := range result.Moves {
		fmt.Printf("\n--- Step %d : %v ---\n", i+1, move)
		printBoard(b, result.Snapshots[i])
	}
}

func printBoard(b *board.Board, s board.State) {
	for r := 0; r < b.Rows; r++ {
		for c := 0; c < b.Cols; c++ {
			if s.Pos == (board.Position{Row: r, Col: c}) {
				fmt.Print("Z")
			} else {
				fmt.Printf("%c", b.Grid[r][c])
			}
		}
		fmt.Println()
	}
}
