package handlers

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
)

const batchSize = 100

// ProduceMailingEvents handles producing events for currency rate update emails.
func (m *Repository) ProduceMailingEvents() error {
	rate, err := m.Services.Fetcher.Fetch(context.Background(), "USD", "UAH")
	if err != nil {
		return fmt.Errorf("failed to retrieve rate: %w", err)
	}

	floatRate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return fmt.Errorf("failed to parse rate: %w", err)
	}

	var offset int
	errChan := make(chan error, 1)
	for {
		subscriptions, err := m.DB.GetSubscriptions(batchSize, offset)
		if err != nil {
			return err
		}
		if len(subscriptions) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, sub := range subscriptions {
			wg.Add(1)
			go func(sub models.Subscription) {
				defer wg.Done()

				data := outbox.Data{Email: sub.Email, Rate: floatRate}
				if localErr := m.App.Outbox.AddEvent(data); localErr != nil {
					select {
					case errChan <- localErr:
					default:
					}
				}
			}(sub)
		}
		wg.Wait()

		select {
		case err := <-errChan:
			return err
		default:
		}

		offset += batchSize
	}
	return nil
}
