package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
	emailpkg "github.com/vladyslavpavlenko/genesis-api-project/internal/email"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/jsonutils"
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
		_ = jsonutils.ErrorJSON(w, fmt.Errorf("error fetching rate update: %w", err), http.StatusServiceUnavailable)
		return
	}

	// AddSubscription a response
	payload := jsonutils.Response{
		Error: false,
		Data: rateUpdate{
			BaseCode:   "USD",
			TargetCode: "UAH",
			Price:      price,
		},
	}

	// Send the response back
	_ = jsonutils.WriteJSON(w, http.StatusOK, payload)
}

// Subscribe handles the `/subscribe` request.
func (m *Repository) Subscribe(w http.ResponseWriter, r *http.Request) {
	email, err := parseEmailFromRequest(r)
	if err != nil {
		_ = jsonutils.ErrorJSON(w, err)
		return
	}

	err = m.Services.Subscriber.AddSubscription(email)
	if err != nil {
		_ = jsonutils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonutils.Response{
		Error:   false,
		Message: "subscribed",
	}

	_ = jsonutils.WriteJSON(w, http.StatusOK, payload)
}

// Unsubscribe handles the `/unsubscribe` request.
func (m *Repository) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	email, err := parseEmailFromRequest(r)
	if err != nil {
		_ = jsonutils.ErrorJSON(w, err)
		return
	}

	err = m.Services.Subscriber.DeleteSubscription(email)
	if err != nil {
		_ = jsonutils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	payload := jsonutils.Response{
		Error:   false,
		Message: "unsubscribed",
	}

	_ = jsonutils.WriteJSON(w, http.StatusOK, payload)
}

// SendEmails handles the `/sendEmails` request.
func (m *Repository) SendEmails(w http.ResponseWriter, _ *http.Request) {
	// Produce mailing events
	err := m.Services.Notifier.Start()
	if err != nil {
		_ = jsonutils.ErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	// AddSubscription a response
	payload := jsonutils.Response{
		Error: false,
	}

	// Send the response back
	_ = jsonutils.WriteJSON(w, http.StatusOK, payload)
}

// Metrics serves the application metrics in the Prometheus format.
func (m *Repository) Metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	metrics.WritePrometheus(w, false)
}

// parseEmailFromRequest parses the email from the multipart form and validates it.
func parseEmailFromRequest(r *http.Request) (string, error) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return "", errors.New("failed to parse form")
	}

	emailAddr := r.FormValue("email")
	if emailAddr == "" {
		return "", errors.New("email is required")
	}

	if !emailpkg.Email(emailAddr).Validate() {
		return "", errors.New("invalid email")
	}

	return emailAddr, nil
}
