package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/eliben/go-sudoku"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	f := sudoku.Generate(22)
	fmt.Println(f)
	fmt.Println(sudoku.Display(f))
	fmt.Println(sudoku.DisplayAsInput(f))

	sols := sudoku.SolveAll(f, -1)
	fmt.Println("Solutions:", len(sols))

	d, err := sudoku.EvaluateDifficulty(f)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)
}
