package gormsubscriber

import (
	"context"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/storage/gormstorage"
)

// deleteSubscription is a compensation to addSubscription that deletes a
// models.Subscription record, queried by an email address.
func deleteSubscription(saga *State, s *Subscriber) error {
	ctx, cancel := context.WithTimeout(context.Background(), gormstorage.RequestTimeout)
	defer cancel()

	result := s.db.WithContext(ctx).Where("email = ?", saga.Email).Delete(&models.Subscription{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
