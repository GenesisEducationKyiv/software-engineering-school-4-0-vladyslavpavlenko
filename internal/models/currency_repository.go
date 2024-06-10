package models

// GetIDbyCode returns the ID of the currency by its Code.
func (c *Currency) GetIDbyCode(code string) (uint, error) {
	var currency Currency
	err := db.Where("code = ?", code).First(&currency).Error
	if err != nil {
		return 0, err
	}

	return currency.ID, nil
}
