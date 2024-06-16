package email_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email/mocks"
	"gopkg.in/gomail.v2"
)

func TestSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDialer := mocks.NewMockDialer(ctrl)
	gomailSender := email.GomailSender{Dialer: mockDialer}
	config := email.Config{Email: "test@example.com", Password: "password"}
	params := email.Params{To: "recipient@example.com", Subject: "Test", Body: "Hello"}

	mockDialer.EXPECT().DialAndSend(gomock.Any()).Return(nil)

	err := gomailSender.Send(config, params)
	if err != nil {
		t.Errorf("Send failed: %v", err)
	}
}

func TestGomailSenderSendFailure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDialer := mocks.NewMockDialer(ctrl)
	sender := email.GomailSender{Dialer: mockDialer}

	cfg := email.Config{Email: "test@example.com", Password: "password"}
	params := email.Params{To: "recipient@example.com", Subject: "Failure Test", Body: "This email should encounter a send error."}

	testError := errors.New("smtp error")
	mockDialer.EXPECT().DialAndSend(gomock.Any()).Return(testError)

	err := sender.Send(cfg, params)

	assert.Equal(t, testError, err, "Expected a specific error, but got a different one")
}

func TestGomailDialerDialAndSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRealDialer := mocks.NewMockDialer(ctrl)
	gomailDialer := &email.GomailDialer{Dialer: mockRealDialer}

	msg := gomail.NewMessage()
	msg.SetHeader("From", "sender@example.com")
	msg.SetHeader("To", "receiver@example.com")
	msg.SetHeader("Subject", "Test Email")
	msg.SetBody("text/plain", "This is a test email.")

	// Set expectation
	mockRealDialer.EXPECT().DialAndSend(gomock.Any()).Return(nil)

	// Execute the method
	err := gomailDialer.DialAndSend(msg)
	if err != nil {
		t.Errorf("DialAndSend failed: %v", err)
	}
}
