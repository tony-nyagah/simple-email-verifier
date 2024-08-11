package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	emailverifier "github.com/AfterShip/email-verifier"
)

type EmailParts struct {
	Username string
	Domain   string
}

var verifier = emailverifier.NewVerifier().EnableSMTPCheck()

func splitEmail(email string) (EmailParts, error) {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return EmailParts{}, fmt.Errorf("Invalid email format!")
	}
	return EmailParts{
		Username: parts[0],
		Domain:   parts[1],
	}, nil
}

func checkEmail(domain string, username string) {
	result, err := verifier.CheckSMTP(domain, username)
	if err != nil {
		fmt.Println("Check SMTP failed: ", err)
		return
	}

	fmt.Println("SMTP validation result: ", result)
}

func main() {
	email := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter an email address: ")

	if email.Scan() {
		emailParts, err := splitEmail(email.Text())
		if err != nil {
			fmt.Println("Error: ", err)
			return
		}
		fmt.Printf("Username: %s\nDomain: %s\n", emailParts.Username, emailParts.Domain)

		checkEmail(emailParts.Domain, emailParts.Username)
	} else {
		fmt.Println("Error reading input: ", email.Err())
	}
}
