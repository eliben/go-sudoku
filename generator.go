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
// Notes:
//   - Make sure the default rand source is seeded if you really want to get
//     random boards.
//   - This function may take a while to run when given a low hintCount.
func Generate(hintCount int) (Values, Values) {
	empty := EmptyBoard()
	board, solved := Solve(empty, SolveOptions{Randomize: true})
	if !solved || !IsSolved(board) {
		log.Fatal("unable to generate solved board from empty")
	}

	removalOrder := rand.Perm(81)
	count := 81
	solvedBoard := make(Values, 81)
	copy(solvedBoard, board)

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
				return board, solvedBoard
			}
		default:
			// The board has multiple solutions with this square emptied, so put it
			// back and try again with the next square.
			count++
			board[sq] = savedDigit
		}
	}

	return board, solvedBoard
}

// GenerateSymmetrical is similar to Generate, but it generates symmetrical
// boards with 180-degree rotational symmetry.
// Because of this additional constraint, it may have more trouble generating
// boards with a small hintCount than Generate, so you'll have to run it more
// times in a loop to find a good low-hint-count board.
func GenerateSymmetrical(hintCount int) (Values, Values) {
	empty := EmptyBoard()
	board, solved := Solve(empty, SolveOptions{Randomize: true})
	if !solved || !IsSolved(board) {
		log.Fatal("unable to generate solved board from empty")
	}

	solvedBoard := make(Values, 81)
	copy(solvedBoard, board)

	// This function works just like Generate, but instead of picking a random
	// square out of all 81, it picks a random square from the first half of the
	// board and then attempts to remove both this square and its reflection.
	removalOrder := rand.Perm(41)
	count := 81

	for _, sq := range removalOrder {
		// Find sq's reflection; note that in the middle row reflectSq could equal
		// sq - we take this into account when counting how many hints remain on
		// the board.
		reflectSq := 80 - sq

		savedDigit := board[sq]
		savedReflect := board[reflectSq]

		board[sq] = FullDigitsSet()
		board[reflectSq] = FullDigitsSet()

		solutions := SolveAll(board, 2)
		switch len(solutions) {
		case 0:
			log.Fatal("got a board without solutions")
		case 1:
			// We may have removed just one or two hints.
			count--
			if sq != reflectSq {
				count--
			}
			if count <= hintCount {
				return board, solvedBoard
			}
		default:
			count++
			board[sq] = savedDigit
			board[reflectSq] = savedReflect
		}
	}

	return board, solvedBoard
}
