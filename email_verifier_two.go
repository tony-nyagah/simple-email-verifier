package main

import (
	"errors"
	"fmt"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
)

type VerificationResult struct {
	IsValid       bool
	IsDeliverable bool
	IsDisabled    bool
	FullInbox     bool
	HostExists    bool
	Message       string
}

type EmailParts struct {
	Username string
	Domain   string
}

var verifier = emailverifier.NewVerifier().EnableSMTPCheck()

func splitEmail(email string) (EmailParts, error) {
	parts := strings.Split(email, "@")

	if len(parts) != 2 {
		return EmailParts{}, errors.New("invalid email format")
	}

	return EmailParts{
		Username: parts[0],
		Domain:   parts[1],
	}, nil
}

func checkEmail(domain string, username string) (VerificationResult, error) {
	result, err := verifier.CheckSMTP(domain, username)

	if err != nil {
		return VerificationResult{Message: "verification failed: " + err.Error()}, errors.New("verification failed: " + err.Error())
	}

	return VerificationResult{
		IsDeliverable: result.Deliverable,
		IsDisabled:    result.Disabled,
		FullInbox:     result.FullInbox,
		HostExists:    result.HostExists,
	}, nil
}

func main() {
	email := "a.nyagah@gmail.com"
	emailParts, err := splitEmail(email)
	if err != nil {
		panic(err)
	}
	result, err := checkEmail(emailParts.Domain, emailParts.Username)
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Is deliverable: %t\n", result.IsDeliverable)
	fmt.Printf("Is disabled: %t\n", result.IsDisabled)
	fmt.Printf("Full inbox: %t\n", result.FullInbox)
	fmt.Printf("Host exists: %t\n", result.HostExists)
}
