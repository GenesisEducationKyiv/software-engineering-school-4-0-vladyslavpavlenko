package gormrepo

import (
	"errors"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository.
func NewGormUserRepository(conn *GormDB) models.UserRepository {
	return &GormUserRepository{db: conn.DB}
}

// Create creates a new User record.
func (repo *GormUserRepository) Create(email string) (*models.User, error) {
	user := models.User{
		Email:     email,
		CreatedAt: time.Now(),
	}
	result := repo.db.Create(&user)
	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) && pgErr.Code == "23505" {
			return nil, gorm.ErrDuplicatedKey
		}
		return nil, result.Error
	}
	return &user, nil
}
