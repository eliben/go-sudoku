package sudoku

import (
	"fmt"
	"log"
	"math/rand"
	"strings"
	"testing"
	"time"

	"golang.org/x/exp/slices"
)

func TestInit(t *testing.T) {
	// Smoke testing for the top-level vars initialized in init()
	if len(unitlist) != 27 {
		t.Errorf("got len=%v, want 27", len(unitlist))
	}

	wantUnits := []Unit{
		Unit{18, 19, 20, 21, 22, 23, 24, 25, 26},
		Unit{2, 11, 20, 29, 38, 47, 56, 65, 74},
		Unit{0, 1, 2, 9, 10, 11, 18, 19, 20}}

	if !slices.EqualFunc(wantUnits, units[20], func(a, b Unit) bool {
		return slices.Equal(a, b)
	}) {
		t.Errorf("got units[20]=%v\nwant %v", units[20], wantUnits)
	}

	gotPeers := peers[20]
	slices.Sort(gotPeers)
	wantPeers := []Index{0, 1, 2, 9, 10, 11, 18, 19, 21, 22, 23, 24, 25, 26, 29, 38, 47, 56, 65, 74}
	if !slices.Equal(wantPeers, gotPeers) {
		t.Errorf("got peers[20]=%v\n want %v", peers[20], wantPeers)
	}
}

