package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall/js"

	"github.com/eliben/go-sudoku"
)

func main() {
	fmt.Println("Console log: Go wasm")
	//fmt.Println(generateBoardSvg(30, false))
	js.Global().Set("generateBoard", jsonWrapper())
	<-make(chan bool)
}

func generateBoardSvg(hintCount int, symmetrical bool) string {
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
}

func jsonWrapper() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return fmt.Sprintf("got %v args, want 2", len(args))
		}
		hintCount := args[0].Int()
		symmetrical := args[1].Bool()
		brd := generateBoardSvg(hintCount, symmetrical)
		return brd
	})
	return jsonFunc
}
