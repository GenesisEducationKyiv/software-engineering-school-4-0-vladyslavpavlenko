package middleware

import (
	"net/http"
	"time"

	"github.com/VictoriaMetrics/metrics"
)

var (
	rateRequestCounter  = metrics.NewCounter("rate_request_count")
	rateRequestDuration = metrics.NewHistogram("rate_request_duration_seconds")
)

var (
	subscribeRequestCounter  = metrics.NewCounter("subscribe_request_count")
	subscribeRequestDuration = metrics.NewHistogram("subscribe_request_duration_seconds")
)

var (
	unsubscribeRequestCounter  = metrics.NewCounter("unsubscribe_request_count")
	unsubscribeRequestDuration = metrics.NewHistogram("unsubscribe_request_duration_seconds")
)

var (
	sendEmailsRequestCounter  = metrics.NewCounter("send_emails_request_count")
	sendEmailsRequestDuration = metrics.NewHistogram("send_emails_request_duration_seconds")
)

// Metrics is a middleware that records the duration and count of each request.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start).Seconds()

		switch r.URL.Path {
		case "/api/v1/rate":
			rateRequestCounter.Inc()
			rateRequestDuration.Update(duration)
		case "/api/v1/subscribe":
			subscribeRequestCounter.Inc()
			subscribeRequestDuration.Update(duration)
		case "/api/v1/unsubscribe":
			unsubscribeRequestCounter.Inc()
			unsubscribeRequestDuration.Update(duration)
		case "/api/v1/sendEmails":
			sendEmailsRequestCounter.Inc()
			sendEmailsRequestDuration.Update(duration)
		}
	})
}
