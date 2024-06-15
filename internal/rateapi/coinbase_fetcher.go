package rateapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// CoinbaseFetcher implements the Fetcher interface for Coinbase.
type CoinbaseFetcher struct {
	Client *http.Client
}

// CoinbaseResponse is the Coinbase API response structure.
type CoinbaseResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Base     string `json:"base"`
		Currency string `json:"currency"`
	} `json:"data"`
}

type Code string

// Validate checks if the given currency code conforms to the standard format, which consists of three uppercase letters.
func (c Code) Validate() bool {
	_, err := regexp.MatchString("^[A-Z]{3}$", string(c))
	return err == nil
}

func (f *CoinbaseFetcher) FetchRate(baseCode, targetCode string) (string, error) {
	if !Code(baseCode).Validate() || !Code(targetCode).Validate() {
		return "", fmt.Errorf("invalid currency code provided")
	}

	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/buy", baseCode, targetCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	resp, err := f.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body: %w", err)
	}

	var response CoinbaseResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling the response: %w, response body: %s", err, string(body))
	}

	return response.Data.Amount, nil
}
