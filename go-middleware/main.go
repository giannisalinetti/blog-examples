package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/user"
	"regexp"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const (
	rePattern = "^.*(curl|bot|[Pp]ython|[Ww]get).*"
)

// printWelcome prints a simple welcome message
func printWelcome(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering welcome funcion")
	fmt.Fprintf(w, "Hello visitor! You are connecting from IP/Port %s with User-Agent %s\n", r.RemoteAddr, r.UserAgent())
	log.Println("Leaving welcome function")
}

// printHelp prints out a simple help to demonstrate different handlers
func printHelp(w http.ResponseWriter, r *http.Request) {
	log.Println("Entering help funcion")
	fmt.Fprintf(w, "Usage: http://<host>:8000/welcome\n")
	log.Println("Leaving help function")
}

// userAgentCheck avoids requests from curl User-Agent
func userAgentCheck(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Begin User-Agent check")
	userAgent := r.UserAgent() // The User-Agent of the HTTP Request
	re, err := regexp.Compile(rePattern)
	if err != nil {
		log.Println("Error compiling regexp pattern")
		panic(1)
	}
	if re.MatchString(userAgent) {
		log.Printf("Refused connection to client with User-Agent %s", userAgent)
		fmt.Fprintf(w, "Error: cannot accept connections from %s User-Agent.\n", userAgent)
		return
	}
	log.Println("Completed User-Agent check")
	next(w, r)
}

// methodCheck verifies that the HTTP Request method is GET
func methodCheck(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("Begin HTTP Request Method check")
	if r.Method != "GET" {
		fmt.Fprintf(w, "Error: wrong http method\n")
		log.Println("Error during HTTP Method check")
		return // We don't need to go through middleware2
	}
	log.Println("Completed HTTP Request method check")
	next(w, r)
}

func main() {
	// Define port binding flag
	portFlag := flag.Int("port", 8000, "Port number")
	flag.Parse()

	// Load current username
	currentUser, err := user.Current()
	if err != nil {
		log.Fatal("Fatal: cannot evaluate current username")
	}

	// Test if current user is root to open ports below 1024
	if *portFlag < 1024 && currentUser.Username != "root" {
		log.Fatalf("Fatal: %s user cannot open ports under 1024\n", currentUser.Username)
	}

	// Create the port binding string
	portBinding := fmt.Sprintf(":%s", strconv.Itoa(*portFlag))

	// Create new mux router
	r := mux.NewRouter()

	// Register new router and associated handler functions
	r.HandleFunc("/welcome", printWelcome).Methods("GET")
	r.HandleFunc("/help", printHelp).Methods("GET")

	// Define a classic Negroni environment with a standard middleware stack:
	// Recovery - Panic Recovery Middleware
	// Logger - Request/Response Logging
	// Static - Static File Serving
	n := negroni.Classic()

	// Append middleware in the stack.
	// Functions are processed as Negroni Handlers in the same order they are passed.
	n.Use(negroni.HandlerFunc(userAgentCheck))
	n.Use(negroni.HandlerFunc(methodCheck))

	// Append a standard http.Handler
	n.UseHandler(r)

	// Negroni.Run is a wrapper to http.ListenAndServe
	n.Run(portBinding)
}
