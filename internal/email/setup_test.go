package email_test

import "github.com/vladyslavpavlenko/genesis-api-project/internal/email"

type MockEmailSender struct {
	SendFunc func(cfg email.Config, params email.Params) error
}

func (m MockEmailSender) Send(cfg email.Config, params email.Params) error {
	if m.SendFunc != nil {
		return m.SendFunc(cfg, params)
	}
	return nil
}
