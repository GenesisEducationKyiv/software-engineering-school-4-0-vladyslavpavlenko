package dbrepo

// UserRepository interface defines methods to access user data.
type UserRepository interface {
	Create(email string) (*User, error)
}
