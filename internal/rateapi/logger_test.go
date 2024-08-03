package rateapi_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
)

type MockFetcher struct {
	FetchFunc func(ctx context.Context, base, target string) (string, error)
}

func (m *MockFetcher) Fetch(ctx context.Context, base, target string) (string, error) {
	return m.FetchFunc(ctx, base, target)
}

func TestFetcherWithLogger_Fetch(t *testing.T) {
	tests := []struct {
		name           string
		mockFetchFunc  func(ctx context.Context, base, target string) (string, error)
		base           string
		target         string
		expectedRate   string
		expectedError  string
		expectedLogMsg string
	}{
		{
			name: "successful fetch",
			mockFetchFunc: func(_ context.Context, _, _ string) (string, error) {
				return "10.0", nil
			},
			base:           "USD",
			target:         "UAH",
			expectedRate:   "10.0",
			expectedError:  "",
			expectedLogMsg: "[TestFetcher]: rate: 10.0\n",
		},
		{
			name: "fetch error",
			mockFetchFunc: func(_ context.Context, _, _ string) (string, error) {
				return "", errors.New("fetch error")
			},
			base:           "USD",
			target:         "UAH",
			expectedRate:   "",
			expectedError:  "fetch error",
			expectedLogMsg: "[TestFetcher]: error: fetch error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockFetcher := &MockFetcher{
				FetchFunc: tt.mockFetchFunc,
			}

			fetcherWithLogger := rateapi.NewFetcherWithLogger("TestFetcher", mockFetcher, nil)
			rate, err := fetcherWithLogger.Fetch(context.Background(), tt.base, tt.target)

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
