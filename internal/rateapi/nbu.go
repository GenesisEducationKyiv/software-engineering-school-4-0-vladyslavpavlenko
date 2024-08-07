package rateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/VictoriaMetrics/metrics"
)

const (
	nbuURL = "https://bank.gov.ua/NBUStatService/v1/statdirectory/exchange?valcode=%s&json"
)

var fetchedViaNBUCounter = metrics.NewCounter("fetched_via_nbu_count")

type (
	NBUFetcher struct {
		client HTTPClient
	}

	// nbuResponse is the bank.gov.ua API response structure.
	nbuResponse struct {
		R030         int     `json:"r030"`
		TXT          string  `json:"txt"`
		Rate         float64 `json:"rate"`
		CC           string  `json:"cc"`
		ExchangeData string  `json:"exchange_data"`
	}
)

// NewNBUFetcher creates and returns a pointer to a new NBUFetcher.
func NewNBUFetcher(client HTTPClient) *NBUFetcher {
	return &NBUFetcher{client: client}
}

// Fetch performs a call to the https://bank.gov.ua to fetch the exchange rate between the
// specified base currency and UAH.
func (f *NBUFetcher) Fetch(ctx context.Context, base, _ string) (string, error) {
	url := fmt.Sprintf(nbuURL, base)

	req, _ := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)

	resp, err := f.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var r []nbuResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling the response: %w, response body: %s", err, string(body))
	}

	if len(r) == 0 {
		return "", fmt.Errorf("no data in response")
	}

	rate := fmt.Sprintf("%f", r[0].Rate)

	fetchedViaNBUCounter.Inc()
	return rate, nil
}
