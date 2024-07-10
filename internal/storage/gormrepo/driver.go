package gormrepo

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const timeout = time.Second * 5

type Connection struct {
	db *gorm.DB
}

// DB returns a pointer to gorm.DB.
func (c *Connection) DB() *gorm.DB {
	return c.db
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
			c.db = db
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
	sqlDB, err := c.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Migrate performs a database migration for given models.
func (c *Connection) Migrate(models ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := c.db.WithContext(ctx).AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("error migrating models: %w", err)
	}

	return nil
}

// BeginTransaction begins a transaction.
func (c *Connection) BeginTransaction() (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
