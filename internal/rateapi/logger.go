package rateapi

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"

	"go.uber.org/zap"
)

type Fetcher interface {
	Fetch(ctx context.Context, base, target string) (string, error)
}

type FetcherWithLogger struct {
	name    string
	fetcher Fetcher
	logger  *logger.Logger
}

// NewFetcherWithLogger creates and returns a pointer to a new FetcherWithLogger.
func NewFetcherWithLogger(name string, f Fetcher, l *logger.Logger) *FetcherWithLogger {
	return &FetcherWithLogger{
		name:    name,
		fetcher: f,
		logger:  l,
	}
}

// Fetch performs a call to the Fetcher.
func (f *FetcherWithLogger) Fetch(ctx context.Context, base, target string) (string, error) {
	rate, err := f.fetcher.Fetch(ctx, base, target)
	if err != nil {
		f.logger.Error("fetch error", zap.String("fetcher", f.name), zap.Error(err))
		return "", err
	}

	f.logger.Info("fetch successful", zap.String("fetcher", f.name), zap.String("rate", rate))
	return rate, nil
}
