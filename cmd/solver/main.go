package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/eliben/go-sudoku"
)

func main() {
	flag.Parse()

	filename := flag.Args()[0]

	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	// Expect one board per line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		board := strings.TrimSpace(scanner.Text())
		if len(board) == 0 {
			continue
		}

		v, err := sudoku.SolveBoard(board)
		if err != nil {
			log.Fatalf("unable to parse/solve board:", err)
		}

		fmt.Println(v)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
