package sudoku

import (
	"testing"
)

func TestIsSet(t *testing.T) {
	d := Digits(0b0000000001100010)

	var i uint16
	for i = 1; i <= 9; i++ {
		if i == 1 || i == 5 || i == 6 {
			if !isSet(d, i) {
				t.Errorf("got isSet(%016b, %d)=false, want true", d, i)
			}
		} else {
			if isSet(d, i) {
				t.Errorf("got isSet(%016b, %d)=true, want false", d, i)
			}
		}
	}
}
