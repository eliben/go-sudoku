package sudoku

import (
	"strings"
	"testing"
)

func FuzzParseboard(f *testing.F) {
	f.Add("4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......")
	f.Add(DisplayAsInput(EmptyBoard()))
	f.Add(`
. . . |. . 6 |. . .
. 5 9 |. . . |. . 8
2 . . |. . 8 |. . .
------+------+------
. 4 5 |. . . |. . .
. . 3 |. . . |. . .
. . 6 |. . 3 |. 5 4
------+------+------
. . . |3 2 5 |. . 6
. . . |. . . |. . .
. . . |. . . |. . .`)
	f.Add("1234")
	f.Add("")
	f.Add(strings.Repeat("1", 81))

	f.Fuzz(func(t *testing.T, bstr string) {
		// Mostly checks that there are no panics, and ParseBoard returns either
		// a board or and error, but not both.
		b, err := ParseBoard(bstr, true)
		if b != nil && err != nil {
			t.Fatalf("expect b or err to be nil, got b=%v, err=%v", b, err)
		}
		if b == nil && err == nil {
			t.Fatalf("expect b or err to be non-nil, got b=%v, err=%v", b, err)
		}
	})
}
