package rateapi

// Fetcher defines an interface for fetching rates.
type Fetcher interface {
	FetchRate(baseCode, targetCode string) (string, error)
}
