package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/config"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"time"
)

var counts int64

func setup(app *config.AppConfig) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	db, err := connectToDB(dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = runDBMigrations(db)
	if err != nil {
		log.Fatal(err)
	}

	app.DB = db
	app.Models = models.New(db)

	app.EmailConfig = email.Config{
		Email:    os.Getenv("GMAIL_EMAIL"),
		Password: os.Getenv("GMAIL_PASSWORD"),
	}

	if app.EmailConfig.Email == "" || app.EmailConfig.Password == "" {
		log.Fatal("Missing email configuration in environment variables")
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
	err := db.AutoMigrate(&models.Currency{}, &models.User{}, &models.Subscription{})
	if err != nil {
		return fmt.Errorf("error during migration: %v", err)
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

	if err := db.Model(&models.Currency{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialCurrencies := []models.Currency{
		{Code: "USD", Name: "United States Dollar"},
		{Code: "UAH", Name: "Ukrainian Hryvnia"},
	}

	if err := db.Create(&initialCurrencies).Error; err != nil {
		return err
	}

	return nil
}
