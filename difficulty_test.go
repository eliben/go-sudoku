package sudoku

import (
	"math/rand"
	"testing"
	"time"
)

func TestEvaluateDifficulty(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	getDifficulty := func(board string) float64 {
		v, err := ParseBoard(board, false)
		if err != nil {
			t.Fatal(err)
		}

		d, err := EvaluateDifficulty(v)
		if err != nil {
			t.Fatal(err)
		}

		return d
	}

	easyD := getDifficulty(easyboard1)
	hardD := getDifficulty(hardboard2)
	hardlongD := getDifficulty(hardlong)

	if easyD > hardD || hardD > hardlongD {
		t.Errorf("got easyD: %v, hardD: %v, hardlongD: %v", easyD, hardD, hardlongD)
	}
}
