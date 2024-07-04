package outbox

import (
	"encoding/json"
	"time"
)

// Event is a query message model stored in the database.
type Event struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Data      string    `json:"data"`
	Published bool      `json:"published"`
	CreatedAt time.Time `json:"created_at"`
}

// Data is an event data model.
type Data struct {
	Email string  `json:"email"`
	Rate  float64 `json:"rate"`
}

// SerializeData takes a Data struct and serializes it to a JSON string for storage.
func (e *Event) SerializeData(data Data) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	e.Data = string(bytes)
	return nil
}

// DeserializeData deserializes JSON string to Data struct
func DeserializeData(jsonData []byte) (Data, error) {
	var data Data
	err := json.Unmarshal(jsonData, &data)
	if err != nil {
		return Data{}, err
	}
	return data, nil
}
