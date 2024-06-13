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
func (m *Repository) SubscribeUser(email, baseCode, targetCode string) (err error, statusCode int) {
	// Validate email
	var emailValidator validator.EmailValidator

	if !emailValidator.Validate(email) {
		return errors.New("invalid email"), http.StatusBadRequest
	}

	// Create a user record (if not already created)
	user, err := m.App.Models.User.Create(email)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("already subscribed"), http.StatusConflict
		}

		return fmt.Errorf("error creating user"), http.StatusInternalServerError
	}

	// Get currency IDs
	baseCurrencyID, err := m.App.Models.Currency.GetIDbyCode(baseCode)
	if err != nil {
		return fmt.Errorf("error retrieving base currency"), http.StatusInternalServerError
	}

	targetCurrencyID, err := m.App.Models.Currency.GetIDbyCode(targetCode)
	if err != nil {
		return fmt.Errorf("error retrieving target currency"), http.StatusInternalServerError
	}

	// Create and save the subscription
	_, err = m.App.Models.Subscription.Create(user.ID, baseCurrencyID, targetCurrencyID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("already subscribed"), http.StatusConflict
		}

		return fmt.Errorf("error creating subscription"), http.StatusInternalServerError
	}

	return nil, http.StatusAccepted
}
