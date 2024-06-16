package rateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient defines the interface for an HTTP client.
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// CoinbaseFetcher implements the Fetcher interface for Coinbase.
type CoinbaseFetcher struct {
	Client HTTPClient
}

// MockHTTPClient defines the interface for a mock HTTP client.
type MockHTTPClient struct {
	Resp *http.Response
	Err  error
}

func (m *MockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return m.Resp, m.Err
}

// CoinbaseResponse is the Coinbase API response structure.
type CoinbaseResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Base     string `json:"base"`
		Currency string `json:"currency"`
	} `json:"data"`
}

func (f *CoinbaseFetcher) FetchRate(baseCode, targetCode string) (string, error) {
	if !Code(baseCode).Validate() || !Code(targetCode).Validate() {
		return "", fmt.Errorf("invalid currency code provided")
	}

	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/buy", baseCode, targetCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, _ := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)

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

	var response CoinbaseResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling the response: %w, response body: %s", err, string(body))
	}

	return response.Data.Amount, nil
}
