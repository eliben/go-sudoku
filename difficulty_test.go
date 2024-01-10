package sudoku

import (
	"log"
	"testing"

	"golang.org/x/exp/slices"
)

func TestEvaluateDifficulty(t *testing.T) {
	getDifficulty := func(board string) float64 {
		v, err := ParseBoard(board, false)
		if err != nil {
			t.Fatal(err)
		}
		vcopy := slices.Clone(v)

		d, err := EvaluateDifficulty(v)
		if err != nil {
			t.Fatal(err)
		}

		if !slices.Equal(v, vcopy) {
			t.Errorf("EvaluateDifficulty modified values")
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

var filled string = `
3 4 5 |7 9 2 |6 1 8 
8 9 6 |3 5 1 |2 7 4 
1 7 2 |8 6 4 |9 5 3 
------+------+------
2 3 8 |6 4 5 |1 9 7 
7 5 4 |9 1 8 |3 2 6 
6 1 9 |2 3 7 |4 8 5 
------+------+------
5 6 7 |4 2 9 |8 3 1 
4 2 1 |5 8 3 |7 6 9 
9 8 3 |1 7 6 |5 4 2 
`

func TestFilled(t *testing.T) {
	v, err := ParseBoard(filled, false)
	if err != nil {
		log.Fatal(err)
	}

	d, err := EvaluateDifficulty(v)
	if err != nil {
		log.Fatal(err)
	}

	if d != 1.0 {
		t.Errorf("got d=%v; expect difficulty of filled board to be 1.0", d)
	}
}
