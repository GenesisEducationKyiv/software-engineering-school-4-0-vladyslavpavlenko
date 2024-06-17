package gormrepo

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDB struct {
	DB *gorm.DB
}

// Connect implements the DB interface for GormDB.
func (g *GormDB) Connect(dsn string) error {
	var counts int64
	for {
		db, err := openDB(dsn)
		if err != nil {
			log.Printf("Postgres not yet ready... Attempt: %d\n", counts)
			counts++
		} else {
			log.Println("Connected to Postgres!")
			g.DB = db
			return nil
		}

		if counts > 10 {
			log.Println("Maximum retry attempts exceeded:", err)
			return err
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}

// openDB initializes a new database connection.
func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (g *GormDB) Close() error {
	sqlDB, err := g.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

func (g *GormDB) Migrate() error {
	err := g.DB.AutoMigrate(&models.Currency{}, &models.User{}, &models.Subscription{})
	if err != nil {
		return fmt.Errorf("error during migration: %w", err)
	}

	err = createInitialCurrencies(g)
	if err != nil {
		return errors.New(fmt.Sprint("error creating initial currencies:", err))
	}

	return nil
}

// createInitialCurrencies creates initial currencies in the `currencies` table.
func createInitialCurrencies(conn *GormDB) error {
	var count int64

	if err := conn.DB.Model(&models.Currency{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	initialCurrencies := []models.Currency{
		{Code: "USD", Name: "United States Dollar"},
		{Code: "UAH", Name: "Ukrainian Hryvnia"},
	}

	if err := conn.DB.Create(&initialCurrencies).Error; err != nil {
		return err
	}

	return nil
}
