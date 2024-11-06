# Email Verifier

A Go web application that verifies email addresses using multiple validation steps. Built with Go, HTMX, and Tailwind CSS.

## Features

- Real-time email verification
- SMTP server validation
- MX record checking
- Disposable email detection
- Modern UI with instant feedback

## How It Works

The email verification process consists of multiple steps:

1. **Format Validation**
   - Checks if the email follows the correct format (user@domain.com)
   - Validates the syntax according to RFC standards

2. **Domain Validation**
   - Verifies that the domain exists
   - Checks MX (Mail Exchange) records
   - Determines if it's a disposable email domain
   - Identifies role-based emails (e.g., admin@, support@, etc.)

3. **SMTP Verification**
   - Connects to the email server using SMTP
   - Simulates sending a message without actually sending one
   - Verifies if the mailbox exists for the specified address
   - Checks for catch-all policies

Note: Some mail servers may not cooperate with the SMTP verification process, which can affect accuracy.

## Installation

1. Install Go (version 1.21 or higher)
2. Clone this repository

## Running the Application

1. Start the server:
```bash
go run .
```

2. Open your browser and visit:
```
http://localhost:8080
```

## Dependencies

- [email-verifier](https://github.com/AfterShip/email-verifier): For comprehensive email verification
- HTMX: For dynamic UI updates
- Tailwind CSS: For styling

## Usage Tips

- The verification process may take a few seconds due to SMTP checks
- Some email providers may block SMTP verification attempts
- Results include:
  - Validity status
  - Detailed error messages
  - MX record availability
  - Disposable email detection
  - Reachability status

## License

MIT License
