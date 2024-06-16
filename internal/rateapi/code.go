package rateapi

import "regexp"

type Code string

// Validate checks if the given currency code conforms to the standard format, which consists of three uppercase letters.
func (c Code) Validate() bool {
	match, _ := regexp.MatchString("^[A-Za-z]{3}$", string(c))
	return match
}
