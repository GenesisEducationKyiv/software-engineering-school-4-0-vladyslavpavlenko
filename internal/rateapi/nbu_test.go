package rateapi_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
)

func TestNBUFetcher_Fetch(t *testing.T) {
	tests := []struct {
		name          string
		mockDoFunc    func(req *http.Request) (*http.Response, error)
		base          string
		expectedRate  string
		expectedError string
	}{
		{
			name: "successful response",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader([]byte(`[
						{
							"r030": 840,
							"txt": "Долар США",
							"rate": 40.4478,
							"cc": "USD",
							"exchange_data": "24.06.2024"
						}
					]`))),
				}
				return resp, nil
			},
			base:          "USD",
			expectedRate:  "40.447800",
			expectedError: "",
		},
		{
			name: "error making request",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				return nil, errors.New("error making request")
			},
			base:          "USD",
			expectedRate:  "",
			expectedError: "error making request: error making request",
		},
		{
			name: "error reading response body",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(&errorReader{}),
				}
				return resp, nil
			},
			base:          "USD",
			expectedRate:  "",
			expectedError: "error reading the response body: read error",
		},
		{
			name: "API request failed with status",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusInternalServerError,
					Body:       io.NopCloser(bytes.NewReader([]byte("internal server error"))),
				}
				return resp, nil
			},
			base:          "USD",
			expectedRate:  "",
			expectedError: "API request failed with status 500: internal server error",
		},
		{
			name: "error unmarshalling response",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
				}
				return resp, nil
			},
			base:          "USD",
			expectedRate:  "",
			expectedError: "error unmarshaling the response: invalid character 'i' looking for beginning of value, response body: invalid json",
		},
		{
			name: "no data in response",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("[]"))),
				}
				return resp, nil
			},
			base:          "USD",
			expectedRate:  "",
			expectedError: "no data in response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockHTTPClient{
				DoFunc: tt.mockDoFunc,
			}
			fetcher := rateapi.NewNBUFetcher(client)
			rate, err := fetcher.Fetch(context.Background(), tt.base, "")

			if rate != tt.expectedRate {
				t.Errorf("expected rate %s, got %s", tt.expectedRate, rate)
			}

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("expected error %v, got %v", tt.expectedError, err)
			} else if err == nil && tt.expectedError != "" {
				t.Errorf("expected error %v, got none", tt.expectedError)
			}
		})
	}
}
