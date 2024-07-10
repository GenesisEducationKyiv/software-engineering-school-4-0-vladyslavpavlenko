package consumer

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"

	"github.com/segmentio/kafka-go"
)

type Sender interface {
	Send(params email.Params) error
}

type dbConnection interface {
	Migrate(models ...any) error
	AddConsumedEvent(event ConsumedEvent) error
}

type KafkaConsumer struct {
	db     dbConnection
	Reader *kafka.Reader
	Sender Sender
}

// NewKafkaConsumer initializes a new KafkaConsumer.
func NewKafkaConsumer(kafkaURL, topic string, partition int, groupID string, sender Sender, db dbConnection) (*KafkaConsumer, error) {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{kafkaURL},
		Topic:          topic,
		Partition:      partition,
		GroupID:        groupID,
		CommitInterval: 0, // disable auto-commit
	})

	err := db.Migrate(&ConsumedEvent{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to migrate offset")
	}

	return &KafkaConsumer{Reader: reader, Sender: sender, db: db}, nil
}

// Consume is a worker that consumes messages from Kafka and processes them
// to send an email using the Sender interface.
func (c *KafkaConsumer) Consume(ctx context.Context) {
	for {
		// Attempt to fetch a message from Kafka
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			log.Printf("Failed to read message: %v", err)
			continue
		}

		// Attempt to deserialize the fetched message
		data, err := outbox.DeserializeData(m.Value)
		if err != nil {
			log.Printf("Failed to deserialize data from message at offset %d: %v", m.Offset, err)
			continue
		}

		// Send a message
		c.sendMessage(data)

		// Create a record of the consumed event
		keyString := string(m.Key)
		eventID, err := strconv.ParseUint(keyString, 10, 64)
		if err != nil {
			log.Printf("Failed to parse event ID from key: %v", err)
			continue
		}

		consumedEvent := ConsumedEvent{
			ID:         uint(eventID),
			Data:       data.Serialize(),
			ConsumedAt: time.Now(),
		}

		// Attempt to add the consumed event to the database
		if err = c.db.AddConsumedEvent(consumedEvent); err != nil {
			log.Printf("Failed to record consumed event at offset %d: %v", m.Offset, err)
			continue
		}

		// Commit the offset back to Kafka to mark the message as processed
		if err = c.Reader.CommitMessages(ctx, m); err != nil {
			log.Printf("Failed to commit message offset %d: %v", m.Offset, err)
		} else {
			log.Printf("Offset committed successfully for message at offset: %d", m.Offset)
		}
	}
}

func (c *KafkaConsumer) sendMessage(data outbox.Data) {
	params := email.Params{
		To:      data.Email,
		Subject: "USD to UAH Exchange Rate",
		Body:    fmt.Sprintf("The current exchange rate for USD to UAH is %.2f.", data.Rate),
	}

	log.Printf("Sending email to %s", data.Email)

	err := c.Sender.Send(params)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	}
}
