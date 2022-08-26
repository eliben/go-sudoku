package main

import (
	"fmt"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":8899", http.FileServer(http.Dir("cmd/wasm/assets")))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
