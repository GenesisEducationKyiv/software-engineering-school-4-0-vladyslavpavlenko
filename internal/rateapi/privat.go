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
	privatURL = "https://api.privatbank.ua/p24api/pubinfo?json&exchange&coursid=5"
)

var fetchedViaPrivat = metrics.NewCounter("fetched_via_privat")

type (
	PrivatFetcher struct {
		client HTTPClient
	}

	// privatResponse is the api.privatbank.ua API response structure.
	privatResponse struct {
		CCY     string `json:"ccy"`
		BaseCCY string `json:"base_ccy"`
		Buy     string `json:"buy"`
		Sale    string `json:"sale"`
	}
)

// NewPrivatFetcher creates and returns a pointer to a new PrivatFetcher.
func NewPrivatFetcher(client HTTPClient) *PrivatFetcher {
	return &PrivatFetcher{client: client}
}

// Fetch performs a call to the https://api.privatbank.ua to fetch the exchange rate between USD and UAH.
func (f *PrivatFetcher) Fetch(ctx context.Context, _, _ string) (string, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", privatURL, http.NoBody)

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

	var r []privatResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling the response: %w, response body: %s", err, string(body))
	}

	if len(r) == 0 {
		return "", fmt.Errorf("no data in response")
	}

	fetchedViaPrivat.Inc()
	return r[1].Buy, nil
}
