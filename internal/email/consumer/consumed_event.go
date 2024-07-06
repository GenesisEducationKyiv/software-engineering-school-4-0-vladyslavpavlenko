package consumer

import "time"

// ConsumedEvent represents an event consumed by the consumer, including the partition information.
type ConsumedEvent struct {
	Topic      string `gorm:"primaryKey"`
	Partition  int    `gorm:"primaryKey"`
	Offset     int64  `gorm:"primaryKey"`
	ConsumedAt time.Time
}
