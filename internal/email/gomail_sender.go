package email

import (
	"gopkg.in/gomail.v2"
)

// GomailSender implements the Sender interface for Gomail.
type GomailSender struct {
	Dialer Dialer
	Config Config
}

type GomailDialer struct {
	Dialer Dialer
}

func (d *GomailDialer) DialAndSend(m ...*gomail.Message) error {
	return d.Dialer.DialAndSend(m...)
}

func (gs *GomailSender) Send(params Params) error {
	m := gomail.NewMessage()
	m.SetHeader("From", gs.Config.Email)
	m.SetHeader("To", params.To)
	m.SetHeader("Subject", params.Subject)
	m.SetBody("text/plain", params.Body)

	if err := gs.Dialer.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
