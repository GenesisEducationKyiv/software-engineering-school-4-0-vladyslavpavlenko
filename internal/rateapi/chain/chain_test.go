package chain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi/chain"
)

type MockFetcher struct {
	FetchFunc func(ctx context.Context, base, target string) (string, error)
}

func (m *MockFetcher) Fetch(ctx context.Context, base, target string) (string, error) {
	return m.FetchFunc(ctx, base, target)
}

func TestNode_Fetch(t *testing.T) {
	tests := []struct {
		name          string
		fetcher       *MockFetcher
		nextFetcher   *MockFetcher
		base          string
		target        string
		expectedRate  string
		expectedError string
	}{
		{
			name: "successful fetch",
			fetcher: &MockFetcher{
				FetchFunc: func(_ context.Context, _, _ string) (string, error) {
					return "10.0", nil
				},
			},
			nextFetcher:   nil,
			base:          "USD",
			target:        "UAH",
			expectedRate:  "10.0",
			expectedError: "",
		},
		{
			name: "fetch error, no next fetcher",
			fetcher: &MockFetcher{
				FetchFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", errors.New("fetch error")
				},
			},
			nextFetcher:   nil,
			base:          "USD",
			target:        "UAH",
			expectedRate:  "",
			expectedError: "fetch error",
		},
		{
			name: "fetch error, delegate to next fetcher",
			fetcher: &MockFetcher{
				FetchFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", errors.New("fetch error")
				},
			},
			nextFetcher: &MockFetcher{
				FetchFunc: func(_ context.Context, _, _ string) (string, error) {
					return "20.0", nil
				},
			},
			base:          "USD",
			target:        "UAH",
			expectedRate:  "20.0",
			expectedError: "",
		},
		{
			name: "fetch error, delegate to next fetcher with error",
			fetcher: &MockFetcher{
				FetchFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", errors.New("fetch error")
				},
			},
			nextFetcher: &MockFetcher{
				FetchFunc: func(_ context.Context, _, _ string) (string, error) {
					return "", errors.New("next fetcher error")
				},
			},
			base:          "USD",
			target:        "UAH",
			expectedRate:  "",
			expectedError: "next fetcher error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := chain.NewNode(tt.fetcher)
			if tt.nextFetcher != nil {
				nextNode := chain.NewNode(tt.nextFetcher)
				node.SetNext(nextNode)
			}

			rate, err := node.Fetch(context.Background(), tt.base, tt.target)
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
