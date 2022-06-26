package sudoku

import (
	"fmt"
	"testing"
)

var h1 = "85...24..72......9..4.........1.7..23.5...9...4...........8..7..17..........36.4."
var h2 = "..53.....8......2..7..1.5..4....53...1..7...6..32...8..6.5....9..4....3......97.."

func TestEvaluateDifficulty(t *testing.T) {
	v, err := ParseBoard(h2, false)
	fmt.Println(DisplayAsInput(v))
	if err != nil {
		t.Fatal(err)
	}

	EvaluateDifficulty(v)
}
