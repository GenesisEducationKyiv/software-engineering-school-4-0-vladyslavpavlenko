package models

import (
	"gorm.io/gorm"
)

type CurrencyRepository struct {
	db *gorm.DB
}

// NewCurrencyRepository creates a new CurrencyRepository.
func NewCurrencyRepository(db *gorm.DB) *CurrencyRepository {
	return &CurrencyRepository{db: db}
}

// GetIDbyCode returns the ID of the currency by its Code.
func (repo *CurrencyRepository) GetIDbyCode(code string) (uint, error) {
	var currency Currency
	err := repo.db.Where("code = ?", code).First(&currency).Error
	if err != nil {
		return 0, err
	}

	return currency.ID, nil
}
