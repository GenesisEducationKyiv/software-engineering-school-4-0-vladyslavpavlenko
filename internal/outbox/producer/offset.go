package producer

// Offset represents the last published event offset for a topic and partition.
type Offset struct {
	Topic     string `gorm:"primaryKey"`
	Partition int    `gorm:"primaryKey"`
	Offset    uint
}
