package middleware

import (
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

var (
	rateRequestsCounter  = metrics.NewCounter("rate_requests_count")
	rateRequestsDuration = metrics.NewHistogram("rate_requests_duration_seconds")
)

var (
	subscribeRequestsCounter  = metrics.NewCounter("subscribe_requests_count")
	subscribeRequestsDuration = metrics.NewHistogram("subscribe_requests_duration_seconds")
)

var (
	unsubscribeRequestsCounter  = metrics.NewCounter("unsubscribe_requests_count")
	unsubscribeRequestsDuration = metrics.NewHistogram("unsubscribe_requests_duration_seconds")
)

var (
	sendEmailsRequestsCounter  = metrics.NewCounter("send_emails_requests_count")
	sendEmailsRequestsDuration = metrics.NewHistogram("send_emails_requests_duration_seconds")
)

// Metrics is a middleware that records the duration and count of each request.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		switch r.URL.Path {
		case "/api/v1/rate":
			rateRequestsCounter.Inc()
			rateRequestsDuration.UpdateDuration(start)
		case "/api/v1/subscribe":
			subscribeRequestsCounter.Inc()
			subscribeRequestsDuration.UpdateDuration(start)
		case "/api/v1/unsubscribe":
			unsubscribeRequestsCounter.Inc()
			unsubscribeRequestsDuration.UpdateDuration(start)
		case "/api/v1/sendEmails":
			sendEmailsRequestsCounter.Inc()
			sendEmailsRequestsDuration.UpdateDuration(start)
		}
	})
}
