package sudoku

import "fmt"

// TODO: doc

// The before/after elimination distinction is very important here...
// 1. Count hints before elimination
// 2. Count hints after elimination
// 3. Count the low bound on empty rows/cols pre (or after?) elimination
// 4. Count how difficult average (maximal?) search is over a few random tries
func EvaluateDifficulty(values Values) (int, error) {
	countHits := func() int {
		hintcount := 0
		for _, d := range values {
			if d.Size() == 1 {
				hintcount++
			}
		}
		return hintcount
	}

	fmt.Println("hintcount before elimination:", countHits())

	if !EliminateAll(values) {
		return 0, fmt.Errorf("contradiction in board")
	}

	fmt.Println("hintcount after elimination:", countHits())
	return 0, nil
}
