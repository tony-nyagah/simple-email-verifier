// Package emailvalidator provides robust email validation functionality with SMTP verification
package main

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	emailverifier "github.com/AfterShip/email-verifier"
)

var (
	// ErrInvalidEmailFormat indicates the email string doesn't match expected format
	ErrInvalidEmailFormat = errors.New("invalid email format: must contain exactly one @ symbol")
	// ErrVerificationFailed indicates SMTP verification failed
	ErrVerificationFailed = errors.New("smtp verification failed")
)

// Config holds the configuration for the email validator
type Config struct {
	EnableSMTP    bool
	RetryAttempts int
	Timeout       int // in seconds
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		EnableSMTP:    true,
		RetryAttempts: 2,
		Timeout:       10,
	}
}

// Validator provides thread-safe email validation functionality
type Validator struct {
	verifier *emailverifier.Verifier
	config   *Config
	mu       sync.RWMutex
}

// VerificationResult represents the outcome of an email verification attempt
type VerificationResult struct {
	IsValid       bool   `json:"is_valid"`
	IsDeliverable bool   `json:"is_deliverable"`
	IsDisabled    bool   `json:"is_disabled"`
	FullInbox     bool   `json:"full_inbox"`
	HostExists    bool   `json:"host_exists"`
	Message       string `json:"message,omitempty"`
}

// EmailParts represents the constituent parts of an email address
type EmailParts struct {
	Username string
	Domain   string
}

// NewValidator creates a new email validator with the provided configuration
func NewValidator(config *Config) *Validator {
	if config == nil {
		config = DefaultConfig()
	}

	v := &Validator{
		config: config,
	}

	verifierOpts := emailverifier.NewVerifier()
	if config.EnableSMTP {
		verifierOpts.EnableSMTPCheck()
	}
	if config.Timeout > 0 {
		verifierOpts.SetTimeout(config.Timeout)
	}

	v.verifier = verifierOpts

	return v
}

// SplitEmail breaks an email address into its username and domain parts
func (v *Validator) SplitEmail(email string) (EmailParts, error) {
	parts := strings.Split(strings.TrimSpace(email), "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return EmailParts{}, ErrInvalidEmailFormat
	}

	return EmailParts{
		Username: parts[0],
		Domain:   parts[1],
	}, nil
}

// CheckEmail performs comprehensive email validation including SMTP verification
func (v *Validator) CheckEmail(domain, username string) (VerificationResult, error) {
	v.mu.RLock()
	defer v.mu.RUnlock()

	var result VerificationResult
	var lastErr error

	// Retry logic for flaky SMTP connections
	for attempt := 0; attempt < v.config.RetryAttempts; attempt++ {
		smtpResult, err := v.verifier.CheckSMTP(domain, username)
		if err == nil {
			return VerificationResult{
				IsValid:       true,
				IsDeliverable: smtpResult.Deliverable,
				IsDisabled:    smtpResult.Disabled,
				FullInbox:     smtpResult.FullInbox,
				HostExists:    smtpResult.HostExists,
			}, nil
		}
		lastErr = err
	}

	if lastErr != nil {
		return VerificationResult{
			Message: fmt.Sprintf("verification failed: %v", lastErr),
		}, fmt.Errorf("%w: %v", ErrVerificationFailed, lastErr)
	}

	return result, nil
}

// ValidateEmail is a convenience method that combines splitting and checking
func (v *Validator) ValidateEmail(email string) (VerificationResult, error) {
	parts, err := v.SplitEmail(email)
	if err != nil {
		return VerificationResult{}, err
	}
	return v.CheckEmail(parts.Domain, parts.Username)
}
