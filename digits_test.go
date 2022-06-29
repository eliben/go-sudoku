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
			if !d.IsMember(i) {
				t.Errorf("got isMember(%016b, %d)=false, want true", d, i)
			}
		} else {
			if d.IsMember(i) {
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
	if d.Add(2) != wantd2 {
		t.Errorf("got wantd2=%v, want=%v", d.Add(2), wantd2)
	}

	wantd5 := Digits(0b0000000001100010)
	if d.Add(5) != wantd5 {
		t.Errorf("got wantd5=%v, want=%v", d.Add(5), wantd5)
	}
}

func TestRemove(t *testing.T) {
	d := Digits(0b0000000001100010)

	wantd5 := Digits(0b0000000001000010)
	if d.Remove(5) != wantd5 {
		t.Errorf("got wantd5=%v, want=%v", d.Remove(5), wantd5)
	}

	wantd8 := Digits(0b0000000001100010)
	if d.Remove(8) != wantd8 {
		t.Errorf("got wantd8=%v, want=%v", d.Remove(8), wantd8)
	}

	wantd6no1 := Digits(0b0000000001000000)
	donly6 := d.Remove(1).Remove(5)
	if donly6 != wantd6no1 {
		t.Errorf("got wantd6no1=%v, want=%v", donly6, wantd6no1)
	}

	if donly6.String() != "6" {
		t.Errorf("got %v, want 6", donly6.String())
	}
}

func TestRemoveAll(t *testing.T) {
	d := Digits(0b0011100110)

	got1 := d.RemoveAll(0b0011000000)
	want1 := Digits(0b0000100110)
	if got1 != want1 {
		t.Errorf("got %v, want %v", got1, want1)
	}

	got2 := d.RemoveAll(0b0001000110)
	want2 := Digits(0b0010100000)
	if got2 != want2 {
		t.Errorf("got %v, want %v", got2, want2)
	}
}

func TestAddRemoveAllSize(t *testing.T) {
	// Exhaustive testing that adds/removes every digits and tests that IsMember
	// also keeps working.

	// Start with zero. Each iteration adds one digit, tests membership, then
	// removes the digit and tests again.
	d := Digits(0)

	testNoMembers := func() {
		if d.Size() != 0 {
			t.Errorf("got size=%v, want 0", d.Size())
		}

		for dig := uint16(1); dig <= 9; dig++ {
			if d.IsMember(dig) {
				t.Errorf("got IsMember=true for %v, want false", dig)
			}
		}
	}
	testNoMembers()

	for dig := uint16(1); dig <= 9; dig++ {
		t.Run(fmt.Sprintf("dig=%v", dig), func(t *testing.T) {
			// Add 'dig' to set
			d = d.Add(dig)

			if d.Size() != 1 {
				t.Errorf("got size=%v, want 1", d.Size())
			}

			off := d.SingleMemberDigit()
			if off != dig {
				t.Errorf("got SingleMemberDigit=%v, want %v", off, dig)
			}

			// For each 'dig2', check set membership
			for dig2 := uint16(1); dig2 <= 9; dig2++ {
				if dig2 == dig {
					if !d.IsMember(dig2) {
						t.Errorf("got IsMember=false for %v, want true", dig2)
					}
				} else {
					if d.IsMember(dig2) {
						t.Errorf("got IsMember=true for %v, want false", dig2)
					}
				}
			}

			d = d.Remove(dig)
			testNoMembers()
		})
	}
}

func TestFullDigitsSet(t *testing.T) {
	d := FullDigitsSet()
	for dig := uint16(1); dig <= 9; dig++ {
		if !d.IsMember(dig) {
			t.Errorf("got IsMember=false for %v, want true", dig)
		}
	}

	if d.String() != "123456789" {
		t.Errorf("got %v, want all digits", d.String())
	}
}

func TestSingleDigitSet(t *testing.T) {
	d := SingleDigitSet(5)
	for dig := uint16(1); dig <= 9; dig++ {
		if dig == 5 {
			if !d.IsMember(dig) {
				t.Errorf("got IsMember=false for 5, want true")
			}
		} else {
			if d.IsMember(dig) {
				t.Errorf("got IsMember=true for %v, want false", dig)
			}
		}
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
