package email_test

//
// import (
//	"errors"
//	"sync"
//	"testing"
//
//	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
//)
//
// func TestSendEmail_InvalidEmail(_ *testing.T) {
//	var wg sync.WaitGroup
//	wg.Add(1)
//
//	cfg := email.Config{
//		Email:    "invalidemail",
//		Password: "password",
//	}
//
//	params := email.Params{
//		To:      "recipient@example.com",
//		Subject: "Test Subject",
//		Body:    "Test Body",
//	}
//
//	mockSender := MockEmailSender{
//		SendFunc: func(_ email.Config, _ email.Params) error {
//			return errors.New("invalid email address")
//		},
//	}
//
//	go email.SendEmail(mockSender, cfg, params)
//	wg.Wait()
// }
