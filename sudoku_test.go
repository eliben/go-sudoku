package sudoku

import (
	"testing"

	"golang.org/x/exp/slices"
)

func TestNew(t *testing.T) {
	// Smoke testing.
	s := New()
	if len(s.unitlist) != 27 {
		t.Errorf("got len=%v, want 27", len(s.unitlist))
	}

	wantUnits := []Unit{
		Unit{18, 19, 20, 21, 22, 23, 24, 25, 26},
		Unit{2, 11, 20, 29, 38, 47, 56, 65, 74},
		Unit{0, 1, 2, 9, 10, 11, 18, 19, 20}}

	if !slices.EqualFunc(wantUnits, s.units[20], func(a, b Unit) bool {
		return slices.Equal(a, b)
	}) {
		t.Errorf("got units[20]=%v\nwant %v", s.units[20], wantUnits)
	}

	gotPeers := s.peers[20]
	slices.Sort(gotPeers)
	wantPeers := []Index{0, 1, 2, 9, 10, 11, 18, 19, 21, 22, 23, 24, 25, 26, 29, 38, 47, 56, 65, 74}
	if !slices.Equal(wantPeers, gotPeers) {
		t.Errorf("got peers[20]=%v\n want %v", s.peers[20], wantPeers)
	}
}