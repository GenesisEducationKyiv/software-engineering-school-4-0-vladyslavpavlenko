package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/gormrepo"

	"github.com/joho/godotenv"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

// envVariables holds environment variables used in the application.
type envVariables struct {
	DbHost    string
	DbPort    string
	DbUser    string
	DbPass    string
	DbName    string
	EmailAddr string
	EmailPass string
}

func setup(app *config.AppConfig) error {
	envs, err := readEnv()
	if err != nil {
		return fmt.Errorf("error reading the .env file: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		envs.DbHost,
		envs.DbPort,
		envs.DbUser,
		envs.DbPass,
		envs.DbName)

	db, err := connectDB(dsn)
	if err != nil {
		return fmt.Errorf("error conntecting to the database: %w", err)
	}

	err = migrateDB(db)
	if err != nil {
		return fmt.Errorf("error runnning database migrations: %w", err)
	}

	app.DB = db
	app.Models = gormrepo.NewModels(db)

	app.EmailConfig, err = email.NewEmailConfig(envs.EmailAddr, envs.EmailPass)
	if err != nil {
		return errors.New("error setting up email configuration")
	}

	repo := handlers.NewRepo(app)
	handlers.NewHandlers(repo)

	return nil
}

// readEnv reads and returns the environmental variables as an envVariables object.
func readEnv() (envVariables, error) {
	err := godotenv.Load()
	if err != nil {
		return envVariables{}, err
	}

	return envVariables{
		DbHost:    os.Getenv("DB_HOST"),
		DbPort:    os.Getenv("DB_PORT"),
		DbUser:    os.Getenv("DB_USER"),
		DbPass:    os.Getenv("DB_PASS"),
		DbName:    os.Getenv("DB_NAME"),
		EmailAddr: os.Getenv("EMAIL_ADDR"),
		EmailPass: os.Getenv("EMAIL_PASS"),
	}, nil
}

// connectDB sets up a GORM database connection and returns an interface.
func connectDB(dsn string) (*gormrepo.GormDB, error) {
	var db gormrepo.GormDB

	err := db.Connect(dsn)
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// migrateDB runs database migrations.
func migrateDB(db *gormrepo.GormDB) error {
	log.Println("Running migrations...")

	err := db.Migrate()
	if err != nil {
		return fmt.Errorf("error running migrations: %w", err)
	}

	log.Println("Database migrated!")

	return nil
}
