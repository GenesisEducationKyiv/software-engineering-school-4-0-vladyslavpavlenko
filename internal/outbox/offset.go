package outbox

// Offset represents the last published Event for a topic.
type Offset struct {
	Topic  string `gorm:"primaryKey"`
	Offset int64
}
