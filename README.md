# Simple Email Verifier
## How do we verify an email ?
* First it checks for email address format.
* Then make sure that domain name is valid. We also check whether it’s a disposable email address or not.
* In the final step, It extracts the MX records from the domain records and connects to the email server (over SMTP and also simulates sending a message) to make sure the mailbox really exists for that user/address. Some mail servers do not cooperate in the process, in such cases, the result of this email verification tool may not be as accurate as expected.

## How do we use it ?
