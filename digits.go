package sudoku

import (
	"fmt"
	"math/bits"
	"strings"
)

// Digits represents a set of possible digits for a Sudoku square. The functions
// in this file perform set operations on Digits, as needed for Sudoku.
type Digits uint16

// Internal representation: digit N is represented by the Nth lowest bit in
// the Digits value, e.g.:
//
//    Digits = 0b0000_0000_0110_0010
//
// has the bits N=1,5,6 are set, so it represents the set of digits {1, 5, 6}

// fullDigitsSet returns a Digits with all possible digits set.
func fullDigitsSet() Digits {
	return 0b0000001111111110
}

// isMember checks whether digit n is a member of the digit set d.
func (d Digits) isMember(n uint16) bool {
	return (d & (1 << n)) != 0
}

// add adds digit n to set d and returns the new set.
func (d Digits) add(n uint16) Digits {
	return d | (1 << n)
}

// remove removes digit n from set d and returns the new set.
func (d Digits) remove(n uint16) Digits {
	return d &^ (1 << n)
}

// size returns the size of the set - the number of digits in it.
func (d Digits) size() int {
	return bits.OnesCount16(uint16(d))
}

// singleMemberDigit returns the digit that's a member of a 1-element set; this
// assumes that the set indeed has a single element.
func (d Digits) singleMemberDigit() uint16 {
	return uint16(bits.TrailingZeros16(uint16(d)))
}

// String implements the fmt.Stringer interface for Digits.
func (d Digits) String() string {
	var parts []string
	for i := uint16(1); i <= uint16(9); i++ {
		if d.isMember(i) {
			parts = append(parts, fmt.Sprintf("%v", i))
		}
	}
	return strings.Join(parts, "")
}
