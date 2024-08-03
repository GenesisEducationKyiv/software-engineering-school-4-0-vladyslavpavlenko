package rateapi

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"

	"go.uber.org/zap"
)

type fetcher interface {
	Fetch(ctx context.Context, base, target string) (string, error)
}

type FetcherWithLogger struct {
	name    string
	fetcher fetcher
	l       *logger.Logger
}

// NewFetcherWithLogger creates and returns a pointer to a new FetcherWithLogger.
func NewFetcherWithLogger(name string, f fetcher, l *logger.Logger) *FetcherWithLogger {
	return &FetcherWithLogger{
		name:    name,
		fetcher: f,
		l:       l,
	}
}

// Fetch performs a call to the Fetcher.
func (f *FetcherWithLogger) Fetch(ctx context.Context, base, target string) (string, error) {
	rate, err := f.fetcher.Fetch(ctx, base, target)
	if err != nil {
		f.l.Error("error fetching rate", zap.String("fetcher", f.name), zap.Error(err))
		return "", err
	}

	f.l.Info("rate fetched", zap.String("fetcher", f.name), zap.String("rate", rate))
	return rate, nil
}
