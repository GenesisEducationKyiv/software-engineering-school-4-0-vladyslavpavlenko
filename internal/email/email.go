package email

import (
	"errors"
	"log"
	"net/mail"
	"sync"

	"gopkg.in/gomail.v2"
)

type Email string

// Validate checks if the given email address is valid.
func (e Email) Validate() bool {
	_, err := mail.ParseAddress(string(e))
	return err == nil
}

// Config holds the email configuration.
type Config struct {
	Email    string
	Password string
}

// NewEmailConfig creates an instance of Config.
func NewEmailConfig(email string, password string) (Config, error) {
	cfg := Config{
		Email:    email,
		Password: password,
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate validates an instance of Config to check whether the provided parameters are correct.
func (cfg *Config) Validate() error {
	if cfg.Email == "" {
		return errors.New("config email is empty")
	}

	if cfg.Password == "" {
		return errors.New("config password is empty")
	}

	if !Email(cfg.Email).Validate() {
		return errors.New("email is invalid")
	}

	return nil
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
