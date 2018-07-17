package main

import (
	"fmt"
	"net/http"
)

func middleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Executing middleware before mainLogic")
		handler.ServeHTTP(w, r)
		fmt.Println("Executing middleware after mainLogic")
	})
}

func mainLogic(w http.ResponseWriter, r *http.Request) {
	// Business logic goes here
	fmt.Println("Executing mainLogic")
	w.Write([]byte("May the Force be with you!\n"))
}

func main() {
	// HandlerFunc returns a HTTP Handler
	mainLogicHandler := http.HandlerFunc(mainLogic)
	http.Handle("/", middleware(mainLogicHandler))
	http.ListenAndServe(":8000", nil)
}
