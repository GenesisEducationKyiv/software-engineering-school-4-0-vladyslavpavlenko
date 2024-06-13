package validator

import (
	"regexp"
)

type CurrencyCodeValidator struct{}

// Validate checks if the given currency code conforms to the standard format, which consists
// of three uppercase letters.
func (v *CurrencyCodeValidator) Validate(code string) bool {
	_, err := regexp.MatchString("^[A-Z]{3}$", code)
	return err == nil
}
