package gormrepo

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo"
	"gorm.io/gorm"
)

type GormCurrencyRepository struct {
	db *gorm.DB
}

// NewGormCurrencyRepository creates a new GormCurrencyRepository.
func NewGormCurrencyRepository(db *gorm.DB) dbrepo.CurrencyRepository {
	return &GormCurrencyRepository{db: db}
}

// GetIDbyCode returns the ID of the currency by its Code.
func (repo *GormCurrencyRepository) GetIDbyCode(code string) (uint, error) {
	var currency dbrepo.Currency
	err := repo.db.Where("code = ?", code).First(&currency).Error
	if err != nil {
		return 0, err
	}

	return currency.ID, nil
}
