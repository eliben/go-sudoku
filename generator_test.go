package sudoku

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	board := Generate(30)
	vs := SolveAll(board, -1)
	if len(vs) != 1 {
		t.Errorf("got %v solutions, want 1", len(vs))
	}

	if !IsSolved(vs[0]) {
		t.Errorf("got unsolved board")
	}
}
