package sudoku

import (
	"math/rand"
	"testing"
	"time"
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

func TestGenerateSymmetrical(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	//for {
	board := GenerateSymmetrical(30)
	vs := SolveAll(board, -1)
	if len(vs) != 1 {
		t.Errorf("got %v solutions, want 1", len(vs))
	}

	if !IsSolved(vs[0]) {
		t.Errorf("got unsolved board")
	}

	// Check symmetry
	for sq := 0; sq < 41; sq++ {
		if board[sq].Size() != board[80-sq].Size() {
			t.Errorf("squares %v != %v on board, expected symmetry", sq, 80-sq)
		}
	}
}
