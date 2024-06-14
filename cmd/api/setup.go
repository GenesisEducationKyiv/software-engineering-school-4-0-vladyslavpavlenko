package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/gormrepo"

	"github.com/joho/godotenv"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var counts int64

func setup(app *config.AppConfig) error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading the .env file: %w", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := connectToDB(dsn)
	if err != nil {
		return fmt.Errorf("error conntecting to the database: %w", err)
	}

	err = runDBMigrations(db)
	if err != nil {
		return fmt.Errorf("error runnning database migrations: %w", err)
	}

	app.DB = db
	app.Models = gormrepo.New(db)

	app.EmailConfig = email.Config{
		Email:    os.Getenv("GMAIL_EMAIL"),
		Password: os.Getenv("GMAIL_PASSWORD"),
	}

	if app.EmailConfig.Email == "" || app.EmailConfig.Password == "" {
		return errors.New("missing email configuration in environment variables")
	}

	repo := handlers.NewRepo(app)
	handlers.NewHandlers(repo)

	return nil
}

// openDB initializes a new database connection.
func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

// connectToDB sets up a GORM database connection.
func connectToDB(dsn string) (*gorm.DB, error) {
	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection, nil
		}

		if counts > 10 {
			log.Println(err)
			return nil, err
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
		continue
	}
}

// runDBMigrations runs database migrations.
func runDBMigrations(db *gorm.DB) error {
	log.Println("Running migrations...")
	// create tables
	err := db.AutoMigrate(&dbrepo.Currency{}, &dbrepo.User{}, &dbrepo.Subscription{})
	if err != nil {
		return fmt.Errorf("error during migration: %w", err)
	}

	// populate tables with initial data
	err = createInitialCurrencies(db)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial currencies:", err))
	}

	log.Println("Database migrated!")

	return nil
}

// createInitialCurrencies creates initial currencies in the `currencies` table.
func createInitialCurrencies(db *gorm.DB) error {
	var count int64

	if err := db.Model(&dbrepo.Currency{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialCurrencies := []dbrepo.Currency{
		{Code: "USD", Name: "United States Dollar"},
		{Code: "UAH", Name: "Ukrainian Hryvnia"},
	}

	if err := db.Create(&initialCurrencies).Error; err != nil {
		return err
	}

	return nil
}
