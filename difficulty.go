package sudoku

import "fmt"

// TODO: doc

// 2. Count hints after elimination
// 3. Count the low bound on empty rows/cols pre (or after?) elimination
// 4. Count how difficult average (maximal?) search is over a few random tries
func EvaluateDifficulty(values Values) int {
	// Count hints.
	hintcount := 0
	for _, d := range values {
		if d.Size() == 1 {
			hintcount++
		}
	}

	fmt.Println("hintcount =", hintcount)

	return 0
}
