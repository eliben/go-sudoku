package sudoku

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	f := Generate(25)
	fmt.Println(f)
	fmt.Println(Display(f))
	fmt.Println(DisplayAsInput(f))

	sols := SolveAll(f, -1)
	fmt.Println(len(sols))
}
