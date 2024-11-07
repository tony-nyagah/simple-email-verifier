package main

import (
	"errors"
	"fmt"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
)

// VerificationResult represents the result of an email verification.
type VerificationResult struct {
	IsValid       bool
	IsDeliverable bool
	IsDisabled    bool
	FullInbox     bool
	HostExists    bool
	Message       string
}

// EmailVerifier is responsible for verifying email addresses.
type EmailVerifier struct {
	verifier *emailverifier.Verifier
}

// NewEmailVerifier returns a new EmailVerifier instance.
func NewEmailVerifier() *EmailVerifier {
	return &EmailVerifier{
		verifier: emailverifier.NewVerifier().EnableSMTPCheck(),
	}
}

// VerifyEmail verifies an email address and returns the result.
func (ev *EmailVerifier) VerifyEmail(email string) (VerificationResult, error) {
	emailParts := strings.Split(email, "@")

	if len(emailParts) != 2 {
		return VerificationResult{}, errors.New("invalid email format")
	}

	username := emailParts[0]
	domain := emailParts[1]

	result, err := ev.verifier.CheckSMTP(domain, username)
	if err != nil {
		return VerificationResult{}, fmt.Errorf("verification failed: %w", err)
	}

	verificationResult := VerificationResult{
		IsValid:       result.Deliverable,
		IsDeliverable: result.Deliverable,
		IsDisabled:    result.Disabled,
		FullInbox:     result.FullInbox,
		HostExists:    result.HostExists,
	}

	if verificationResult.IsValid {
		verificationResult.Message = "Email is deliverable. The email address is valid"
	} else {
		verificationResult.Message = "Email is not deliverable. The email address is invalid"
	}

	return verificationResult, nil
}
