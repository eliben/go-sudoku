package sudoku

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	f := Generate()
	fmt.Println(f)
	fmt.Println(Display(f))
}
