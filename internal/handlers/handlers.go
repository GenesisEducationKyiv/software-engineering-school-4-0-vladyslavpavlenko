package handlers

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

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

const (
	subscribed   = "subscribed"
	unsubscribed = "unsubscribed"
)

var (
	errFetchingRate  = errors.New("failed to fetch rate")
	errSubscribing   = errors.New("failed to subscribe")
	errUnsubscribing = errors.New("failed to unsubscribe")
	errInvalidEmail  = errors.New("invalid email")
	errSendingEmails = errors.New("failed to send emails")
)

// GetRate handles the `/rateapi` request.
func (h *Handlers) GetRate(w http.ResponseWriter, r *http.Request) {
	// Perform the fetching operation
	price, err := h.Services.Fetcher.Fetch(r.Context(), "USD", "UAH")
	if err != nil {
		reqID := r.Context().Value(middleware.RequestIDKey).(string)
		h.l.Error("failed to fetch rate",
			zap.Error(err),
			zap.String("request_id", reqID),
		)

		_ = jsonutils.ErrorJSON(w, errFetchingRate, http.StatusServiceUnavailable)
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
func (h *Handlers) Subscribe(w http.ResponseWriter, r *http.Request) {
	h.handleSubscription(w, r, h.Services.Subscriber.AddSubscription, subscribed, errSubscribing)
}

// Unsubscribe handles the `/unsubscribe` request.
func (h *Handlers) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	h.handleSubscription(w, r, h.Services.Subscriber.DeleteSubscription, unsubscribed, errUnsubscribing)
}

// SendEmails handles the `/sendEmails` request.
func (h *Handlers) SendEmails(w http.ResponseWriter, r *http.Request) {
	// Produce mailing events
	err := h.Services.Notifier.Start()
	if err != nil {
		reqID := r.Context().Value(middleware.RequestIDKey).(string)
		h.l.Error("failed to send emails",
			zap.Error(err),
			zap.String("request_id", reqID),
		)

		_ = jsonutils.ErrorJSON(w, errSendingEmails, http.StatusInternalServerError)
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
func Metrics(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	metrics.WritePrometheus(w, false)
}

// handleSubscription processes subscription or unsubscription based on the provided handler function.
func (h *Handlers) handleSubscription(w http.ResponseWriter, r *http.Request, action func(string) error,
	successMessage string, errorMessage error,
) {
	email, err := parseEmail(r)
	if err != nil {
		h.handleError(w, r, err, http.StatusBadRequest, errInvalidEmail.Error())
		return
	}

	err = action(email)
	if err != nil {
		h.handleError(w, r, err, http.StatusInternalServerError, errorMessage.Error())
		return
	}

	payload := jsonutils.Response{
		Error:   false,
		Message: successMessage,
	}
	_ = jsonutils.WriteJSON(w, http.StatusOK, payload)
}

// handleError handles errors and logs them with the request ID.
func (h *Handlers) handleError(w http.ResponseWriter, r *http.Request, err error, statusCode int, logMessage string) {
	reqID, ok := r.Context().Value(middleware.RequestIDKey).(string)
	if !ok {
		reqID = "unknown"
	}
	h.l.Error(logMessage,
		zap.Error(err),
		zap.String("request_id", reqID),
	)

	_ = jsonutils.ErrorJSON(w, err, statusCode)
}

// parseEmail parses and validates the email from the multipart form of the http.Request.
func parseEmail(r *http.Request) (string, error) {
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
