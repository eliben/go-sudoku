package sudoku

import (
	"log"
)

func Generate() Values {
	empty := EmptyBoard()
	full, solved := Solve(empty, SolveOptions{Randomize: true})
	if !solved || !IsSolved(full) {
		log.Fatal("unable to generate solved board from empty")
	}
	return full
}
