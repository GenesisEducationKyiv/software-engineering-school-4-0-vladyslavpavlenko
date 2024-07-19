package consumer

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
	"go.uber.org/zap"

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
	logger *logger.Logger
}

// NewKafkaConsumer initializes a new KafkaConsumer.
func NewKafkaConsumer(kafkaURL, topic string, partition int, groupID string,
	sender Sender, db dbConnection, l *logger.Logger,
) (*KafkaConsumer, error) {
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

	if l == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &KafkaConsumer{Reader: reader, Sender: sender, db: db, logger: l}, nil
}

// Consume is a worker that consumes messages from Kafka and processes them
// to send an email using the Sender interface.
func (c *KafkaConsumer) Consume(ctx context.Context) {
	for {
		// Check if context is canceled before attempting to fetch a message
		select {
		case <-ctx.Done():
			c.logger.Warn("shutting down consumer...", zap.String("cause", "context canceled"))
			return
		default:
		}

		// Attempt to fetch a message from Kafka
		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.logger.Warn("shutting down consumer...", zap.String("cause", "context canceled"))
				return
			}
			c.logger.Error("failed to read message", zap.Error(err))
			continue
		}

		// Attempt to deserialize the fetched message
		data, err := outbox.DeserializeData(m.Value)
		if err != nil {
			c.logger.Error("failed to deserialize data", zap.Int64("offset", m.Offset), zap.Error(err))
			continue
		}

		// Send a message
		c.sendMessage(data)

		// Create a record of the consumed event
		keyString := string(m.Key)
		eventID, err := strconv.ParseUint(keyString, 10, 64)
		if err != nil {
			c.logger.Error("failed to parse event id", zap.Error(err))
			continue
		}

		sData, err := data.Serialize()
		if err != nil {
			c.logger.Error("failed to serialize data", zap.Error(err))
		}

		consumedEvent := ConsumedEvent{
			ID:         uint(eventID),
			Data:       sData,
			ConsumedAt: time.Now(),
		}

		// Attempt to add the consumed event to the database
		if err = c.db.AddConsumedEvent(consumedEvent); err != nil {
			c.logger.Error("failed to record consumed event", zap.Int64("offset", m.Offset), zap.Error(err))
			continue
		}

		// Commit the offset back to Kafka to mark the message as processed
		if err = c.Reader.CommitMessages(ctx, m); err != nil {
			c.logger.Error("failed to commit message", zap.Int64("offset", m.Offset), zap.Error(err))
		} else {
			c.logger.Error("offset committed", zap.Int64("offset", m.Offset))
		}
	}
}

func (c *KafkaConsumer) sendMessage(data outbox.Data) {
	params := email.Params{
		To:      data.Email,
		Subject: "USD to UAH Exchange Rate",
		Body:    fmt.Sprintf("The current exchange rate for USD to UAH is %.2f.", data.Rate),
	}

	c.logger.Debug("sending email", zap.String("email", data.Email))

	err := c.Sender.Send(params)
	if err != nil {
		c.logger.Error("failed to send email", zap.Error(err))
	}
}
