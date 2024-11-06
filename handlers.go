package main

import (
	"html/template"
	"net/http"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	result := verifyEmail(email)

	tmpl := template.Must(template.ParseFiles("templates/result.html"))
	tmpl.Execute(w, result)
}
