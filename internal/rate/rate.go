package rate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"time"
)

// TODO: rework this package & implement abstraction.

// CoinbaseResponse is the Coinbase API response structure.
type CoinbaseResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Base     string `json:"base"`
		Currency string `json:"currency"`
	} `json:"data"`
}

// GetRate returns the exchange rate between the base currency and the target currency using Coinbase API.
func GetRate(baseCode, targetCode string) (string, error) {
	if !isValidCurrencyCode(baseCode) || !isValidCurrencyCode(targetCode) {
		return "", fmt.Errorf("invalid currency code provided")
	}

	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/buy", baseCode, targetCode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
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

// isValidCurrencyCode validates whether the currency code is valid.
func isValidCurrencyCode(code string) bool {
	ok, _ := regexp.MatchString("^[A-Z]{3}$", code)
	return ok
}
