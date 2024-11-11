package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"
)

// Global template cache
var (
	templates   *template.Template
	validator   *Validator
	configOnce  sync.Once
	templateDir = "templates" // Update this to match your template directory
)

// Initialize templates and validator
func init() {
	// Initialize templates
	var err error
	templates, err = template.ParseFiles(
		filepath.Join(templateDir, "index.html"),
		filepath.Join(templateDir, "result.html"),
	)
	if err != nil {
		log.Fatalf("Error parsing templates: %v", err)
	}

	// Initialize validator with default config
	validator = NewValidator(DefaultConfig())
}

func handleHome(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func handleVerify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	email := r.FormValue("email")
	if email == "" {
		renderResult(w, VerificationResult{
			IsValid: false,
			Message: "Email address is required",
		})
		return
	}

	// Validate the email
	result, err := validator.ValidateEmail(email)
	if err != nil {
		// Handle validation errors gracefully
		result = VerificationResult{
			IsValid: false,
			Message: err.Error(),
		}
	}

	renderResult(w, result)
}

func renderResult(w http.ResponseWriter, result VerificationResult) {
	err := templates.ExecuteTemplate(w, "result.html", result)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func main() {
	mux := http.NewServeMux()

	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Routes
	mux.HandleFunc("/", handleHome)
	mux.HandleFunc("/verify", handleVerify)

	// Start server
	log.Println("Server starting on port 8052")
	log.Fatal(http.ListenAndServe(":8052", mux))
}
