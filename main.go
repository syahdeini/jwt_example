package main

import (
	"log"
	"net/http"
)

func main() {
	// "Signin" and "welcome" are the handlers that will implement
	http.HandleFunc("/signin", Signin)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refersh", Refresh)

	// Start the server on port 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
