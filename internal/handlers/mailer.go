package handlers

import (
	"fmt"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/models"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rate"
	"log"
	"sync"
)

// NotifySubscribers handles sending emails to all the subscribers.
func (m *Repository) NotifySubscribers() {
	var subscriptions []models.Subscription
	m.App.DB.Preload("User").Preload("BaseCurrency").Preload("TargetCurrency").Find(&subscriptions)

	var wg sync.WaitGroup
	for _, subscription := range subscriptions {
		wg.Add(1)

		baseCode := subscription.BaseCurrency.Code
		targetCode := subscription.TargetCurrency.Code

		price, err := rate.GetRate(baseCode, targetCode)
		if err != nil {
			log.Printf("Failed to retrieve rate for %s to %s: %v", baseCode, targetCode, err)
			wg.Done()
			continue
		}

		params := email.Params{
			To:      subscription.User.Email,
			Subject: fmt.Sprintf("%s to %s Exchange Rate", baseCode, targetCode),
			Body:    fmt.Sprintf("The current exchange rate for %s to %s is %s.", baseCode, targetCode, price),
		}

		go email.SendEmail(&wg, m.App.EmailConfig, params)
	}
	wg.Wait()
}
