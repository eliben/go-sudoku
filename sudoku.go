package sudoku

import (
	"fmt"
	"math/rand"
	"strings"

	"golang.org/x/exp/slices"
)

// Index represents a square on the Sudoku board; it's a number in the inclusive
// range [0, 80] that stands for row*9+col.
//
// These are the squares designated by an Index:
//
//  0  1  2 |  3  4  5 |  6  7  8
//  9 10 11 | 12 13 14 | 15 16 17
// 18 19 20 | 21 22 23 | 24 25 26
// ---------+----------+---------
// 27 28 29 | 30 31 32 | 33 34 35
// 36 37 38 | 39 40 41 | 42 43 44
// 45 46 47 | 48 49 50 | 51 52 53
// ---------+----------+---------
// 54 55 56 | 57 58 59 | 60 61 62
// 63 64 65 | 66 67 68 | 69 70 71
// 72 73 74 | 75 76 77 | 78 79 80
type Index = int

// Unit is a list of square indices that belong to the same Sudoku
// unit - a row, column or 3x3 block which should contain unique digits.
type Unit = []Index

// Values represents a Sudoku board in a format that's usable for solving.
// An element at index [i] in Values represents Sudoku square i (see the
// documentation of the Index type), and contains a set of all candidate
// digits for this square.
type Values []Digits

// unitlist is the list of all units that exist on the board.
var unitlist []Unit

// units maps an index to a list of units that contain that square.
// The mapping is a slice, i.e. units[i] is a list of all the units
// that contain the square with index i.
var units [][]Unit

// peers maps an index to a list of unique peers - other indices that share
// some unit with this index (it won't contain the index itself).
var peers [][]Index

func init() {
	index := func(row, col int) Index {
		return row*9 + col
	}

	// row units
	for row := 0; row < 9; row++ {
		var rowUnit []Index
		for col := 0; col < 9; col++ {
			rowUnit = append(rowUnit, index(row, col))
		}
		unitlist = append(unitlist, rowUnit)
	}

	// column units
	for col := 0; col < 9; col++ {
		var colUnit []Index
		for row := 0; row < 9; row++ {
			colUnit = append(colUnit, index(row, col))
		}
		unitlist = append(unitlist, colUnit)
	}

	// 3x3 block units
	for blockRow := 0; blockRow < 3; blockRow++ {
		for blockCol := 0; blockCol < 3; blockCol++ {
			var blockUnit []Index

			for row := 0; row < 3; row++ {
				for col := 0; col < 3; col++ {
					blockUnit = append(blockUnit, index(blockRow*3+row, blockCol*3+col))
				}
			}
			unitlist = append(unitlist, blockUnit)
		}
	}

	// For each index i, units[i] is a list of all units that contain i.
	units = make([][]Unit, 81)
	for i := 0; i < 81; i++ {
		for _, unit := range unitlist {
			if slices.Index(unit, i) >= 0 {
				units[i] = append(units[i], slices.Clone(unit))
			}
		}
	}

	// For each index i, peers[i] is a list of unique indices that share some
	// unit with i.
	peers = make([][]Index, 81)
	for i := 0; i < 81; i++ {
		for _, unit := range units[i] {
			for _, candidate := range unit {
				// This uses linear search to ensure uniqueness, but this calculation is
				// only done once at solver creation so we don't particularly care about
				// its speed.
				if candidate != i && slices.Index(peers[i], candidate) < 0 {
					peers[i] = append(peers[i], candidate)
				}
			}
		}
	}
}

// ParseBoard parses a Sudoku board given in textual representation, and returns
// it as Values. The textual representation is as described in
// http://norvig.com/sudoku.html: a string with a sequence of 81 runes in the
// set [0123456789.], where 0 or . mean "unassigned". All other runes in the
// string are ignored.
// If runElimination is false, the board is returned immediately after parsing.
// If runElimination is true, ParseBoard will invoke EliminateAll on the board
// and return the result. This is recommended when the board is then passed to
// a solver.
// It returns an error if there was an issue parsing the board, of if the board
// isn't a valid Sudoku board (e.g. contradictions exist).
func ParseBoard(str string, runElimination bool) (Values, error) {
	var dgs []uint16

	// Iterate and grab only the supported runes; ignore all others.
	for _, r := range str {
		if r >= '0' && r <= '9' {
			dgs = append(dgs, uint16(r)-uint16('0'))
		} else if r == '.' {
			dgs = append(dgs, 0)
		}
	}

	if len(dgs) != 81 {
		return nil, fmt.Errorf("got only %v digits in board, want 81", len(dgs))
	}

	// Start with an empty board.
	values := EmptyBoard()

	// Assign square digits based on the parsed board. Note that this runs
	// constraint propagation and may discover contradictions.
	for sq, d := range dgs {
		if d != 0 {
			values[sq] = SingleDigitSet(d)
		}
	}
	//fmt.Println(Display(values))

	if runElimination && !EliminateAll(values) {
		return nil, fmt.Errorf("contradiction when eliminating board")
	}
	//fmt.Println(Display(values))

	return values, nil
}

