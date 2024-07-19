package main

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
)

func main() {
	a := config.New()
	a.Logger = logger.New()

	err := app.Run(a)
	if err != nil {
		a.Logger.Fatal(err.Error())
	}
}
