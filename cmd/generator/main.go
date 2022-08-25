package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/eliben/go-sudoku"
)

var symFlag = flag.Bool("sym", false, "generate a symmetrical puzzle")
var diffFlag = flag.Float64("diff", 2.5, "minimal difficulty for generated puzzle")
var hintCountFlag = flag.Int("hintcount", 28, "hint count for generation; higher counts lead to easier puzzles")
var svgOutFlag = flag.String("svgout", "", "file name for SVG output, if needed")

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Println(out, "usage: generator [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	count := 0
	maxDifficultySeen := 0.0

	for {
		var board sudoku.Values

		if *symFlag {
			board = sudoku.GenerateSymmetrical(*hintCountFlag)
		} else {
			board = sudoku.Generate(*hintCountFlag)
		}

		sols := sudoku.SolveAll(board, -1)
		if len(sols) != 1 {
			continue
		}

		d, err := sudoku.EvaluateDifficulty(board)
		if err != nil {
			log.Fatal(err)
		}

		if d >= *diffFlag {
			fmt.Println(sudoku.DisplayAsInput(board))
			fmt.Printf("Difficulty: %.2f\n", d)

			if len(*svgOutFlag) > 0 {
				f, err := os.Create(*svgOutFlag)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()
				sudoku.DisplayAsSVG(f, board, d)
				fmt.Println("Wrote SVG output to", *svgOutFlag)
			}

			break
		} else {
			count++
			if d > maxDifficultySeen {
				maxDifficultySeen = d
			}

			if count > 0 && count%10 == 0 {
				fmt.Printf("Tried %v boards; max difficulty seen %.2f\n", count, maxDifficultySeen)
			}
		}
	}
}
