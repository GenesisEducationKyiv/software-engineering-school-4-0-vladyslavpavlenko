package gormsubscriber

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormstorage"
)

var ErrorInvalidEmail = errors.New("invalid email")

// validateSubscription is an action that validates a subscription by validating an
// email address and checking if it already exists.
func validateSubscription(saga *State, s *Subscriber) error {
	// Validate the email format
	if !email.Email(saga.Email).Validate() {
		return ErrorInvalidEmail
	}

	// Check if the subscription already exists
	ctx, cancel := context.WithTimeout(context.Background(), gormstorage.RequestTimeout)
	defer cancel()

	var count int64
	err := s.db.WithContext(ctx).Model(&models.Subscription{}).Where("email = ?", saga.Email).Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return ErrorDuplicateSubscription
	}

	return nil
}

// addSubscription is an action that creates a new models.Subscription record.
func addSubscription(saga *State, s *Subscriber) error {
	ctx, cancel := context.WithTimeout(context.Background(), gormstorage.RequestTimeout)
	defer cancel()

	subscription := models.Subscription{
		Email:     saga.Email,
		CreatedAt: time.Now(),
	}
	result := s.db.WithContext(ctx).Create(&subscription)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return ErrorDuplicateSubscription
		}
		return result.Error
	}
	return nil
}
