package sudoku

import (
	"fmt"
	"testing"
)

func TestEvaluateDifficulty(t *testing.T) {
	t.Skip()
	h := "85...24..72......9..4.........1.7..23.5...9...4...........8..7..17..........36.4."
	v, err := ParseBoard(h, false)
	fmt.Println(DisplayAsInput(v))
	if err != nil {
		t.Fatal(err)
	}

	EvaluateDifficulty(v)
}
