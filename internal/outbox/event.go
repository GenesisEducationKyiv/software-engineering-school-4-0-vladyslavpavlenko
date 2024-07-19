package outbox

import (
	"encoding/json"
	"time"
)

// Event is a query message model stored in the database.
type Event struct {
	ID        uint `gorm:"primaryKey"`
	Data      string
	CreatedAt time.Time
}

// Data is an event data model.
type Data struct {
	Email string  `json:"email"`
	Rate  float64 `json:"rate"`
}

// Serialize takes a Data struct and serializes it to a JSON string.
func (d Data) Serialize() (string, error) {
	bytes, err := json.Marshal(d)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
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
