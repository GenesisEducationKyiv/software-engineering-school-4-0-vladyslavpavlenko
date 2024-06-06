package rate

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CoinbaseResponse is the Coinbase API response structure.
type CoinbaseResponse struct {
	Data struct {
		Amount   string `json:"amount"`
		Base     string `json:"base"`
		Currency string `json:"currency"`
	} `json:"data"`
}

// GetRate returns the exchange rate between the base currency and the target currency using Coinbase API.
func GetRate(baseCode string, targetCode string) (string, error) {
	url := fmt.Sprintf("https://api.coinbase.com/v2/prices/%s-%s/buy", baseCode, targetCode)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading the response body: %v", err)
	}

	var response CoinbaseResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("error unmarshaling the response: %v", err)
	}

	return response.Data.Amount, nil
}
