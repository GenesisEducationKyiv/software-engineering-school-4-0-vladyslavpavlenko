package coinbase_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi/coinbase"
)

func TestFetchRate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := io.WriteString(w, `{"data":{"amount":"123.45","base":"USD","currency":"EUR"}}`)
		if err != nil {
			t.Log("Failed to write response body:", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer ts.Close()

	fetcher := coinbase.NewCoinbaseFetcher(&http.Client{})

	tests := []struct {
		name       string
		baseCode   string
		targetCode string
		wantErr    bool
	}{
		{"valid codes", "USD", "UAH", false},
		{"invalid base code", "US1", "UAH", true},
		{"invalid target code", "USD", "HRYVNIA", true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := fetcher.FetchRate(tc.baseCode, tc.targetCode)
			if (err != nil) != tc.wantErr {
				t.Errorf("FetchRate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if err == nil && got == "" {
				t.Errorf("FetchRate() = %v, want non-empty", got)
			}
		})
	}
}

func TestFetchRate_NetworkError(t *testing.T) {
	mockClient := &coinbase.MockHTTPClient{
		Err: fmt.Errorf("network error"),
	}
	fetcher := coinbase.Fetcher{Client: mockClient}
	_, err := fetcher.FetchRate("USD", "UAH")
	if err == nil {
		t.Errorf("Expected network error, got none")
	}
}

func TestFetchRate_NonOKStatusCode(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
	}
	mockClient := &coinbase.MockHTTPClient{
		Resp: mockResp,
	}
	fetcher := coinbase.Fetcher{Client: mockClient}
	_, err := fetcher.FetchRate("USD", "UAH")
	if err == nil {
		t.Errorf("Expected error for non-OK status code, got none")
	}
}

func TestFetchRate_ReadBodyError(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(errReader{}),
	}
	mockClient := &coinbase.MockHTTPClient{
		Resp: mockResp,
	}
	fetcher := coinbase.Fetcher{Client: mockClient}
	_, err := fetcher.FetchRate("USD", "EUR")
	if err == nil {
		t.Errorf("Expected error reading the body, got none")
	}
}

type errReader struct{}

func (e errReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestFetchRate_UnmarshalError(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"invalid_json"}`)),
	}
	mockClient := &coinbase.MockHTTPClient{
		Resp: mockResp,
	}
	fetcher := coinbase.Fetcher{Client: mockClient}
	_, err := fetcher.FetchRate("USD", "EUR")
	if err == nil {
		t.Errorf("Expected JSON unmarshal error, got none")
	}
}
