package dbrepo

// CurrencyRepository interface defines methods to access currency data.
type CurrencyRepository interface {
	GetIDbyCode(string) (uint, error)
}