func TestAssignElimination(t *testing.T) {
	vals := EmptyBoard()

	if IsSolved(vals) {
		t.Errorf("an empty board is solved")
	}

	// Assign a digit to square 20; check that this digit is the only candidate
	// in square 20, and that it was eliminated from all the peers of 20.
	assign(vals, 20, 5)

	if vals[20].Size() != 1 || vals[20].SingleMemberDigit() != 5 {
		t.Errorf("got vals[20]=%v", vals[20])
	}

	for sq := 0; sq <= 80; sq++ {
		if slices.Contains(peers[20], sq) {
			if vals[sq].IsMember(5) {
				t.Errorf("got member 5 in peer square %v", sq)
			}
		} else {
			if !vals[sq].IsMember(5) {
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

// This is the sudoku board Norvig reported takes the longest for his program
// to solve. Note that this board has multiple solutions.
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
	v, err := ParseBoard(easyboard1)
	if err != nil {
		t.Fatal(err)
	}

	if !IsSolved(v) {
		t.Errorf("expect easy board to be solved")
	}

	// Harder board that isn't fully solved without search.
	v2, err := ParseBoard(hardboard1)
	if err != nil {
		t.Fatal(err)
	}

	if IsSolved(v2) {
		t.Errorf("expect hard board to not be solved")
	}

	// Count how many squares are solved immediately in this puzzle and compare
	// to the number Norvig got.
	var solvedSquares int
	for _, d := range v2 {
		if d.Size() == 1 {
			solvedSquares++
		}
	}

	if solvedSquares != 20 {
		t.Errorf("got %v solved squares, want 20", solvedSquares)
	}
}

func TestSolveBoard(t *testing.T) {
	v, err := ParseBoard(hardboard1)
	if err != nil {
		log.Fatal(err)
	}
	vcopy := slices.Clone(v)
	vs, success := Solve(v, SolveOptions{})
	if !slices.Equal(v, vcopy) {
		t.Errorf("Solve modified board; before=%v, after=%v", vcopy, v)
	}

	if !success || !IsSolved(vs) {
		t.Errorf("expect hardboard1 to be solved by search")
	}

	// Should work on the easy board also (even though it's solved with the
	// initial parse)
	v2, err := ParseBoard(easyboard1)
	if err != nil {
		log.Fatal(err)
	}
	v2, success2 := Solve(v2, SolveOptions{})

	if !success2 || !IsSolved(v2) {
		t.Errorf("expect easy board to be solved by search")
	}

	// And the other hard board
	v3, err := ParseBoard(hardboard2)
	if err != nil {
		log.Fatal(err)
	}
	v3, success3 := Solve(v3, SolveOptions{})

	if !success3 || !IsSolved(v3) {
		t.Errorf("expect hardboard2 to be solved by search")
	}
}

func TestSolveWithStats(t *testing.T) {
	// The easy board is solved just by calling ParseBoard, needing no search.
	WithStats(func() {
		_, err := ParseBoard(easyboard1)
		if err != nil {
			t.Fatal(err)
		}

		if Stats.NumAssigns == 0 {
			t.Errorf("got NumAssigns==0")
		}
		if Stats.NumSearches != 0 {
			t.Errorf("got NumSearches=%v, want 0", Stats.NumSearches)
		}

		// For the hard board, we'll find both assigns and searches
		Stats.Reset()

		v, err := ParseBoard(hardboard1)
		if err != nil {
			t.Fatal(err)
		}
		_, _ = Solve(v, SolveOptions{})

		if Stats.NumAssigns == 0 {
			t.Errorf("got NumAssigns==0")
		}
		if Stats.NumSearches == 0 {
			t.Errorf("got NumSearches==0")
		}
	})
}

func TestIsSolved(t *testing.T) {
	v, err := ParseBoard(easyboard1)
	if err != nil {
		t.Fatal(err)
	}

	if !IsSolved(v) {
		t.Errorf("expect easy board to be solved")
	}

	// Now modify the board and make sure it's not considered "solved" any more.
	// ... modify by trying to add options to each square separately.
	for sq := range v {
		vcopy := slices.Clone(v)
		vcopy[sq] = vcopy[sq].Add(6).Add(8)

		if IsSolved(vcopy) {
			t.Errorf("expect board to not be solved after modification: %v", vcopy)
		}
	}
}

func TestSolveAll(t *testing.T) {
	v, err := ParseBoard(hardboard1)
	if err != nil {
		log.Fatal(err)
	}
	vs := SolveAll(v, -1)

	if len(vs) != 1 {
		t.Errorf("got %v solutions, want 1", len(vs))
	}
	if !IsSolved(vs[0]) {
		t.Errorf("got %v, want solved board", vs[0])
	}

	// Now generate a multiple-solution board, by replacing all instances
	// of 1 and 2 in the solved board by the digits set "12", permitting any
	// combination of them to solve the board.
	board := vs[0]
	for sq, d := range board {
		if d.Size() == 1 && (d.IsMember(1) || d.IsMember(2)) {
			board[sq] = d.Add(1).Add(2)
		}
	}

	vs = SolveAll(board, -1)
	if len(vs) != 2 {
		t.Errorf("got %v solved boards, want 2", len(vs))
	}

	if !IsSolved(vs[0]) || !IsSolved(vs[1]) {
		t.Errorf("got unsolved boards")
	}

	// Now try to limit max=1, see that SolveAll only returns a single solution.
	vs = SolveAll(board, 1)
	if len(vs) != 1 {
		t.Errorf("got %v solved boards, want 1", len(vs))
	}

	if !IsSolved(vs[0]) {
		t.Errorf("got unsolved boards")
	}

	// Create a board with a contradiction on purpose, and verify that SolveAll
	// returns an empty list.
	v, err = ParseBoard(hardboard1)
	if err != nil {
		log.Fatal(err)
	}
	v[30] = SingleDigitSet(1)
	v[31] = SingleDigitSet(2)
	v[32] = SingleDigitSet(3)
	vs = SolveAll(v, -1)
	if len(vs) != 0 {
		t.Errorf("expect unsolvable, got %v", vs)
	}
}

func TestHardlong(t *testing.T) {
	v, err := ParseBoard(hardlong)
	if err != nil {
		log.Fatal(err)
	}

	// The "hardlong" puzzle has multiple solutions. Norvig says he found 13
	// different solutions, but there are vastly more. Use SolveAll to explore
	// the first 1000 or so.

	// Find the first 1000 solutions
	vs := SolveAll(v, 1000)

	if len(vs) < 1000 {
		t.Errorf("got %v solutions, expected at least 1000", len(vs))
	}

	for _, v := range vs {
		if !IsSolved(v) {
			t.Errorf("got unsolved board %v", v)
		}
	}
}

// This board is unsolvable, but it takes the search a while to figure this
// out.
var impossible string = `
. . . |. . 5 |. 8 .
. . . |6 . 1 |. 4 3
. . . |. . . |. . .
------+------+------
. 1 . |5 . . |. . .
. . . |1 . 6 |. . .
3 . . |. . . |. . 5
------+------+------
5 3 . |. . . |. 6 1
. . . |. . . |. . 4
. . . |. . . |. . .`

// Run this test but skip in "not short" mode
func TestImpossible(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	WithStats(func() {
		v, err := ParseBoard(impossible)
		if err != nil {
			log.Fatal(err)
		}
		v, success := Solve(v, SolveOptions{})

		if success || IsSolved(v) {
			t.Errorf("got solved board for impossible")
		}
		fmt.Printf("searches=%v, assigns=%v\n", Stats.NumSearches, Stats.NumAssigns)
	})
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
	for _, board := range strings.Split(hardest, "\n") {
		board = strings.TrimSpace(board)
		if len(board) > 0 {
			v, err := ParseBoard(board)
			if err != nil {
				log.Fatalf("error for board %v: %v", board, err)
			}
			vs, success := Solve(v, SolveOptions{})
			if !success || !IsSolved(vs) {
				t.Errorf("not solved board %v", board)
			}

			// Now solve again, with randomization
			vsr, success := Solve(v, SolveOptions{Randomize: true})
			if !success || !IsSolved(vsr) {
				t.Errorf("not solved randomized board %v", board)
			}
		}
	}
}

func TestSolveEmpty(t *testing.T) {
	vals := EmptyBoard()
	vres, solved := Solve(vals, SolveOptions{})
	if !solved {
		t.Errorf("want Solve(empty) to report success")
	}

	if !IsSolved(vres) {
		t.Errorf("want solved result board; got:\n%v", Display(vres))
	}

	// Try a few randomized solutions
	for i := 0; i < 10; i++ {
		vs, solved := Solve(vals, SolveOptions{Randomize: true})
		if !solved {
			t.Errorf("want Solve(empty) to report success")
		}
		if !IsSolved(vs) {
			t.Errorf("want solved result board; got:\n%v", Display(vs))
		}
	}
}

func BenchmarkParseBoardAssign(b *testing.B) {
	// Benchmark how long it takes to parse a board and run full constraint
	// propagation. We know that for easyboard1 it's fully solved with
	// constraint propagation after parsing.
	for i := 0; i < b.N; i++ {
		_, _ = ParseBoard(easyboard1)
	}
}

func BenchmarkSolveBoardHardlong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		v, err := ParseBoard(hardlong)
		if err != nil {
			log.Fatal(err)
		}
		v, success := Solve(v, SolveOptions{})
		if !success {
			log.Fatal("not solved")
		}
	}
}

func BenchmarkSolveBoardHardlongRandomized(b *testing.B) {
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < b.N; i++ {
		v, err := ParseBoard(hardlong)
		if err != nil {
			log.Fatal(err)
		}
		v, success := Solve(v, SolveOptions{Randomize: true})
		if !success {
			log.Fatal("not solved")
		}
	}
}

func BenchmarkSolveEmpty(b *testing.B) {
	// Benchmark how long it takes to "solve" an empty board.
	empty := EmptyBoard()
	for i := 0; i < b.N; i++ {
		_, _ = Solve(empty, SolveOptions{})
	}
}

func BenchmarkSolveEmptyRandomized(b *testing.B) {
	// Benchmark how long it takes to "solve" an empty board,
	// with randomization. Each solution will be different.
	rand.Seed(time.Now().UnixNano())
	empty := EmptyBoard()
	for i := 0; i < b.N; i++ {
		_, _ = Solve(empty, SolveOptions{})
	}
}