// EliminateAll runs elimination on all assigned squares in values. It applies
// first-order Sudoku heuristics on the entire board. Returns true if the
// elimination is successful, and false if the boards has a contradiction.
func EliminateAll(values Values) bool {
	for sq, d := range values {
		if d.Size() == 1 {
			// Because of how eliminate() works, we prepare for it by remembering
			// which digit this square has assigned, setting the square to the full
			// set of digits and then calling eliminate on all digits except the
			// assigned one.
			digit := d.SingleMemberDigit()
			values[sq] = FullDigitsSet()
			for dn := uint16(1); dn <= 9; dn++ {
				if dn != digit {
					if !eliminate(values, sq, dn) {
						return false
					}
				}
			}
		}
	}
	return true
}

// assign attempts to assign digit to values[square], propagating
// constraints from the assignment. values is modified.
// It returns true if the assignment succeeded, and false if the assignment
// fails resulting in an invalid Sudoku board.
func assign(values Values, square Index, digit uint16) bool {
	if EnableStats {
		Stats.NumAssigns++
	}

	for d := uint16(1); d <= 9; d++ {
		// For each d 1..9 that's != digit, if d is set in
		// values[square], try to eliminate it.
		if values[square].IsMember(d) && d != digit {
			if !eliminate(values, square, d) {
				return false
			}
		}
	}
	return true
}

// eliminate removes digit from the candidates in values[square], propagating
// constraints. values is modified.
// It returns false if this results in an invalid Sudoku board; otherwise
// returns true.
func eliminate(values Values, square Index, digit uint16) bool {
	if !values[square].IsMember(digit) {
		// Already eliminated
		return true
	}

	// Remove digit from the candidates in square.
	values[square] = values[square].Remove(digit)

	switch values[square].Size() {
	case 0:
		// No remaining options for square -- this is a contradiction.
		return false
	case 1:
		// A single digit candidate remaining in the square -- this creates a new
		// constraint. Eliminate this digit from all peer squares.
		remaining := values[square].SingleMemberDigit()
		for _, peer := range peers[square] {
			if !eliminate(values, peer, remaining) {
				return false
			}
		}
	}

	// Since digit was eliminated from square, it's possible that we'll find a
	// position for this digit in one of the units the square belongs to.
UnitLoop:
	for _, unit := range units[square] {
		// Looking for a single square in this unit that has 'digit' as one of its
		// candidates. sqd marks the square, or -1 if no such square was found.
		sqd := -1
		for _, sq := range unit {
			if values[sq].IsMember(digit) {
				if sqd == -1 {
					sqd = sq
				} else {
					// More than one square has 'digit' as a candidate, so we won't be
					// able to simplify things.
					continue UnitLoop
				}
			}
		}
		if sqd == -1 {
			// Contradiction: no places left in this unit for 'digit'
			return false
		}

		// There's only a single place left in the unit for 'digit' to go, so
		// assign it.
		if !assign(values, sqd, digit) {
			return false
		}
	}

	return true
}

