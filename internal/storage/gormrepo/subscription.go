package gormrepo

import (
	"errors"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
)

var (
	ErrorDuplicateSubscription   = errors.New("subscription already exists")
	ErrorNonExistentSubscription = errors.New("subscription does not exist")
)

// AddSubscription creates a new models.Subscription record.
func (c *Connection) AddSubscription(email string) error {
	subscription := models.Subscription{
		Email:     email,
		CreatedAt: time.Now(),
	}
	result := c.DB.Create(&subscription)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return ErrorDuplicateSubscription
		}
		return result.Error
	}
	return nil
}

// DeleteSubscription deletes a models.Subscription record.
func (c *Connection) DeleteSubscription(email string) error {
	result := c.DB.Where("email = ?", email).Delete(&models.Subscription{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrorNonExistentSubscription
	}

	return nil
}

// GetSubscriptions returns a paginated list of subscriptions. Limit specify the number of records to be retrieved
// Limit conditions can be canceled by using `Limit(-1)`. Offset specify the number of records to skip before starting
// to return the records. Offset conditions can be canceled by using `Offset(-1)`.
func (c *Connection) GetSubscriptions(limit, offset int) ([]models.Subscription, error) {
	var subscriptions []models.Subscription
	result := c.DB.Limit(limit).Offset(offset).Find(&subscriptions)
	if result.Error != nil {
		return nil, result.Error
	}

	return subscriptions, nil
}
