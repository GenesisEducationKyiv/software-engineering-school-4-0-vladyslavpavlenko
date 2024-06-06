package email

import (
	"gopkg.in/gomail.v2"
	"log"
	"sync"
)

// Config holds the email configuration.
type Config struct {
	Email    string
	Password string
}

// Params holds the email message data.
type Params struct {
	To      string
	Subject string
	Body    string
}

// SendEmail sends an email using the provided configuration and message data.
func SendEmail(wg *sync.WaitGroup, config Config, params Params) {
	defer wg.Done()

	msg := Params{
		To:      params.To,
		Subject: params.Subject,
		Body:    params.Body,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Email)
	m.SetHeader("To", msg.To)
	m.SetHeader("Subject", msg.Subject)
	m.SetBody("text/plain", msg.Body)

	d := gomail.NewDialer("smtp.gmail.com", 587, config.Email, config.Password)

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Could not send email to %s: %v", msg.To, err)
	} else {
		log.Printf("Email sent successfully to %s!", msg.To)
	}
}
