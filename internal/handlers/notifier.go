package handlers

import (
	"context"
	"fmt"
	"log"
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

	floatRate, err := strconv.ParseFloat(rate, 32)
	if err != nil {
		return fmt.Errorf("failed to parse price: %w", err)
	}

	var offset int
	for {
		subscriptions, err := m.DB.GetSubscriptions(batchSize, offset)
		if err != nil {
			return err
		}
		if len(subscriptions) == 0 {
			break
		}

		var wg sync.WaitGroup
		for _, subscription := range subscriptions {
			wg.Add(1)
			go func(sub models.Subscription) {
				defer wg.Done()

				data := outbox.Data{Email: sub.Email, Rate: floatRate}

				if err = m.App.Outbox.AddEvent(data); err != nil {
					log.Printf("error adding event: %v", err)
				}
			}(subscription)
		}
		wg.Wait()

		offset += batchSize
	}

	return nil
}
