package outbox

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

type Outbox interface {
	AddEvent(data Data) error
	GetUnpublishedEvents() ([]Event, error)
	MarkEventAsPublished(eventID uint) error
	Cleanup()
}

type KafkaProducer struct {
	Writer *kafka.Writer
	Outbox Outbox
}

// NewKafkaProducer initializes a new KafkaProducer.
func NewKafkaProducer(kafkaURL string, outbox Outbox) *KafkaProducer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	return &KafkaProducer{Writer: writer, Outbox: outbox}
}

// NewTopic creates a new kafka.TopicConfig if it does not exist.
func (p *KafkaProducer) NewTopic(topic string, partitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", p.Writer.Addr.String())
	if err != nil {
		return err
	}
	defer conn.Close()

	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     partitions,
		ReplicationFactor: replicationFactor,
	}

	err = conn.CreateTopics(topicConfig)
	if err != nil {
		return err
	}
	return nil
}

// SetTopic changes the topic of the current kafka.Writer
func (p *KafkaProducer) SetTopic(topic string) {
	p.Writer.Topic = topic
}

// Produce fetches for unpublished events, publishes them, and marks them as published.
func (p *KafkaProducer) Produce(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down worker...")
			return
		case <-ticker.C:
			p.processEvents(p.Outbox)
		}
	}
}

// processEvents handles the retrieval and publishing of events.
func (p *KafkaProducer) processEvents(outbox Outbox) {
	events, err := outbox.GetUnpublishedEvents()
	if err != nil {
		log.Println("Error fetching events:", err)
		return
	}

	for _, event := range events {
		if err := p.publishMessage([]byte(event.Data)); err == nil {
			if err = outbox.MarkEventAsPublished(event.ID); err != nil {
				log.Printf("Failed to mark event %d as published: %v", event.ID, err)
				continue
			}
		} else {
			log.Printf("Failed to publish message for event %d: %v", event.ID, err)
		}
	}

	outbox.Cleanup()
}

func (p *KafkaProducer) publishMessage(message []byte) error {
	// Establishing connection to the leader of the partition
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		p.Writer.Addr.String(),
		p.Writer.Topic,
		0)
	if err != nil {
		log.Printf("Failed to dial leader: %v", err)
		return err
	}
	defer conn.Close()

	// Writing message to the leader
	_, err = conn.WriteMessages(
		kafka.Message{
			Key:   []byte(fmt.Sprintf("Key-%d", time.Now().Unix())),
			Value: message,
		},
	)
	if err != nil {
		log.Printf("Failed to write to leader: %v", err)
		return err
	}
	return nil
}
