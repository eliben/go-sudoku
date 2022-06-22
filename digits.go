package sudoku

// Digits represents a set of possible digits for a Sudoku square. The functions
// in this file perform set operations on Digits, as needed for Sudoku.
type Digits uint16

// Internal representation: digit N is represented by the Nth lowest bit in
// the Digits value, e.g.:
//
//    Digits = 0b0000_0000_0110_0010
//
// has the bits N=1,5,6 are set, so it represents the set of digits {1, 5, 6}

// isMember checks whether digit n is a member of the digit set d.
func (d Digits) isMember(n uint16) bool {
	return (d & (1 << n)) != 0
}

// add adds digit n to set d and returns the new set.
func (d Digits) add(n uint16) Digits {
	return d | (1 << n)
}
