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

func TestPrivatFetcher_Fetch(t *testing.T) {
	tests := []struct {
		name          string
		mockDoFunc    func(req *http.Request) (*http.Response, error)
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
							"ccy": "EUR",
							"base_ccy": "UAH",
							"buy": "24.5",
							"sale": "25.0"
						},
						{
							"ccy": "USD",
							"base_ccy": "UAH",
							"buy": "24.5",
							"sale": "25.0"
						}
					]`))),
				}
				return resp, nil
			},
			expectedRate:  "24.5",
			expectedError: "",
		},
		{
			name: "error making request",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				return nil, errors.New("error making request")
			},
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
			expectedRate:  "",
			expectedError: "no data in response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockHTTPClient{
				DoFunc: tt.mockDoFunc,
			}
			fetcher := rateapi.NewPrivatFetcher(client)
			rate, err := fetcher.Fetch(context.Background(), "", "")

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
