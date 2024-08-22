//go:build js && wasm

package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/eliben/go-sudoku"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	fmt.Println("go-sudoku wasm")

	// Export the jsGenerateBoard function to JS.
	js.Global().Set("generateBoard", jsGenerateBoard)

	// For the Go code to be usable from JS, the main function has to run forever.
	<-make(chan bool)
}

// jsGenerateBoard wraps the functionality we need from this package, for use
// in the web interface. It creates a function that takes two parameters:
// an integer hint count, and a boolean "is symmetrical" flag. It returns
// the SVG generated for the board as a string.
var jsGenerateBoard = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
	if len(args) != 2 {
		return fmt.Sprintf("got %v args, want 2", len(args))
	}
	hintCount := args[0].Int()
	symmetrical := args[1].Bool()

	var board sudoku.Values
	if symmetrical {
		board = sudoku.GenerateSymmetrical(hintCount)
	} else {
		board = sudoku.Generate(hintCount)
	}

	d, err := sudoku.EvaluateDifficulty(board)
	if err != nil {
		log.Fatal(err)
	}

	var buf bytes.Buffer
	sudoku.DisplayAsSVG(&buf, board, d)
	return buf.String()
})
