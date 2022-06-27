package sudoku

import (
	"log"
	"math/rand"
)

// Generate generates a random Sudoku board that has a single solution, with
// at-most hintCount hints remaining on the board. Note that this cannot be
// always reliably done when the count is low (lower than 23 or so), because
// generating a board with a single solution that has a low number of initial
// hints is very hard.
// There are no guarantees made about the difficulty of the generated board,
// though higher hint counts generally correlate with easier boards. It's
// recommended to generate a large number of boards using this function and
// evaluate their difficulty separately using EvaluateDifficulty.
// Note: make sure the default rand source is seeded if you really want to get
// random boards.
func Generate(hintCount int) Values {
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
		board[sq] = FullDigitsSet()

		solutions := SolveAll(board, 2)
		switch len(solutions) {
		case 0:
			// Some sort of bug, because removing a square from a solved board should
			// never result in an unsolvable board.
			log.Fatal("got a board without solutions")
		case 1:
			count--
			if count <= hintCount {
				return board
			}
		default:
			// The board has multiple solutions with this square emptied, so put it
			// back and try again with the next square.
			board[sq] = savedDigit
		}
	}

	return board
}
