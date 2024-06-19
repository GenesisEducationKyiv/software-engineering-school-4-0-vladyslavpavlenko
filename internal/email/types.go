package email

import (
	"net/mail"

	"gopkg.in/gomail.v2"
)

type Email string

// Validate checks if the given email address is valid.
func (e Email) Validate() bool {
	_, err := mail.ParseAddress(string(e))
	return err == nil
}

// Dialer defines an interface for a dialer to an SMTP server.
type Dialer interface {
	DialAndSend(m ...*gomail.Message) error
}

// Params holds the email message data.
type Params struct {
	To      string
	Subject string
	Body    string
}
