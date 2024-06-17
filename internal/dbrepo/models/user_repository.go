package models

// UserRepository interface defines methods to access User data.
type UserRepository interface {
	Create(email string) (uint, error)
}