// Display returns a visual representation of values, with all the digit
// candidates as a string in each cell.
func Display(values Values) string {
	// Find maximum length of one square.
	var maxlen int = 0
	for _, d := range values {
		if d.Size() > maxlen {
			maxlen = d.Size()
		}
	}
	width := maxlen + 1

	line := strings.Join([]string{
		strings.Repeat("-", width*3),
		strings.Repeat("-", width*3),
		strings.Repeat("-", width*3)}, "+")

	var sb strings.Builder
	for sq, d := range values {
		fmt.Fprintf(&sb, "%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+d.Size())/2, d))
		if sq%9 == 2 || sq%9 == 5 {
			sb.WriteString("|")
		}
		if sq%9 == 8 {
			sb.WriteRune('\n')
		}
		if sq == 26 || sq == 53 {
			sb.WriteString(line + "\n")
		}
	}
	return sb.String()
}

// DisplayAsInput returns a 2D Sudoku input board corresponding to values.
// It treats solved squares (with one candidate) as hints that are filled into
// the board, and unsolved squares (with more than one candidate) as empty.
func DisplayAsInput(values Values) string {
	line := strings.Join([]string{
		strings.Repeat("-", 6),
		strings.Repeat("-", 6),
		strings.Repeat("-", 6)}, "+")

	var sb strings.Builder
	for sq, d := range values {
		ds := d.String()
		if d.Size() > 1 {
			ds = "."
		}
		fmt.Fprintf(&sb, "%s ", ds)
		if sq%9 == 2 || sq%9 == 5 {
			sb.WriteString("|")
		}
		if sq%9 == 8 {
			sb.WriteRune('\n')
		}
		if sq == 26 || sq == 53 {
			sb.WriteString(line + "\n")
		}
	}
	return sb.String()
}

// EmptyBoard creates an "empty" Sudoku board, where each square can potentially
// contain any digit.
func EmptyBoard() Values {
	vals := make(Values, 81)
	for sq := range vals {
		vals[sq] = FullDigitsSet()
	}
	return vals
}

// IsSolved checks whether values is a properly solved Sudoku board, with all
// the constraints satisfied.
func IsSolved(values Values) bool {
	for _, unit := range unitlist {
		var dset Digits
		for _, sq := range unit {
			// Some squares have more than a single candidate? Not solved.
			if values[sq].Size() != 1 {
				return false
			}
			dset = dset.Add(values[sq].SingleMemberDigit())
		}
		// Not all digits covered by this unit? Not solved.
		if dset != FullDigitsSet() {
			return false
		}
	}
	return true
}

// findSquareWithFewestCandidates finds a square in values with more than one
// digit candidate, but the smallest number of such candidates.
func findSquareWithFewestCandidates(values Values) Index {
	var squareToTry Index = -1
	var minSize int = 10
	for sq, d := range values {
		if d.Size() > 1 && d.Size() < minSize {
			minSize = d.Size()
			squareToTry = sq
		}
	}
	return squareToTry
}

// SolveOptions is a container of options that can be taken by the
// Solve function.
type SolveOptions struct {
	// Randomize tells the solver to randomly shuffle its digit selection when
	// attempting to guess a value for a square. For actual randomness, the
	// rand package's default randomness source should be properly seeded before
	// invoking Solve.
	Randomize bool
}

// Solve runs a backtracking search to solve the board given in values.
// It returns true and the solved values if the search succeeded and we ended up
// with a board with only a single candidate per square; otherwise, it returns
// false. The input values is not modified.
// The solution process can be configured by providing SolveOptions.
// Consider making SolveOptions ... so they're not mandatory, but panic if more
// than 1.
func Solve(values Values, options SolveOptions) (Values, bool) {
	if EnableStats {
		Stats.NumSearches++
	}

	squareToTry := findSquareWithFewestCandidates(values)

	// If we didn't find any square with more than one candidate, the board is
	// solved!
	if squareToTry == -1 {
		return values, true
	}

	var candidates = []uint16{1, 2, 3, 4, 5, 6, 7, 8, 9}
	if options.Randomize {
		rand.Shuffle(len(candidates), func(i, j int) {
			candidates[i], candidates[j] = candidates[j], candidates[i]
		})
	}

	for _, d := range candidates {
		// Try to assign sq with each one of its candidate digits. If this results
		// in a successful Solve() - we've solved the board!
		if values[squareToTry].IsMember(d) {
			vcopy := slices.Clone(values)
			if assign(vcopy, squareToTry, d) {
				if vresult, solved := Solve(vcopy, options); solved {
					return vresult, true
				}
			}
		}
	}
	return values, false
}

// SolveAll finds all solutions to the given board and returns them. If no
// solutions were found, an empty list is returned. max can specify the
// (approximate) maximal number of solutions to find; a value <= 0 means "all of
// them". Often more solutions than max will be returned, but not a lot more
// (maybe 2-3x as many).
// Warning: this function can take a LONG time to run for boards with multiple
// solutions, and it can consume enormous amounts of memory because it has to
// remember each solution it finds. For some boards it will run forever (e.g.
// finding all solutions on an empty board). If in doubt, use the max parameter
// to restrict the number.
func SolveAll(values Values, max int) []Values {
	squareToTry := findSquareWithFewestCandidates(values)

	// If we didn't find any square with more than one candidate, the board is
	// solved!
	if squareToTry == -1 {
		return []Values{values}
	}

	var allSolved []Values

	for d := uint16(1); d <= 9; d++ {
		// Try to assign sq with each one of its candidate digits. If this results
		// in a successful Solve() - we've solved the board!
		if values[squareToTry].IsMember(d) {
			vcopy := slices.Clone(values)
			if assign(vcopy, squareToTry, d) {
				if vsolved := SolveAll(vcopy, max); len(vsolved) > 0 {
					allSolved = append(allSolved, vsolved...)
					if max > 0 && len(allSolved) >= max {
						return allSolved
					}
				}
			}
		}
	}
	return allSolved
}

// EnableStats enables statistics collection during the processes of solving.
// When stats are enabled, solving will be slightly slower.
//
// Note: statistics collection is NOT SAFE FOR CONCURRENT ACCESS.
var EnableStats bool = false

type StatsCollector struct {
	NumSearches uint64
	NumAssigns  uint64
}

// Stats is the global variable for accessing statistics from this package.
// It's recommended to call Stats.Reset() before solving a board, and access
// the Stats fields after it's done.
var Stats StatsCollector

func (s *StatsCollector) Reset() {
	s.NumSearches = 0
	s.NumAssigns = 0
}

// WithStats helps run any block of code with stats enabled.
func WithStats(f func()) {
	EnableStats = true
	defer func() {
		EnableStats = false
	}()
	Stats.Reset()

	f()
}
