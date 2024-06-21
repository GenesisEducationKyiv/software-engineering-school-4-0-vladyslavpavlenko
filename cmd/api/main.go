package main

import (
	"log"
	"os"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
)

const webPort = 8080

var app config.AppConfig

func main() {
	err := run()
	if err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
