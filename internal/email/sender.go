package email

import (
	"log"
	"sync"

	"gopkg.in/gomail.v2"
)

// Sender defines an interface for sending emails.
type Sender interface {
	Send(cfg Config, params Params) error
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

// SendEmail sends an email using the provided configuration and message data.
func SendEmail(wg *sync.WaitGroup, sender Sender, cfg Config, params Params) {
	defer wg.Done()

	if err := sender.Send(cfg, params); err != nil {
		log.Printf("Could not send email to %s: %v", params.To, err)
	}
}
