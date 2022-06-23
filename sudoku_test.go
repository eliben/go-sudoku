package sudoku

import (
	"log"
	"strings"
	"testing"

	"golang.org/x/exp/slices"
)

func TestNew(t *testing.T) {
	// Smoke testing.
	s := New()
	if len(s.unitlist) != 27 {
		t.Errorf("got len=%v, want 27", len(s.unitlist))
	}

	wantUnits := []Unit{
		Unit{18, 19, 20, 21, 22, 23, 24, 25, 26},
		Unit{2, 11, 20, 29, 38, 47, 56, 65, 74},
		Unit{0, 1, 2, 9, 10, 11, 18, 19, 20}}

	if !slices.EqualFunc(wantUnits, s.units[20], func(a, b Unit) bool {
		return slices.Equal(a, b)
	}) {
		t.Errorf("got units[20]=%v\nwant %v", s.units[20], wantUnits)
	}

	gotPeers := s.peers[20]
	slices.Sort(gotPeers)
	wantPeers := []Index{0, 1, 2, 9, 10, 11, 18, 19, 21, 22, 23, 24, 25, 26, 29, 38, 47, 56, 65, 74}
	if !slices.Equal(wantPeers, gotPeers) {
		t.Errorf("got peers[20]=%v\n want %v", s.peers[20], wantPeers)
	}
}

func TestAssignElimination(t *testing.T) {
	s := New()
	vals := emptyBoard()

	if s.isSolved(vals) {
		t.Errorf("an empty board is solved")
	}

	// Assign a digit to square 20; check that this digit is the only candidate
	// in square 20, and that it was eliminated from all the peers of 20.
	s.assign(vals, 20, 5)

	if vals[20].size() != 1 || vals[20].singleMemberDigit() != 5 {
		t.Errorf("got vals[20]=%v", vals[20])
	}

	for sq := 0; sq <= 80; sq++ {
		if slices.Contains(s.peers[20], sq) {
			if vals[sq].isMember(5) {
				t.Errorf("got member 5 in peer square %v", sq)
			}
		} else {
			if !vals[sq].isMember(5) {
				t.Errorf("got no member 5 in non-peer square %v", sq)
			}
		}
	}
}

// Easy board from Norvig's example that's solved by constraint propagation
// w/o any search.
var easyboard1 string = "003020600900305001001806400008102900700000008006708200002609500800203009005010300"

var hardboard1 string = "4.....8.5.3..........7......2.....6.....8.4......1.......6.3.7.5..2.....1.4......"
var hardboard2 string = "..53.....8......2..7..1.5..4....53...1..7...6..32...8..6.5....9..4....3......97.."

// This is the longest Norvig's program took to solve a puzzle
var hardlong string = `
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
. . . |. . . |. . .`

func TestParseBoard(t *testing.T) {
	s := New()
	v, err := s.parseBoard(easyboard1)
	if err != nil {
		t.Fatal(err)
	}

	if !s.isSolved(v) {
		t.Errorf("expect easy board to be solved")
	}

	// Harder board that isn't fully solved without search.
	v2, err := s.parseBoard(hardboard1)
	if err != nil {
		t.Fatal(err)
	}

	if s.isSolved(v2) {
		t.Errorf("expect hard board to not be solved")
	}

	// Count how many squares are solved immediately in this puzzle and compare
	// to the number Norvig got.
	var solvedSquares int
	for _, d := range v2 {
		if d.size() == 1 {
			solvedSquares++
		}
	}

	if solvedSquares != 20 {
		t.Errorf("got %v solved squares, want 20", solvedSquares)
	}
}

func TestSolveBoard(t *testing.T) {
	s := New()
	v, err := s.solveBoard(hardboard1)
	if err != nil {
		log.Fatal(err)
	}

	if !s.isSolved(v) {
		t.Errorf("expect hardboard1 to be solved by search")
	}

	// Should work on the easy board also (even though it's solved with the
	// initial parse)
	v2, err := s.solveBoard(easyboard1)
	if err != nil {
		log.Fatal(err)
	}

	if !s.isSolved(v2) {
		t.Errorf("expect easy board to be solved by search")
	}

	// And the other hard board
	v3, err := s.solveBoard(hardboard2)
	if err != nil {
		log.Fatal(err)
	}

	if !s.isSolved(v3) {
		t.Errorf("expect hardboard2 to be solved by search")
	}
}

func TestIsSolved(t *testing.T) {
	s := New()
	v, err := s.parseBoard(easyboard1)
	if err != nil {
		t.Fatal(err)
	}

	if !s.isSolved(v) {
		t.Errorf("expect easy board to be solved")
	}

	// Now modify the board and make sure it's not considered "solved" any more.
	// ... modify by trying to add options to each square separately.
	for sq := range v {
		vcopy := slices.Clone(v)
		vcopy[sq] = vcopy[sq].add(6).add(8)

		if s.isSolved(vcopy) {
			t.Errorf("expect board to not be solved after modification: %v", vcopy)
		}
	}
}

func TestSolveHardest(t *testing.T) {
	// The "hardest" puzzles Norvig found online (taken from
	// https://norvig.com/hardest.txt)
	hardest := `
85...24..72......9..4.........1.7..23.5...9...4...........8..7..17..........36.4.
..53.....8......2..7..1.5..4....53...1..7...6..32...8..6.5....9..4....3......97..
12..4......5.69.1...9...5.........7.7...52.9..3......2.9.6...5.4..9..8.1..3...9.4
...57..3.1......2.7...234......8...4..7..4...49....6.5.42...3.....7..9....18.....
7..1523........92....3.....1....47.8.......6............9...5.6.4.9.7...8....6.1.
1....7.9..3..2...8..96..5....53..9...1..8...26....4...3......1..4......7..7...3..
1...34.8....8..5....4.6..21.18......3..1.2..6......81.52..7.9....6..9....9.64...2
...92......68.3...19..7...623..4.1....1...7....8.3..297...8..91...5.72......64...
.6.5.4.3.1...9...8.........9...5...6.4.6.2.7.7...4...5.........4...8...1.5.2.3.4.
7.....4...2..7..8...3..8.799..5..3...6..2..9...1.97..6...3..9...3..4..6...9..1.35
....7..2.8.......6.1.2.5...9.54....8.........3....85.1...3.2.8.4.......9.7..6....
`
	s := New()
	for _, board := range strings.Split(hardest, "\n") {
		board = strings.TrimSpace(board)
		if len(board) > 0 {
			v, err := s.solveBoard(board)
			if err != nil {
				log.Fatal(err)
			}

			if !s.isSolved(v) {
				t.Errorf("not solved board %v", board)
			}
		}
	}
}

func BenchmarkSudokuNew(b *testing.B) {
	// Benchmarking initialization.
	for i := 0; i < b.N; i++ {
		bn := New()
		_ = bn
	}
}

func BenchmarkParseBoardAssign(b *testing.B) {
	// Benchmark how long it takes to parse a board and run full constraint
	// propagation. We know that for easyboard1 it's fully solved with
	// constraint propagation after parsing.
	bn := New()
	for i := 0; i < b.N; i++ {
		_, _ = bn.parseBoard(easyboard1)
	}
}

func BenchmarkSolveBoard(b *testing.B) {
	bn := New()
	for i := 0; i < b.N; i++ {
		_, _ = bn.solveBoard(hardboard2)
	}
}

func BenchmarkSolveBoardHardlong(b *testing.B) {
	bn := New()
	for i := 0; i < b.N; i++ {
		v, err := bn.solveBoard(hardlong)
		if err != nil || !bn.isSolved(v) {
			panic("not solved")
		}
	}
}
