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

type (
	MockHTTPClient struct {
		DoFunc func(req *http.Request) (*http.Response, error)
	}

	errorReader struct{}
)

func (e *errorReader) Read(_ []byte) (n int, err error) {
	return 0, errors.New("read error")
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestCoinbaseFetcher_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		mockDoFunc     func(req *http.Request) (*http.Response, error)
		expectedAmount string
		expectedError  error
	}{
		{
			name: "successful response",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body: io.NopCloser(bytes.NewReader([]byte(`{
						"data": {
							"amount": "8",
							"base": "USD",
							"currency": "UAH"
						}
					}`))),
				}
				return resp, nil
			},
			expectedAmount: "8",
			expectedError:  nil,
		},
		{
			name: "error making request",
			mockDoFunc: func(_ *http.Request) (*http.Response, error) {
				return nil, errors.New("error making request")
			},
			expectedAmount: "",
			expectedError:  errors.New("error making request: error making request"),
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
			expectedAmount: "",
			expectedError:  errors.New("error reading the response body: read error"),
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
			expectedAmount: "",
			expectedError:  errors.New("API request failed with status 500: internal server error"),
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
			expectedAmount: "",
			expectedError: errors.New("error unmarshaling the response: invalid character 'i'" +
				"looking for beginning of value, response body: invalid json"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &MockHTTPClient{
				DoFunc: tt.mockDoFunc,
			}
			fetcher := rateapi.NewCoinbaseFetcher(client)
			amount, err := fetcher.Fetch(context.Background(), "USD", "UAH")

			if amount != tt.expectedAmount {
				t.Errorf("expected amount %s, got %s", tt.expectedAmount, amount)
			}

			if err == nil && tt.expectedError != nil {
				t.Errorf("expected error %v, got none", tt.expectedError)
			}
		})
	}
}
