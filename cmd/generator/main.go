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
	f := sudoku.Generate(30)
	fmt.Println(f)
	fmt.Println(sudoku.Display(f))
	fmt.Println(sudoku.DisplayAsInput(f))

	c := 0
	for _, d := range f {
		if d.Size() == 1 {
			c++
		}
	}
	fmt.Println("Hint squares:", c)

	if !sudoku.EliminateAll(f) {
		log.Fatal("contradiction")
	}
	fmt.Println(sudoku.DisplayAsInput(f))

	sols := sudoku.SolveAll(f, -1)
	fmt.Println("Solutions:", len(sols))
}
