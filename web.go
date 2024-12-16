package main

import (
	"fmt"
	"net/http"
)

// homeHandler serves the main HTML page
func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}

// validateEmailHandler processes email validation requests
func validateEmailHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	valid, message := validateEmail(email)
	w.Header().Set("Content-Type", "text/html")
	class := "invalid"
	if valid {
		class = "valid"
	}
	fmt.Fprintf(w, `<div id="result" class="%s">%s</div>`, class, message)
}
