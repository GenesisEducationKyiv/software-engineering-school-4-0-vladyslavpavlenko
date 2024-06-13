package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/validator"
	"gorm.io/gorm"
)

// SubscribeUser subscribes a user to the rate update mailing list by adding a new email to the database and
// creating a corresponding subscription record. TODO: remove hardcoded currency codes.
func (m *Repository) SubscribeUser(email, baseCode, targetCode string) (statusCode int, err error) {
	// Validate email
	var emailValidator validator.EmailValidator

	if !emailValidator.Validate(email) {
		return http.StatusBadRequest, errors.New("invalid email")
	}

	// Create a user record (if not already created)
	user, err := m.App.Models.User.Create(email)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return http.StatusConflict, fmt.Errorf("already subscribed")
		}

		return http.StatusInternalServerError, fmt.Errorf("error creating user")
	}

	// Get currency IDs
	baseCurrencyID, err := m.App.Models.Currency.GetIDbyCode(baseCode)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error retrieving base currency")
	}

	targetCurrencyID, err := m.App.Models.Currency.GetIDbyCode(targetCode)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("error retrieving target currency")
	}

	// Create and save the subscription
	_, err = m.App.Models.Subscription.Create(user.ID, baseCurrencyID, targetCurrencyID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return http.StatusConflict, fmt.Errorf("already subscribed")
		}

		return http.StatusInternalServerError, fmt.Errorf("error creating subscription")
	}

	return http.StatusAccepted, nil
}
