package handlers_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
)

func TestSubscribeUser(t *testing.T) {
	appConfig := &config.AppConfig{}
	mockFetcher := new(MockFetcher)
	mockSubscriber := new(MockSubscriber)
	mockSender := new(MockEmailSender)
	repo := handlers.NewRepo(appConfig, mockFetcher, mockSubscriber, mockSender)

	tests := []struct {
		name         string
		email        email.Email
		mockResponse error
		wantStatus   int
		wantErr      string
		expectCreate bool
	}{
		{"Invalid Email", "bademail", nil, http.StatusBadRequest, "invalid email", false},
		{"Valid Email Success", "good@email.com", nil, http.StatusAccepted, "", true},
		{"Duplicate Email", "duplicate@email.com", gorm.ErrDuplicatedKey, http.StatusConflict, "already subscribed", true},
		{"DB Error", "error@email.com", errors.New("db error"), http.StatusInternalServerError, "error creating user", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expectCreate {
				mockSubscriber.On("AddSubscription", string(tc.email)).Return(tc.mockResponse)
			}

			statusCode, err := repo.SubscribeUser(string(tc.email))
			assert.Equal(t, tc.wantStatus, statusCode)
			if err != nil {
				assert.Contains(t, err.Error(), tc.wantErr)
			}

			mockSubscriber.AssertExpectations(t)
		})
	}
}
