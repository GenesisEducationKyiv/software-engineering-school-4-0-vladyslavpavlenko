package gormrepo

import (
	"github.com/vladyslavpavlenko/genesis-api-project/internal/dbrepo/models"
	"gorm.io/gorm"
)

type CurrencyRepository struct {
	DB *gorm.DB
}

// NewGormCurrencyRepository creates a new GormCurrencyRepository.
func NewGormCurrencyRepository(conn *GormDB) models.CurrencyRepository {
	return &CurrencyRepository{DB: conn.DB}
}

// GetIDbyCode returns the ID of the currency by its Code.
func (c *CurrencyRepository) GetIDbyCode(code string) (uint, error) {
	var currency models.Currency
	err := c.DB.Where("code = ?", code).First(&currency).Error
	if err != nil {
		return 0, err
	}

	return currency.ID, nil
}
