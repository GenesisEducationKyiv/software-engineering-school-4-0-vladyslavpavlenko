package gormrepo

import (
	"errors"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository.
func NewGormUserRepository(conn *GormDB) models.UserRepository {
	return &UserRepository{DB: conn.DB}
}

// Create creates a new User record.
func (u *UserRepository) Create(email string) (uint, error) {
	user := models.User{
		Email:     email,
		CreatedAt: time.Now(),
	}
	result := u.DB.Create(&user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return 0, gorm.ErrDuplicatedKey
		}
		return 0, result.Error
	}
	return user.ID, nil
}
