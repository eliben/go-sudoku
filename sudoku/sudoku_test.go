package sudoku

import (
	"testing"
)

func TestNew(t *testing.T) {
	s := New()
	if len(s.unitlist) != 27 {
		t.Errorf("got len=%v, want 27", len(s.unitlist))
	}
}
