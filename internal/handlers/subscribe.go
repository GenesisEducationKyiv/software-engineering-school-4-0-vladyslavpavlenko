package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"gorm.io/gorm"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

// SubscribeUser adds a user to the `subscriptions` table.
func (m *Repository) SubscribeUser(emailAddr string) (statusCode int, err error) {
	if !email.Email(emailAddr).Validate() {
		return http.StatusBadRequest, errors.New("invalid email")
	}

	err = m.Services.Subscriber.AddSubscription(emailAddr)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return http.StatusConflict, fmt.Errorf("already subscribed")
		}

		return http.StatusInternalServerError, err
	}

	return http.StatusAccepted, nil
}
