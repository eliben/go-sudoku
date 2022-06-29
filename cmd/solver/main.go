package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/eliben/go-sudoku"
)

// TODO: add real help / flags
var statsFlag = flag.Bool("stats", false, "enable stats for solving")
var randomizeFlag = flag.Bool("randomize", false, "randomize solving order")
var actionFlag = flag.String("action", "solve", "action to perform: solve, count")

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintln(out, "usage: solver [options] <input in stdin>")
		fmt.Fprintln(out, "Options:")
		flag.PrintDefaults()
	}

	flag.Parse()

	switch *actionFlag {
	case "solve":
		solveAndReport()
	case "count":
		countHints()
	default:
		flag.Usage()
		log.Fatal("Please select one of the supported actions.")
	}
}

func solveAndReport() {
	var totalDuration time.Duration = 0
	var maxDuration time.Duration = 0
	var totalSearches uint64 = 0
	var totalDifficulty float64
	var maxSearches uint64 = 0
	var numBoards int = 0
	var numSolved int = 0

	if *statsFlag {
		sudoku.EnableStats = true
	}

	if *randomizeFlag {
		rand.Seed(time.Now().UnixNano())
	}

	boards := getInputBoards()
	for _, board := range boards {
		numBoards++

		v, err := sudoku.ParseBoard(board, false)
		if err != nil {
			log.Fatal(err)
		}
		d, err := sudoku.EvaluateDifficulty(v)
		if err != nil {
			log.Fatal(err)
		}
		totalDifficulty += d

		tStart := time.Now()
		sudoku.EliminateAll(v)
		v, _ = sudoku.Solve(v, sudoku.SolveOptions{Randomize: *randomizeFlag})
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
			if sudoku.Stats.NumSearches > maxSearches {
				maxSearches = sudoku.Stats.NumSearches
			}
			sudoku.Stats.Reset()
		}
	}

	fmt.Printf("Solved %v/%v boards\n", numSolved, numBoards)
	fmt.Printf("Average difficulty: %.2v\n", totalDifficulty/float64(numBoards))
	fmt.Printf("Duration average=%-15v max=%v\n", totalDuration/time.Duration(numBoards), maxDuration)
	if *statsFlag {
		fmt.Printf("Searches average=%-15.2f max=%v\n", float64(totalSearches)/float64(numBoards), maxSearches)
	}
}

func countHints() {
	boards := getInputBoards()
	for _, board := range boards {
		fmt.Println("board:", board)
		v, err := sudoku.ParseBoard(board, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("|")

		initialNumHints := sudoku.CountHints(v)
		fmt.Printf("\\ initial num hints:           %v\n", initialNumHints)

		sudoku.EliminateAll(v)
		afterElimNumHints := sudoku.CountHints(v)
		fmt.Printf("  num hints after elimination: %v\n", afterElimNumHints)

		sudoku.ApplyTwinsStrategy(v)
		afterTwinsNumHints := sudoku.CountHints(v)
		fmt.Printf("  num hints after twins:       %v\n", afterTwinsNumHints)
		fmt.Println("")
	}
}

// getInputBoards reads input boards from stdin, ignores comments and empty
// lines and returns them.
func getInputBoards() []string {
	var boards []string

	scanner := bufio.NewScanner(os.Stdin)
	// Expect one board per line, ignoring whitespace and lines starting with '#'.
	for scanner.Scan() {
		board := strings.TrimSpace(scanner.Text())
		if len(board) == 0 || strings.HasPrefix(board, "#") {
			continue
		}
		boards = append(boards, board)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return boards
}
