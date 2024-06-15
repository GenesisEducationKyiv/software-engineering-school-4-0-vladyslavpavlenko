package models

// CurrencyRepository interface defines methods to access Currency data.
type CurrencyRepository interface {
	GetIDbyCode(string) (uint, error)
}
