package consumer

import (
	"context"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metrics"

	"github.com/vladyslavpavlenko/genesis-api-project/pkg/logger"
	"go.uber.org/zap"

	"github.com/pkg/errors"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/outbox"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"

	"github.com/segmentio/kafka-go"
)

var (
	sentEmailsCounter    = metrics.NewCounter("sent_emails_count")
	notSentEmailsCounter = metrics.NewCounter("not_sent_emails_count")
	emailSendingDuration = metrics.NewHistogram("email_sending_duration_seconds")
	_                    = metrics.NewGauge("email_sending_success_rate", calculateEmailSuccessRate)
)

type (
	sender interface {
		Send(params email.Params) error
	}

	dbConnection interface {
		Migrate(models ...any) error
		AddConsumedEvent(event ConsumedEvent) error
	}
)

type KafkaConsumer struct {
	db     dbConnection
	Reader *kafka.Reader
	Sender sender
	l      *logger.Logger
}

// NewKafkaConsumer initializes a new KafkaConsumer.
func NewKafkaConsumer(kafkaURL, topic string, partition int, groupID string,
	sender sender, db dbConnection, l *logger.Logger,
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

	return &KafkaConsumer{Reader: reader, Sender: sender, db: db, l: l}, nil
}

// Consume is a worker that consumes messages from Kafka and processes them
// to send an email using the Sender interface.
func (c *KafkaConsumer) Consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			c.l.Info("shutting down consumer...", zap.String("cause", "context canceled"))
			return
		default:
		}

		m, err := c.Reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				c.l.Info("shutting down consumer...", zap.String("cause", "context canceled"))
				return
			}
			c.l.Error("failed to read message", zap.Error(err))
			continue
		}

		start := time.Now()

		data, err := outbox.DeserializeData(m.Value)
		if err != nil {
			c.l.Error("failed to deserialize data", zap.Int64("offset", m.Offset), zap.Error(err))
			continue
		}

		err = c.sendMessage(data)
		if err != nil {
			notSentEmailsCounter.Inc()
			c.l.Error("failed to send email", zap.Int64("offset", m.Offset), zap.Error(err))
		} else {
			sentEmailsCounter.Inc()
		}

		// Record the duration it took to process the email
		emailSendingDuration.UpdateDuration(start)

		if err = c.Reader.CommitMessages(ctx, m); err != nil {
			c.l.Error("failed to commit message", zap.Int64("offset", m.Offset), zap.Error(err))
			continue
		}

		c.l.Debug("offset committed", zap.Int64("offset", m.Offset))
	}
}

func (c *KafkaConsumer) sendMessage(data outbox.Data) error {
	params := email.Params{
		To:      data.Email,
		Subject: "USD to UAH Exchange Rate",
		Body:    fmt.Sprintf("The current exchange rate for USD to UAH is %.2f.", data.Rate),
	}

	err := c.Sender.Send(params)
	if err != nil {
		return err
	}

	return nil
}
