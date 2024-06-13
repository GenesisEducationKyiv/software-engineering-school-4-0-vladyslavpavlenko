package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rate"
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

// GetRate handles the `/rate` request.
func (m *Repository) GetRate(w http.ResponseWriter, _ *http.Request) {
	// Create a new Coinbase fetcher
	fetcher := rate.CoinbaseFetcher{
		Client: &http.Client{},
	}

	// Perform the fetching operation
	price, err := fetcher.FetchRate("USD", "UAH")
	if err != nil {
		_ = m.errorJSON(w, fmt.Errorf("error getting rate update: %w", err), http.StatusServiceUnavailable)
		return
	}

	// Create a response
	update := rateUpdate{
		BaseCode:   "USD",
		TargetCode: "UAH",
		Price:      price,
	}

	payload := jsonResponse{
		Error: false,
		Data:  update,
	}

	// Send the response back
	_ = m.writeJSON(w, http.StatusOK, payload)
}

// Subscribe handles the `/subscribe` request.
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	// Parse the form
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

	// Perform the subscription operation
	err, code := m.SubscribeUser(body.Email, "USD", "UAH")
	if err != nil {
		_ = m.errorJSON(w, err, code)
		return
	}

	// Create a response
	payload := jsonResponse{
		Error:   false,
		Message: "subscribed",
	}

	// Send the response back
	_ = m.writeJSON(w, http.StatusOK, payload)
}

// SendEmails handles the `/sendEmails` request.
func (m *Repository) SendEmails(w http.ResponseWriter, _ *http.Request) {
	// Perform the mailing operation
	err := m.NotifySubscribers()
	if err != nil {
		_ = m.errorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Create a response
	payload := jsonResponse{
		Error:   false,
		Message: "sent",
	}

	// Send the response back
	_ = m.writeJSON(w, http.StatusOK, payload)
}
