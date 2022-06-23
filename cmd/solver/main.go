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
	flag.Parse()

	filename := flag.Args()[0]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	var totalDuration time.Duration = 0
	var maxDuration time.Duration = 0
	var numBoards int = 0

	// Expect one board per line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		board := strings.TrimSpace(scanner.Text())
		if len(board) == 0 || strings.HasPrefix(board, "#") {
			continue
		}

		numBoards++

		tStart := time.Now()
		v, err := sudoku.SolveBoard(board)
		tElapsed := time.Now().Sub(tStart)
		if err != nil {
			log.Fatalf("unable to parse/solve board:", err)
		}

		totalDuration += tElapsed
		if tElapsed > maxDuration {
			maxDuration = tElapsed
		}

		if !sudoku.IsSolved(v) {
			log.Fatalf("board %v is not solved", v)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Total boards solved:", numBoards)
	fmt.Println("Average duration:", totalDuration/time.Duration(numBoards))
	fmt.Println("Max duration:", maxDuration)
}
