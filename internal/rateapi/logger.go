package rateapi

import (
	"context"
	"log"
)

type (
	FetcherWithLogger struct {
		name    string
		fetcher Fetcher
	}

	Fetcher interface {
		Fetch(ctx context.Context, base, target string) (string, error)
	}
)

// NewFetcherWithLogger creates and returns a pointer to a new FetcherWithLogger.
func NewFetcherWithLogger(name string, fetcher Fetcher) *FetcherWithLogger {
	return &FetcherWithLogger{
		name:    name,
		fetcher: fetcher,
	}
}

// Fetch performs a call to the Fetcher.
func (f *FetcherWithLogger) Fetch(ctx context.Context, base, target string) (string, error) {
	rate, err := f.fetcher.Fetch(ctx, base, target)
	if err != nil {
		log.Printf("[%s]: error: %v", f.name, err)
		return "", err
	}

	log.Printf("[%s]: rate: %+v", f.name, rate)
	return rate, nil
}
