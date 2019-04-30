// +build !js

package main

import (
	"log"
	"net/http"
)

// Serve the wasm ui binary.
//	The file main.wasm must be in the current directory.
//	go run serve.go
func main() {
	println("serving on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", http.FileServer(http.Dir("."))))
}
