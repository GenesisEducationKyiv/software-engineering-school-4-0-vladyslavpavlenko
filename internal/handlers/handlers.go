package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers/json"
)

// rateUpdate holds the exchange rateapi update data.
type rateUpdate struct {
	BaseCode   string `json:"base_code"`
	TargetCode string `json:"target_code"`
	Price      string `json:"price"`
}

// GetRate handles the `/rateapi` request.
func (m *Repository) GetRate(w http.ResponseWriter, r *http.Request) {
	// Perform the fetching operation
	price, err := m.Services.Fetcher.Fetch(r.Context(), "USD", "UAH")
	if err != nil {
		_ = json.ErrorJSON(w, fmt.Errorf("error fetching rate update: %w", err), http.StatusServiceUnavailable)
		return
	}

	// AddSubscription a response
	payload := json.Response{
		Error: false,
		Data: rateUpdate{
			BaseCode:   "USD",
			TargetCode: "UAH",
			Price:      price,
		},
	}

	// Send the response back
	_ = json.WriteJSON(w, http.StatusOK, payload)
}

// subscriptionBody is the email subscription request body structure.
type subscriptionBody struct {
	Email string `json:"email"`
}

// Subscribe handles the `/subscribe` request.
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	var body subscriptionBody

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		_ = json.ErrorJSON(w, errors.New("failed to parse form"))
		return
	}

	body.Email = r.FormValue("email")
	if body.Email == "" {
		_ = json.ErrorJSON(w, errors.New("email is required"))
		return
	}

	// Perform the subscription operation
	code, err := m.SubscribeUser(body.Email)
	if err != nil {
		_ = json.ErrorJSON(w, err, code)
		return
	}

	// AddSubscription a response
	payload := json.Response{
		Error:   false,
		Message: "subscribed",
	}

	// Send the response back
	_ = json.WriteJSON(w, http.StatusOK, payload)
}

// SendEmails handles the `/sendEmails` request.
func (m *Repository) SendEmails(w http.ResponseWriter, _ *http.Request) {
	// Perform the mailing operation
	err := m.NotifySubscribers()
	if err != nil {
		_ = json.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// AddSubscription a response
	payload := json.Response{
		Error:   false,
		Message: "sent",
	}

	// Send the response back
	_ = json.WriteJSON(w, http.StatusOK, payload)
}
