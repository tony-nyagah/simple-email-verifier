package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/verify", handleVerify)

	log.Println("Server starting on port 8052")
	log.Fatal(http.ListenAndServe(":8052", mux))
}
