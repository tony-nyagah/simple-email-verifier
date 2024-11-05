package main

import (
	"fmt"
	"strings"

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

func verifyEmail(domain string, username string) (VerificationResult, error) {
	verifier = emailverifier.NewVerifier().EnableSMTPCheck()

	result, err := verifier.CheckSMTP(domain, username)
	if err != nil {
		return VerificationResult{}, fmt.Errorf("failed to verify email: %v", err)
	}

	return VerificationResult{
		IsDeliverable: result.Deliverable,
		HostExists:    result.HostExists,
		Disabled:      result.Disabled,
	}, nil
}

func main() {

}
