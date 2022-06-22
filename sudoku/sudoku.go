package sudoku

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
	unitlist []Unit
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

	return &Sudoku{
		unitlist: unitlist}
}
