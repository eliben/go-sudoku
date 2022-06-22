package sudoku

import "github.com/eliben/go-sudoku/slices"

// Digits represents a set of possible digits for a Sudoku square.
type Digits = uint16

// Index represents a square on the Sudoku board; it's a number in the inclusive
// range [0, 80] that stands for row*9+col.
type Index = int

// Unit is a list of square indices that belong to the same Sudoku
// unit - a row, column or 3x3 block which should contain unique digits.
type Unit = []Index

func index(row, col int) Index {
	return row*9 + col
}

type Sudoku struct {
	// unitlist is the list of all units that exist on the board.
	unitlist []Unit

	// units maps an index to a list of units that contain that square.
	// The mapping is a slice, i.e. units[i] is a list of all the units
	// that contain the square with index i.
	units [][]Unit
}

func New() *Sudoku {
	var unitlist []Unit

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
	units := make([][]Unit, 81)
	for i := 0; i < 81; i++ {
		for _, unit := range unitlist {
			if slices.Index(unit, i) >= 0 {
				units[i] = append(units[i], slices.Clone(unit))
			}
		}
	}

	return &Sudoku{
		unitlist: unitlist,
		units:    units,
	}
}
