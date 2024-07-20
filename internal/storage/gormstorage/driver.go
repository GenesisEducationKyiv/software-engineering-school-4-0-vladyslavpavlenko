package gormstorage

import (
	"context"
	"fmt"
	"time"

	glogger "gorm.io/gorm/logger"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const RequestTimeout = time.Second * 5

type Connection struct {
	db *gorm.DB
	l  *logger.Logger
}

// DB returns a pointer to gorm.DB.
func (c *Connection) DB() *gorm.DB {
	return c.db
}

// Setup sets up a new Connection with a logger.
func (c *Connection) Setup(dsn string, l *logger.Logger) error {
	c.l = l
	var counts int64
	for {
		db, err := openDB(dsn)
		if err != nil {
			c.l.Error("Postgres not yet ready...", zap.Int64("attempt", counts), zap.Error(err))
			counts++
		} else {
			c.l.Debug("connected to Postgres!")
			c.db = db
			return nil
		}

		if counts > 10 {
			c.l.Error("maximum retry attempts exceeded", zap.Error(err))
			return err
		}

		c.l.Debug("backing off for two seconds...")
		time.Sleep(2 * time.Second)
	}
}

// openDB initializes a new gorm.DB database connection.
func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: glogger.Default.LogMode(glogger.Silent),
	})
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
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	err := c.db.WithContext(ctx).AutoMigrate(models...)
	if err != nil {
		return fmt.Errorf("error migrating models: %w", err)
	}

	return nil
}

// BeginTransaction begins a transaction.
func (c *Connection) BeginTransaction() (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
	defer cancel()

	tx := c.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}
	return tx, nil
}
