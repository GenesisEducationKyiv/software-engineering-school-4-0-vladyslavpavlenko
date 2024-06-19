package rateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const URL = "https://api.coinbase.com/v2/prices/USD-UAH/buy"

// HTTPClient defines the interface for an HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// Fetcher implements the Fetcher interface for Coinbase.
type Fetcher struct {
	Client HTTPClient
}

// NewCoinbaseFetcher creates and returns a pointer to a new CoinbaseFetcher.
func NewCoinbaseFetcher(client HTTPClient) *Fetcher {
	return &Fetcher{Client: client}
}

// Response is the Coinbase API response structure.
type Response struct {
	Data struct {
		Amount   string `json:"amount"`
		Base     string `json:"base"`
		Currency string `json:"currency"`
	} `json:"data"`
}

func (f *Fetcher) Fetch() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", URL, http.NoBody)

	resp, err := f.Client.Do(req)
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

	var r Response
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling the response: %w, response body: %s", err, string(body))
	}

	return r.Data.Amount, nil
}
