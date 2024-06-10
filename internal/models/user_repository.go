package models

import (
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
	"time"
)

// Create creates a new User record.
func (u *User) Create(email string) (*User, error) {
	user := User{
		Email:     email,
		CreatedAt: time.Now(),
	}
	result := db.Create(&user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, result.Error
	}
	return &user, nil
}
