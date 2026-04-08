//go:build !js

package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Serving on http://localhost:8080")
	if err := http.ListenAndServe(":8080", http.FileServer(http.Dir("."))); err != nil {
		fmt.Println("Server error:", err)
	}
}
