package sudoku

import (
	"fmt"
)

// TODO: doc

// The before/after elimination distinction is very important here...
// 1. Count hints before elimination
// 2. Count hints after elimination
// 3. Count the low bound on empty rows/cols pre (or after?) elimination
// 4. Count how difficult average (maximal?) search is over a few random tries
func EvaluateDifficulty(values Values) (int, error) {
	countHints := func() int {
		hintcount := 0
		for _, d := range values {
			if d.Size() == 1 {
				hintcount++
			}
		}
		return hintcount
	}

	hintsBeforeElimination := countHints()

	// Count the lower bound (minimal number) of hints in rows and cols, pre
	// elimination.
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

	if !EliminateAll(values) {
		return 0, fmt.Errorf("contradiction in board")
	}
	hintsAfterElimination := countHints()

	EnableStats = true
	var totalSearches uint64 = 0
	iterations := 100
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
	} else if hintsAfterElimination > 40 {
		hintsAfterDifficulty = 2.0
	} else if hintsAfterElimination > 36 {
		hintsAfterDifficulty = 3.0
	} else if hintsAfterElimination > 32 {
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
	} else if minHints == 2 {
		minHintsDifficulty = 4.0
	} else {
		minHintsDifficulty = 5.0
	}

	var searchDifficulty float64
	if averageSearches <= 1.0 {
		searchDifficulty = 1.0
	} else if averageSearches < 3.0 {
		searchDifficulty = 2.0
	} else if averageSearches < 8.0 {
		searchDifficulty = 3.0
	} else if averageSearches < 25.0 {
		searchDifficulty = 4.0
	} else {
		searchDifficulty = 5.0
	}

	fmt.Printf("hintsBeforeElimination: %v, hintsBeforeDifficulty: %v\n", hintsBeforeElimination, hintsBeforeDifficulty)
	fmt.Printf("hintsAfterlimination: %v, hintsAfterifficulty: %v\n", hintsAfterElimination, hintsAfterDifficulty)
	fmt.Printf("minHints: %v, minHintsDifficulty: %v\n", minHints, minHintsDifficulty)
	fmt.Printf("averageSearches: %v, searchDifficulty: %v\n", averageSearches, searchDifficulty)

	return 0, nil
}
