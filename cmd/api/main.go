package main

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	a := config.New()
	l := logger.New()

	err := app.Run(a, l)
	if err != nil {
		l.Fatal("failed to start application", zap.Error(err))
	}
}
