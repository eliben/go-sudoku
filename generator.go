package sudoku

import (
	"log"
	"math/rand"
)

// TODO: document that the rng has to be seeded
func Generate(fillCount int) Values {
	empty := EmptyBoard()
	board, solved := Solve(empty, SolveOptions{Randomize: true})
	if !solved || !IsSolved(board) {
		log.Fatal("unable to generate solved board from empty")
	}

	removalOrder := rand.Perm(81)
	count := 81

	for _, sq := range removalOrder {
		savedDigit := board[sq]
		// Try to remove the number from square sq.
		board[sq] = fullDigitsSet()

		solutions := SolveAll(board, 2)
		switch len(solutions) {
		case 0:
			log.Fatal("got a board without solutions")
		case 1:
			count--
			if count <= fillCount {
				return board
			}
		default:
			board[sq] = savedDigit
		}
	}

	return board
}
