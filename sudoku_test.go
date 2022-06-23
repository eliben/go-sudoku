package sudoku

import (
	"log"
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
		t.Errorf("expect hard board to be solved by search")
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
		_, _ = bn.solveBoard(hardboard1)
	}
}
