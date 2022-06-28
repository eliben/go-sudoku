package sudoku

import (
	"fmt"
)

// EvaluateDifficulty evaluates the difficulty of a Sudoku puzzle heuristically
// and returns the score on a scale from 1.0 (easiest) to 5.0 hardest. It can
// also return an error if the given board has contradictions, is unsolvable,
// etc. It should be passed a board that didn't have elimination applied to it.
//
// The heuristics are based on 4 factors:
//
// 1. How many hints (filled-in squares) the board has.
// 2. How many hints remain after running a round of elimination (first-order
//    Sudoku solving value deduction).
// 3. How many hints does a row or column with the minimal number of hints have
// 4. How many guesses a backtracking search requires to solve the board
//    (averaged over multiple runs).
//
// This approach was partially inspired by the paper "Sudoku Puzzles Generating:
// from Easy to Evil" by Xiang-Sun ZHANG's research group.
func EvaluateDifficulty(values Values) (float64, error) {
	hintsBeforeElimination := CountHints(values)

	// Count the lower bound (minimal number) of hints in individual rows and
	// cols, pre elimination.
	minHints := 9

	index := func(row, col int) Index {
		return row*9 + col
	}

	// ... first the rows.
	for row := 0; row < 9; row++ {
		rowCount := 0
		for col := 0; col < 9; col++ {
			if values[index(row, col)].Size() == 1 {
				rowCount++
			}
		}
		if rowCount < minHints {
			minHints = rowCount
		}
	}

	// ... then the columns.
	for col := 0; col < 9; col++ {
		colCount := 0
		for row := 0; row < 9; row++ {
			if values[index(row, col)].Size() == 1 {
				colCount++
			}
		}
		if colCount < minHints {
			minHints = colCount
		}
	}

	// Run elimination and count how many hints are on the board after it.
	if !EliminateAll(values) {
		return 0, fmt.Errorf("contradiction in board")
	}
	hintsAfterElimination := CountHints(values)

	// Run a number of randomized searches and count the average search count.
	EnableStats = true
	var totalSearches uint64 = 0
	iterations := 10
	for i := 0; i < iterations; i++ {
		Stats.Reset()
		_, solved := Solve(values, SolveOptions{Randomize: true})
		if !solved {
			return 0, fmt.Errorf("cannot solve")
		}
		totalSearches += Stats.NumSearches
	}
	EnableStats = false
	averageSearches := float64(totalSearches) / float64(iterations)

	// Assign difficulty scores based on ranges in each category.
	var hintsBeforeDifficulty float64
	if hintsBeforeElimination > 50 {
		hintsBeforeDifficulty = 1.0
	} else if hintsBeforeElimination > 35 {
		hintsBeforeDifficulty = 2.0
	} else if hintsBeforeElimination > 31 {
		hintsBeforeDifficulty = 3.0
	} else if hintsBeforeElimination > 27 {
		hintsBeforeDifficulty = 4.0
	} else {
		hintsBeforeDifficulty = 5.0
	}

	var hintsAfterDifficulty float64
	if hintsAfterElimination > 55 {
		hintsAfterDifficulty = 1.0
	} else if hintsAfterElimination > 42 {
		hintsAfterDifficulty = 2.0
	} else if hintsAfterElimination > 37 {
		hintsAfterDifficulty = 3.0
	} else if hintsAfterElimination > 33 {
		hintsAfterDifficulty = 4.0
	} else {
		hintsAfterDifficulty = 5.0
	}

	var minHintsDifficulty float64
	if minHints >= 5 {
		minHintsDifficulty = 1.0
	} else if minHints == 4 {
		minHintsDifficulty = 2.0
	} else if minHints == 3 {
		minHintsDifficulty = 3.0
	} else if minHints >= 1 {
		minHintsDifficulty = 4.0
	} else {
		minHintsDifficulty = 5.0
	}

	var searchDifficulty float64
	if averageSearches <= 1.0 {
		searchDifficulty = 1.0
	} else if averageSearches < 3.0 {
		searchDifficulty = 2.0
	} else if averageSearches < 10.0 {
		searchDifficulty = 3.0
	} else if averageSearches < 40.0 {
		searchDifficulty = 4.0
	} else {
		searchDifficulty = 5.0
	}

	// Assign final difficulty with weights
	difficulty := 0.5*hintsAfterDifficulty +
		0.3*hintsBeforeDifficulty +
		0.05*minHintsDifficulty +
		0.15*searchDifficulty

	return difficulty, nil
}
