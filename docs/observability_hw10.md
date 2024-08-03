# hw10-observability

For metrics storage, I used the VictoriaMetrics database alongside `github.com/VictoriaMetrics/metrics`, a lightweight and an easy-to-use package for exporting metrics in Prometheus format. For logging, I chose UBER’s zap package (`github.com/uber-go/zap`).

## Logs
After adding the `/metrics` endpoint to the main server, I began questioning how I can protect this endpoint from remaining openly accessible in the future and ensure it is accessible only by authorized services, such as the VictoriaMetrics scraper. As suggested by @vladpistun, I separated the **API** (`:8080`) and the **Metrics** (`:8081`) servers.
## Metrics
Before discussing the [alerts](#alerts) I would set up, let us first talk about the basic metrics. Below, you can see all the metrics I added so far, organized by packages.

### handlers
These are basic metrics that allow us to track the number of requests and the duration of those requests for each of the existing endpoints.

For the `/rate` endpoint:
```
"rate_requests_count"                   // counter
"rate_requests_duration_seconds"        // histogram
```

For the `/subscribe` endpoint:
```
"subscribe_requests_count"              // counter
"subscribe_requests_duration_seconds"   // histogram
```

For the `/unsubscribe` endpoint:
```
"unsubscribe_requests_count"            // counter
"unsubscribe_requests_duration_seconds" // histogram
```

For the `/sendEmails` endpoint:
```
"send_emails_requests_count"            // counter
"send_emails_requests_duration_seconds" // histogram
```

It might also be useful to set up metrics to track the spread of HTTP response codes for every endpoint.

### email/consumer
When dealing with sending emails, I think it is pretty much trivial yet useful to know the following: how many emails were sent, how many weren't, the success rate, and the time duration it takes to send an email. These are the exact four metrics I added to the `email/sender` package.
```
sent_emails_count              // counter
not_sent_emails_count          // counter
email_sending_duration_seconds // histogram
email_sending_success_rate     // gauge
```

This is the function that computes the success rate of email sending operations:
```go
func calculateEmailSuccessRate() float64 {
	totalAttempts := sentEmailsCounter.Get() + notSentEmailsCounter.Get()
	if totalAttempts == 0 {
		return 0
	}
	return float64(sentEmailsCounter.Get()) / float64(totalAttempts)
}
```

### rateapi
Having implemented the Chain of Responsibility pattern in the `rateapi` package, I thought it might be interesting (and useful) to track the number of requests made to each of the rate providers.

Therefore, I have set up these metrics:
```
fetched_via_coinbase_count // counter
fetched_via_privat_count   // counter
fetched_via_nbu_count      // counter
```

## 🚨 Alerts
Speaking of alerts, I would add them for the following metrics:

### [endpoint]_requests_count
If the count is 0 over a predefined interval (for example, 5 or 10 minutes), this suggests that there are no requests being made to the rate fetching service.

### (possible) [endpoint]_requests_400_count
As mentioned above, let us consider we have a metric that measures the HTTP response statuses for a requested endpoint. If more than, for example, 5% of requests result in a 400 error within a 10-minute frame, this might indicate that a significant portion of incoming requests are malformed or incorrect, which could be due to some issues with the service.

### sent_emails_count
If the number of emails sent is significantly lower than expected (for example, more than 20% decrease) compared to a rolling average, it might indicate an issue with email dispatch.

### email_sending_success_rate
If the success rate drops below a critical threshold (for example, below 90%) during email dispatch, this would indicate a high rate of email sending failures.

### fetched_via_[fetcher]_count
If the number of fetches via a certain rate fetcher reaches 0, then there must be an issue with this particular fetcher that requires maintenance. If there is a request limit for this particular service, then setting up an alert for approaching this limit is also a good idea.

## Conclusion
Somewhat unexpectedly for me, setting up the logging helped me, firstly, to identify some gaps in error handling, which I successfully fixed, and, secondly, pointed out places in the project that I should refactor a bit, which I did. 

In fact, the commits I've submitted to this pool request make a lot of small, yet significant, changes to this project.