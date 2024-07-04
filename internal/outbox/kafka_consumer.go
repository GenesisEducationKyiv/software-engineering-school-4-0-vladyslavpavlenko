package outbox

import (
	"context"
	"fmt"
	"log"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"

	"github.com/segmentio/kafka-go"
)

// NewKafkaReader initializes a new kafka.Reader with a specific topic and group.
func NewKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaURL},
		Topic:     topic,
		GroupID:   groupID,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
}

// ConsumeMessages reads messages from Kafka, deserializes them into Event structs, and processes them.
func ConsumeMessages(ctx context.Context, reader *kafka.Reader) {
	for {
		m, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		log.Printf("Message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))

		data, err := DeserializeData(m.Value)
		if err != nil {
			log.Printf("Failed to deserialize data from message: %v", err)
		}

		sendMessage(data)
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
