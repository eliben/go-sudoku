package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/eliben/go-sudoku"
)

var symFlag = flag.Bool("sym", false, "generate a symmetrical puzzle")
var diffFlag = flag.Float64("diff", 2.5, "minimal difficulty for generated puzzle")
var hintCountFlag = flag.Int("hintcount", 28, "hint count for generation; higher counts lead to easier puzzles")

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Println(out, "usage: generator [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	for {
		var f sudoku.Values

		if *symFlag {
			f = sudoku.GenerateSymmetrical(*hintCountFlag)
		} else {
			f = sudoku.Generate(*hintCountFlag)
		}

		sols := sudoku.SolveAll(f, -1)
		if len(sols) != 1 {
			continue
		}

		d, err := sudoku.EvaluateDifficulty(f)
		if err != nil {
			log.Fatal(err)
		}
		if d >= *diffFlag {
			fmt.Println(sudoku.DisplayAsInput(f))
			fmt.Printf("Difficulty: %.2f\n", d)
			break
		}
	}
}
