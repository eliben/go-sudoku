package sudoku

// Digits represents a set of possible digits for a Sudoku square. The functions
// in this file perform set operations on Digits, as needed for Sudoku.
type Digits = uint16

// Internal representation: digit N is represented by the Nth lowest bit in
// the Digits value, e.g.:
//
//    Digits = 0b0000_0000_0110_0010
//
// has the bits N=1,5,6 are set, so it represents the set of digits {1, 5, 6}

func isSet(d Digits, n uint16) bool {
	return (d & (1 << n)) != 0
}
