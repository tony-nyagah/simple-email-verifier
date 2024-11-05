package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"text/template"

	emailverifier "github.com/AfterShip/email-verifier"
)

type EmailParts struct {
	Username string
	Domain   string
}

type VerificationResult struct {
	IsDeliverable bool
	HostExists    bool
	Disabled      bool
}

func splitEmail(email string) (EmailParts, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return EmailParts{}, fmt.Errorf("invalid email format")
	}

	return EmailParts{
		Username: parts[0],
		Domain:   parts[1],
	}, nil
}

func verifyEmail(domain string, username string, email string) (VerificationResult, error) {
	var verifier = emailverifier.NewVerifier().EnableSMTPCheck()

	result, err := verifier.CheckSMTP(domain, username)
	if err != nil {
		return VerificationResult{
			IsDeliverable: result.Deliverable,
			HostExists:    result.HostExists,
			Disabled:      result.Disabled,
		}, fmt.Errorf("failed to verify email: %v, %v", email, err)
	}

	return VerificationResult{
		IsDeliverable: result.Deliverable,
		HostExists:    result.HostExists,
		Disabled:      result.Disabled,
	}, nil
}

func main() {
	port := ":8037"
	log.Println("Starting server on port", strings.Split(port, ":")[1])

	// Styles
	styles := http.FileServer(http.Dir("./views/stylesheets"))
	http.Handle("/styles/", http.StripPrefix("/styles/", styles))

	// Health check
	healthCheck := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "It works!")
	}
	http.HandleFunc("/healthcheck", healthCheck)

	// Index page
	indexHandler := func(response http.ResponseWriter, request *http.Request) {
		tmpl := template.Must(template.ParseFiles("views/index.html"))
		tmpl.Execute(response, nil)
	}
	http.HandleFunc("/", indexHandler)

	log.Fatal(http.ListenAndServe(port, nil))

	// email := bufio.NewScanner(os.Stdin)
	// fmt.Print("Enter an email address: ")

	// if email.Scan() {
	// 	emailParts, err := splitEmail(email.Text())
	// 	if err != nil {
	// 		fmt.Println("Error: ", err)
	// 		return
	// 	}
	// 	fmt.Printf("Username: %s\nDomain: %s\n", emailParts.Username, emailParts.Domain)

	// 	fmt.Println(verifyEmail(emailParts.Domain, emailParts.Username, email.Text()))
	// } else {
	// 	fmt.Println("Error reading input: ", email.Err())
	// }
}
