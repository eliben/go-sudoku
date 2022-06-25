package sudoku

import (
	"fmt"
	"log"
	"math/rand"
)

// TODO: document that the rng has to be seeded
func Generate(hintCount int) Values {
	empty := EmptyBoard()
	board, solved := Solve(empty, SolveOptions{Randomize: true})
	if !solved || !IsSolved(board) {
		log.Fatal("unable to generate solved board from empty")
	}

	removalOrder := rand.Perm(81)
	count := 81

	for _, sq := range removalOrder {
		fmt.Println("... trying to remove from", sq)
		savedDigit := board[sq]
		// Try to remove the number from square sq.
		board[sq] = FullDigitsSet()

		solutions := SolveAll(board, 2)
		fmt.Println("... # solutions:", len(solutions))
		switch len(solutions) {
		case 0:
			log.Fatal("got a board without solutions")
		case 1:
			count--
			if count <= hintCount {
				return board
			}
		default:
			board[sq] = savedDigit
		}
	}

	return board
}
