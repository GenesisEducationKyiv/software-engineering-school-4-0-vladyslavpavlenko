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
	coinbaseURL = "https://api.coinbase.com/v2/prices/%s-%s/buy"
)

var fetchedViaCoinbaseCounter = metrics.NewCounter("fetched_via_coinbase_count")

type (
	CoinbaseFetcher struct {
		client HTTPClient
	}

	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}

	// coinbaseResponse is the Coinbase API response structure.
	coinbaseResponse struct {
		Data struct {
			Amount   string `json:"amount"`
			Base     string `json:"base"`
			Currency string `json:"currency"`
		} `json:"data"`
	}
)

// NewCoinbaseFetcher creates and returns a pointer to a new CoinbaseFetcher.
func NewCoinbaseFetcher(client HTTPClient) *CoinbaseFetcher {
	return &CoinbaseFetcher{client: client}
}

// Fetch performs a call to the Coinbase API to fetch the exchange rate between the
// specified base and target currencies.
func (f *CoinbaseFetcher) Fetch(ctx context.Context, base, target string) (string, error) {
	url := fmt.Sprintf(coinbaseURL, base, target)

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

	var r coinbaseResponse
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling the response: %w, response body: %s", err, string(body))
	}

	fetchedViaCoinbaseCounter.Inc()

	return r.Data.Amount, nil
}
