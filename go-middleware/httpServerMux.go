package main

import (
	"fmt"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
}

func main() {
	// Define a custom multiplexer
	r := http.NewServeMux()

	// HandleFunc is now a method of the new multiplexer
	r.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
