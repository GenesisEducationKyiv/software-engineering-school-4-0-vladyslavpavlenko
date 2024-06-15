package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
)

// NotifySubscribers handles sending currency update rateapi emails to all the subscribers.
func (m *Repository) NotifySubscribers() error {
	subscriptions, err := m.App.Models.Subscription.GetSubscriptions()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, subscription := range subscriptions {
		wg.Add(1)

		baseCode := subscription.BaseCurrency.Code
		targetCode := subscription.TargetCurrency.Code

		fetcher := rateapi.CoinbaseFetcher{
			Client: &http.Client{},
		}

		price, err := fetcher.FetchRate(baseCode, targetCode)
		if err != nil {
			log.Printf("Failed to retrieve rateapi for %s to %s: %v", baseCode, targetCode, err)
			wg.Done()
			continue
		}

		params := email.Params{
			To:      subscription.User.Email,
			Subject: fmt.Sprintf("%s to %s Exchange Rate", baseCode, targetCode),
			Body:    fmt.Sprintf("The current exchange rateapi for %s to %s is %s.", baseCode, targetCode, price),
		}

		go email.SendEmail(&wg, m.App.EmailConfig, params)
	}
	wg.Wait()

	return nil
}
