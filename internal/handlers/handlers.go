package handlers

import (
	"errors"
	"fmt"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rate"
	"gorm.io/gorm"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

// rateUpdate holds the exchange rate update data.
type rateUpdate struct {
	BaseCode   string `json:"base_code"`
	TargetCode string `json:"target_code"`
	Price      string `json:"price"`
}

// subscriptionBody is the email subscription request body structure.
type subscriptionBody struct {
	Email string `json:"email"`
	// BaseCurrencyCode   string `json:"base_currency_code"`
	// TargetCurrencyCode string `json:"target_currency_code"`
}

// GetRate gets the current USD to UAH exchange rate.
func (m *Repository) GetRate(w http.ResponseWriter, r *http.Request) {
	price, err := rate.GetRate("USD", "UAH")
	if err != nil {
		_ = m.errorJSON(w, errors.New("error calling Coinbase API"), http.StatusBadRequest) // http.StatusServiceUnavailable
		return
	}

	rateResp := rateUpdate{
		BaseCode:   "USD",
		TargetCode: "UAH",
		Price:      price,
	}

	// Send response
	payload := jsonResponse{
		Error: false,
		Data:  rateResp,
	}

	_ = m.writeJSON(w, http.StatusOK, payload)
}

// Subscribe handles email subscriptions by adding a new email to the database and creating a corresponding subscription
// record. By default, it sets up a USD to UAH exchange rate subscription, but the implementation allows for working
// with any currency, as long as it exists in the `currencies` table.
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	var body subscriptionBody

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		_ = m.errorJSON(w, errors.New("failed to parse form"))
		return
	}

	body.Email = r.FormValue("email")
	if body.Email == "" {
		_ = m.errorJSON(w, errors.New("email is required"))
		return
	}

	if !validateEmail(body.Email) {
		_ = m.errorJSON(w, errors.New("email is invalid"))
		return
	}

	// Create and save the user
	user, err := m.App.Models.User.Create(body.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			_ = m.errorJSON(w, fmt.Errorf("already subscribed"), http.StatusConflict)
		} else {
			_ = m.errorJSON(w, fmt.Errorf("error creating user"), http.StatusInternalServerError)
		}
		return
	}

	// Get currency IDs
	baseCurrencyID, err := m.getCurrencyID("USD")
	if err != nil {
		_ = m.errorJSON(w, fmt.Errorf("error retrieving base currency"))
		return
	}

	targetCurrencyID, err := m.getCurrencyID("UAH")
	if err != nil {
		_ = m.errorJSON(w, fmt.Errorf("error retrieving target currency"))
		return
	}

	// Create and save the subscription
	_, err = m.App.Models.Subscription.Create(user.ID, baseCurrencyID, targetCurrencyID)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			_ = m.errorJSON(w, fmt.Errorf("already subscribed"), http.StatusConflict)
		} else {
			_ = m.errorJSON(w, fmt.Errorf("error creating subscription"), http.StatusInternalServerError)
		}
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "subscribed",
	}

	_ = m.writeJSON(w, http.StatusOK, payload)
}

// SendEmails handles the /sendEmails request.
func (m *Repository) SendEmails(w http.ResponseWriter, r *http.Request) {
	m.NotifySubscribers()
}
