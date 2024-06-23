package main

import (
	"log"
	"os"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/app"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/app/config"
)

func main() {
	appConfig := config.NewAppConfig()

	err := app.Run(appConfig)
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
