package validator

// Validator defines an interface for validators.
type Validator interface {
	Validate(input string) (bool, string)
}
