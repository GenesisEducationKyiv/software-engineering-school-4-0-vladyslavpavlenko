package email

import "net/mail"

type Email string

// Validate checks if the given email address is valid.
func (e Email) Validate() bool {
	_, err := mail.ParseAddress(string(e))
	return err == nil
}
