package main

import (
	"fmt"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

// validateEmail checks if the email is well-formed and has valid MX records.
func validateEmail(email string) (bool, string) {
	// Check basic structure
	_, err := mail.ParseAddress(email)
	if err != nil {
		return false, "Invalid email format"
	}

	// Extract domain
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false, "Invalid email format"
	}
	domain := parts[1]

	// Check MX records
	mxRecords, err := net.LookupMX(domain)
	if err != nil || len(mxRecords) == 0 {
		return false, fmt.Sprintf("No MX records found for domain: %s", domain)
	}

	// Check SMTP connection
	smtpHost := fmt.Sprintf("%s:%d", mxRecords[0].Host, 25)
	if err := smtpConnect(smtpHost, email); err != nil {
		return false, fmt.Sprintf("Failed to connect to SMTP server: %v", err)
	}

	return true, "Email is valid"
}

func smtpConnect(smtpHost string, email string) error {
	auth := smtp.PlainAuth("", email, "", strings.Split(smtpHost, ":")[0])
	conn, err := smtp.Dial(smtpHost)
	if err != nil {
		return err
	}
	defer conn.Close()

	if err := conn.Auth(auth); err != nil {
		return err
	}

	return nil
}
