package main

import (
	"fmt"
	"net/http"
)

func main() {
	port := ":8899"
	fmt.Println("Serving on", port)

	err := http.ListenAndServe(port, http.FileServer(http.Dir("assets")))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
