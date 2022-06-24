package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/eliben/go-sudoku"
)

func main() {
	statsFlag := flag.Bool("stats", false, "enable stats for solving")
	flag.Parse()

	filename := flag.Args()[0]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	var totalDuration time.Duration = 0
	var maxDuration time.Duration = 0
	var totalSearches uint64 = 0
	var numBoards int = 0
	var numSolved int = 0

	if *statsFlag {
		sudoku.EnableStats = true
	}

	// Expect one board per line, ignoring whitespace and lines starting with '#'.
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		board := strings.TrimSpace(scanner.Text())
		if len(board) == 0 || strings.HasPrefix(board, "#") {
			continue
		}

		numBoards++

		tStart := time.Now()
		v, err := sudoku.ParseBoard(board)
		v, _ = sudoku.Solve(v)
		if err != nil {
			log.Fatal(err)
		}
		tElapsed := time.Now().Sub(tStart)

		totalDuration += tElapsed
		if tElapsed > maxDuration {
			maxDuration = tElapsed
		}

		if sudoku.IsSolved(v) {
			numSolved++
		}

		if *statsFlag {
			totalSearches += sudoku.Stats.NumSearches
			sudoku.Stats.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Solved %v boards out of %v\n", numSolved, numBoards)
	fmt.Println("Average duration:", totalDuration/time.Duration(numBoards))
	fmt.Println("Max duration:", maxDuration)
	if *statsFlag {
		fmt.Printf("Average searches: %.2f per board\n", float64(totalSearches)/float64(numBoards))
	}
}
