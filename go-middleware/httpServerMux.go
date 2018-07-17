package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func printHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!\n")
}

func printDate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Current date: %s\n", time.Now().String())
}

func main() {
	// Define a custom multiplexer
	r := http.NewServeMux()

	// Register printHello function
	r.HandleFunc("/hello", printHello)
	// Register printDate function
	r.HandleFunc("/date", printDate)

	// Use r as the default handler
	log.Fatal(http.ListenAndServe(":8080", r))
}
