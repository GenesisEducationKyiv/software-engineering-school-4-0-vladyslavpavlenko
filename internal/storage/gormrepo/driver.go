package gormrepo

import (
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Connection struct {
	DB *gorm.DB
}

// Setup sets up a new Connection.
func (c *Connection) Setup(dsn string) error {
	var counts int64
	for {
		db, err := openDB(dsn)
		if err != nil {
			log.Printf("Postgres not yet ready... Attempt: %d\n", counts)
			counts++
		} else {
			log.Println("Connected to Postgres!")
			c.DB = db
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

// openDB initializes a new gorm.DB database connection.
func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Close closes a database connection.
func (c *Connection) Close() error {
	sqlDB, err := c.DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate performs a database migration for given models.
func (c *Connection) Migrate(models ...any) error {
	err := c.DB.AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("error migrating models: %w", err)
	}

	return nil
}
