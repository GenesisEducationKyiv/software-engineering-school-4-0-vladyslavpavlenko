package models

// Currency is a GORM currency model.
type Currency struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
