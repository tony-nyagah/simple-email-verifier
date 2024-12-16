package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Static file server for frontend assets
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// Define routes
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/validate-email", validateEmailHandler)

	// Start the server
	port := "8080"
	fmt.Printf("Server running on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
