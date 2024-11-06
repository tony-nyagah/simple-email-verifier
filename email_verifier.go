package main

import (
	emailverifier "github.com/AfterShip/email-verifier"
)

type VerificationResult struct {
	IsValid      bool   `json:"isValid"`
	Message      string `json:"message"`
	MXAvailable  bool   `json:"mxAvailable"`
	Disposable   bool   `json:"disposable"`
	Reachable    string `json:"reachable"`
	FullResponse string `json:"fullResponse"`
}

var verifier = emailverifier.NewVerifier().
	EnableSMTPCheck().
	EnableDomainSuggest().
	EnableCatchAllCheck()

func verifyEmail(email string) VerificationResult {
	if email == "" {
		return VerificationResult{
			IsValid:   false,
			Message:   "Email cannot be empty",
			Reachable: "invalid",
		}
	}



	ret, err := verifier.Verify(email)
	if err != nil {
		return VerificationResult{
			IsValid:   false,
			Message:   "Verification failed: " + err.Error(),
			Reachable: "error",
		}
	}

	message := "Email appears valid"
	isValid := true
	reachable := "yes"

	if ret.Disposable {
		message = "Email is from a disposable service"
		isValid = false
		reachable = "disposable"
	}

	return VerificationResult{
		IsValid:    isValid,
		Message:    message,
		Disposable: ret.Disposable,
		Reachable:  reachable}
}
