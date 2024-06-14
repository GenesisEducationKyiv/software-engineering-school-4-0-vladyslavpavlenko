package gormrepo

import (
	"errors"
	"time"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository creates a new GormUserRepository.
func NewGormUserRepository(db *gorm.DB) dbrepo.UserRepository {
	return &GormUserRepository{db: db}
}

// Create creates a new User record.
func (repo *GormUserRepository) Create(email string) (*dbrepo.User, error) {
	user := dbrepo.User{
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
