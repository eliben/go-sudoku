package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/eliben/go-sudoku"
)

// Note: trying to generate difficult boards with low hintcount may take a long
// time.
const defaultSeed = 0xDEADBEEF

var symFlag = flag.Bool("sym", false, "generate a symmetrical puzzle")
var diffFlag = flag.Float64("diff", 2.5, "minimal difficulty for generated puzzle")
var hintCountFlag = flag.Int("hintcount", 28, "hint count for generation; higher counts lead to easier puzzles")
var svgOutFlag = flag.String("svgout", "", "file name for SVG output, if needed")
var seedFlag = flag.Int64("seed", defaultSeed, "random number generator seed")
var diffMaxFlag = flag.Float64("diffMax", *diffFlag+1.0, "maximum difficulty for a puzzle")

func main() {
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Println(out, "usage: generator [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
	}
	flag.Parse()

	rand.Seed(*seedFlag)

	count := 0
	maxDifficultySeen := 0.0

	for {
		var board, solvedBoard sudoku.Values

		if *symFlag {
			board, solvedBoard = sudoku.GenerateSymmetrical(*hintCountFlag)
		} else {
			board, solvedBoard = sudoku.Generate(*hintCountFlag)
		}

		d, err := sudoku.EvaluateDifficulty(board)
		if err != nil {
			log.Fatal(err)
		}

		if d >= *diffFlag && d <= *diffMaxFlag {
			fmt.Println(sudoku.DisplayAsInput(board))
			fmt.Printf("Difficulty: %.2f\n", d)

			fmt.Println(sudoku.DisplayAsInput(solvedBoard))
			boardString := sudoku.DisplayAsString(board)
			solvedBoardString := sudoku.DisplayAsString(solvedBoard)
			fmt.Println("PuzzleString:", boardString)
			fmt.Println("SolvedString:", solvedBoardString)
			fmt.Printf("CartaginaOutput:%v,3x3,%v,%v,%v,%v,%v\n", *seedFlag, *diffFlag, *diffMaxFlag, *hintCountFlag, boardString, solvedBoardString)
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
