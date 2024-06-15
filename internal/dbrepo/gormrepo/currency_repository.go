package gormrepo

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"
	"gorm.io/gorm"
)

type GormCurrencyRepository struct {
	db *gorm.DB
}

// NewGormCurrencyRepository creates a new GormCurrencyRepository.
func NewGormCurrencyRepository(conn *GormDB) models.CurrencyRepository {
	return &GormCurrencyRepository{db: conn.DB}
}

// GetIDbyCode returns the ID of the currency by its Code.
func (repo *GormCurrencyRepository) GetIDbyCode(code string) (uint, error) {
	var currency models.Currency
	err := repo.db.Where("code = ?", code).First(&currency).Error
	if err != nil {
		return 0, err
	}

	return currency.ID, nil
}
