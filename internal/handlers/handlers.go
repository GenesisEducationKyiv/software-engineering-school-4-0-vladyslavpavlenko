package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/utils"
)

// GetRate handles the `/rateapi` request.
func (m *Repository) GetRate(w http.ResponseWriter, _ *http.Request) {
	// Create a new Coinbase fetcher
	fetcher := rateapi.Fetcher{
		Client: &http.Client{},
	}

	// Perform the fetching operation
	price, err := fetcher.Fetch()
	if err != nil {
		_ = utils.ErrorJSON(w, fmt.Errorf("error fetching rate update: %w", err), http.StatusServiceUnavailable)
		return
	}

	// Create a response
	payload := utils.JSONResponse{
		Error: false,
		Data: rateUpdate{
			BaseCode:   "USD",
			TargetCode: "UAH",
			Price:      price,
		},
	}

	// Send the response back
	_ = utils.WriteJSON(w, http.StatusOK, payload)
}

// Subscribe handles the `/subscribe` request.
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	// Parse the form
	var body subscriptionBody

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		_ = utils.ErrorJSON(w, errors.New("failed to parse form"))
		return
	}

	body.Email = r.FormValue("email")
	if body.Email == "" {
		_ = utils.ErrorJSON(w, errors.New("email is required"))
		return
	}

	// Perform the subscription operation
	code, err := m.SubscribeUser(body.Email)
	if err != nil {
		_ = utils.ErrorJSON(w, err, code)
		return
	}

	// Create a response
	payload := utils.JSONResponse{
		Error:   false,
		Message: "subscribed",
	}

	// Send the response back
	_ = utils.WriteJSON(w, http.StatusOK, payload)
}

// SendEmails handles the `/sendEmails` request.
func (m *Repository) SendEmails(w http.ResponseWriter, _ *http.Request) {
	// Perform the mailing operation
	err := m.NotifySubscribers()
	if err != nil {
		_ = utils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// Create a response
	payload := utils.JSONResponse{
		Error:   false,
		Message: "sent",
	}

	// Send the response back
	_ = utils.WriteJSON(w, http.StatusOK, payload)
}
