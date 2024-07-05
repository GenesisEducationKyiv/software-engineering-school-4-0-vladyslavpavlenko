package consumer

import (
	"context"
	"fmt"
	"log"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"

	"github.com/segmentio/kafka-go"
)

type Sender interface {
	Send(params email.Params) error
}

type KafkaConsumer struct {
	Reader *kafka.Reader
	Sender Sender
}

// NewKafkaConsumer initializes a new KafkaConsumer.
func NewKafkaConsumer(kafkaURL, topic string, partition int) *KafkaConsumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaURL},
		Topic:          topic,
		Partition:      partition,
		CommitInterval: 0, // disable auto-commit
	})

	return &KafkaConsumer{Reader: reader}
}

// ConsumeMessages reads messages from Kafka, deserializes them into Event structs, and processes them.
func (c *KafkaConsumer) ConsumeMessages(ctx context.Context) {
	for {
		// Read a message from Kafka
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		// Deserialize the data from the message
		data, err := outbox.DeserializeData(m.Value)
		if err != nil {
			log.Printf("Failed to deserialize data from message: %v", err)
			continue
		}

		// Process the message
		sendMessage(data, c.Sender)

		// Commit the offset after processing the message
		if err = c.Reader.CommitMessages(ctx, m); err != nil {
			log.Printf("Failed to commit message offset: %v", err)
		}
	}
}

func sendMessage(data outbox.Data, sender Sender) {
	params := email.Params{
		To:      data.Email,
		Subject: "USD to UAH Exchange Rate",
		Body:    fmt.Sprintf("The current exchange rate for USD to UAH is %.2f.", data.Rate),
	}

	log.Printf("Sending email to %s", data.Email)

	err := sender.Send(params)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}
