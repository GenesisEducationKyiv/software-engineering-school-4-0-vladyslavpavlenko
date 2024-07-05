package outbox

import (
	"context"
	"fmt"
	"log"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"

	"github.com/segmentio/kafka-go"
)

// NewKafkaReader initializes a new kafka.Reader with a specific topic and group.
func NewKafkaReader(kafkaURL, topic string, partition int) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaURL},
		Topic:          topic,
		Partition:      partition,
		CommitInterval: 0, // disable auto-commit
	})
}

// ConsumeMessages reads messages from Kafka, deserializes them into Event structs, and processes them.
func ConsumeMessages(ctx context.Context, reader *kafka.Reader) {
	for {
		// Read a message from Kafka
		m, err := reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		// Deserialize the data from the message
		data, err := DeserializeData(m.Value)
		if err != nil {
			log.Printf("Failed to deserialize data from message: %v", err)
			continue
		}

		// Process the message
		sendMessage(data)

		// Commit the offset after processing the message
		if err = reader.CommitMessages(ctx, m); err != nil {
			log.Printf("Failed to commit message offset: %v", err)
		}
	}
}

func sendMessage(data Data) {
	params := email.Params{
		To:      data.Email,
		Subject: "USD to UAH Exchange Rate",
		Body:    fmt.Sprintf("The current exchange rate for USD to UAH is %.2f.", data.Rate),
	}

	log.Printf("Sending email to %s", data.Email)

	err := email.SenderService.Send(params)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}
