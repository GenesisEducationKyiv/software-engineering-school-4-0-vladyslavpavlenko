package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/kelseyhightower/envconfig"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/gormrepo"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

// envVariables holds environment variables used in the application.
type envVariables struct {
	DBURL     string `envconfig:"DB_URL"`
	DBPort    string `envconfig:"DB_PORT"`
	DBUser    string `envconfig:"DB_USER"`
	DBPass    string `envconfig:"DB_PASS"`
	DBName    string `envconfig:"DB_NAME"`
	EmailAddr string `envconfig:"EMAIL_ADDR"`
	EmailPass string `envconfig:"EMAIL_PASS"`
}

func setup(app *config.AppConfig) (dbrepo.DB, error) {
	envs, err := readEnv()
	if err != nil {
		return nil, fmt.Errorf("error reading the .env file: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		envs.DBURL,
		envs.DBPort,
		envs.DBUser,
		envs.DBPass,
		envs.DBName)

	db, err := connectDB(dsn)
	if err != nil {
		return nil, fmt.Errorf("error conntecting to the database: %w", err)
	}

	err = migrateDB(db)
	if err != nil {
		return nil, fmt.Errorf("error runnning database migrations: %w", err)
	}

	app.Models = gormrepo.NewModels(db)

	app.EmailConfig, err = email.NewEmailConfig(envs.EmailAddr, envs.EmailPass)
	if err != nil {
		return nil, errors.New("error setting up email configuration")
	}

	repo := handlers.NewRepo(app, db)
	handlers.NewHandlers(repo)

	return db, nil
}

// readEnv reads and returns the environmental variables as an envVariables object.
func readEnv() (envVariables, error) {
	var envs envVariables
	err := envconfig.Process("", &envs)
	if err != nil {
		return envVariables{}, err
	}
	return envs, nil
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
