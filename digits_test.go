package sudoku

import (
	"fmt"
	"testing"
)

func TestIsMember(t *testing.T) {
	d := Digits(0b0000000001100010)

	var i uint16
	for i = 1; i <= 9; i++ {
		if i == 1 || i == 5 || i == 6 {
			if !d.isMember(i) {
				t.Errorf("got isMember(%016b, %d)=false, want true", d, i)
			}
		} else {
			if d.isMember(i) {
				t.Errorf("got isMember(%016b, %d)=true, want false", d, i)
			}
		}
	}
}

func TestAdd(t *testing.T) {
	d := Digits(0b0000000001100010)

	if d.String() != "156" {
		t.Errorf("got %v, want 156", d.String())
	}

	wantd2 := Digits(0b0000000001100110)
	if d.add(2) != wantd2 {
		t.Errorf("got wantd2=%v, want=%v", d.add(2), wantd2)
	}

	wantd5 := Digits(0b0000000001100010)
	if d.add(5) != wantd5 {
		t.Errorf("got wantd5=%v, want=%v", d.add(5), wantd5)
	}
}

func TestRemove(t *testing.T) {
	d := Digits(0b0000000001100010)

	wantd5 := Digits(0b0000000001000010)
	if d.remove(5) != wantd5 {
		t.Errorf("got wantd5=%v, want=%v", d.remove(5), wantd5)
	}

	wantd8 := Digits(0b0000000001100010)
	if d.remove(8) != wantd8 {
		t.Errorf("got wantd8=%v, want=%v", d.remove(8), wantd8)
	}

	wantd6no1 := Digits(0b0000000001000000)
	donly6 := d.remove(1).remove(5)
	if donly6 != wantd6no1 {
		t.Errorf("got wantd6no1=%v, want=%v", donly6, wantd6no1)
	}

	if donly6.String() != "6" {
		t.Errorf("got %v, want 6", donly6.String())
	}
}

func TestRemoveAll(t *testing.T) {
	d := Digits(0b0011100110)

	got1 := d.removeAll(0b0011000000)
	want1 := Digits(0b0000100110)
	if got1 != want1 {
		t.Errorf("got %v, want %v", got1, want1)
	}

	got2 := d.removeAll(0b0001000110)
	want2 := Digits(0b0010100000)
	if got2 != want2 {
		t.Errorf("got %v, want %v", got2, want2)
	}
}

func TestAddRemoveAllSize(t *testing.T) {
	// Exhaustive testing that adds/removes every digits and tests that isMember
	// also keeps working.

	// Start with zero. Each iteration adds one digit, tests membership, then
	// removes the digit and tests again.
	d := Digits(0)

	testNoMembers := func() {
		if d.size() != 0 {
			t.Errorf("got size=%v, want 0", d.size())
		}

		for dig := uint16(1); dig <= 9; dig++ {
			if d.isMember(dig) {
				t.Errorf("got isMember=true for %v, want false", dig)
			}
		}
	}
	testNoMembers()

	for dig := uint16(1); dig <= 9; dig++ {
		t.Run(fmt.Sprintf("dig=%v", dig), func(t *testing.T) {
			// Add 'dig' to set
			d = d.add(dig)

			if d.size() != 1 {
				t.Errorf("got size=%v, want 1", d.size())
			}

			off := d.singleMemberDigit()
			if off != dig {
				t.Errorf("got singleMemberDigit=%v, want %v", off, dig)
			}

			// For each 'dig2', check set membership
			for dig2 := uint16(1); dig2 <= 9; dig2++ {
				if dig2 == dig {
					if !d.isMember(dig2) {
						t.Errorf("got isMember=false for %v, want true", dig2)
					}
				} else {
					if d.isMember(dig2) {
						t.Errorf("got isMember=true for %v, want false", dig2)
					}
				}
			}

			d = d.remove(dig)
			testNoMembers()
		})
	}
}

func TestTwoMemberDigits(t *testing.T) {
	d := Digits(0b0000000000100100)
	d1, d2 := d.twoMemberDigits()
	if d1 != 2 || d2 != 5 {
		t.Errorf("got %v,%v, want 2 and 5", d1, d2)
	}

	d = Digits(0b0000001000000010)
	d1, d2 = d.twoMemberDigits()
	if d1 != 1 || d2 != 9 {
		t.Errorf("got %v,%v, want 1 and 9", d1, d2)
	}
}

func TestFullDigitsSet(t *testing.T) {
	d := fullDigitsSet()
	for dig := uint16(1); dig <= 9; dig++ {
		if !d.isMember(dig) {
			t.Errorf("got isMember=false for %v, want true", dig)
		}
	}

	if d.String() != "123456789" {
		t.Errorf("got %v, want all digits", d.String())
	}
}

// TODO: add tests for intersection, or remove it
