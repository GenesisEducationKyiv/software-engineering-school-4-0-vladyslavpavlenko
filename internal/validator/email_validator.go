package validator

import (
	"net/mail"
)

type EmailValidator struct{}

// Validate checks if the given email address is valid.
func (v *EmailValidator) Validate(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
