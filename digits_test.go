package sudoku

import (
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

	wantd2 := Digits(0b0000000001100110)
	if d.add(2) != wantd2 {
		t.Errorf("got wantd2=%v, want=%v", d.add(2), wantd2)
	}

	wantd5 := Digits(0b0000000001100010)
	if d.add(5) != wantd5 {
		t.Errorf("got wantd5=%v, want=%v", d.add(5), wantd5)
	}
}
