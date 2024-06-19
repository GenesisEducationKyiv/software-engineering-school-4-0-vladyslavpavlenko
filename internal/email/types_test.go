package email_test

import (
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

func TestEmail_Validate(t *testing.T) {
	cases := []struct {
		email string
		valid bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
	}

	for _, tc := range cases {
		e := email.Email(tc.email)
		if e.Validate() != tc.valid {
			t.Errorf("Expected %v for email validation of %s", tc.valid, tc.email)
		}
	}
}
