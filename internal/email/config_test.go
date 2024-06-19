package email_test

import (
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/email"
)

func TestConfig_Validate(t *testing.T) {
	cases := []struct {
		name    string
		cfg     email.Config
		isValid bool
	}{
		{
			name: "valid config",
			cfg: email.Config{
				Email:    "test@example.com",
				Password: "password",
			},
			isValid: true,
		},
		{
			name: "empty email",
			cfg: email.Config{
				Password: "password",
			},
		},
		{
			name: "empty password",
			cfg: email.Config{
				Email: "test@example.com",
			},
		},
		{
			name: "invalid email",
			cfg: email.Config{
				Email:    "bad-email",
				Password: "password",
			},
		},
	}

	for _, tc := range cases {
		err := tc.cfg.Validate()
		if (err == nil) != tc.isValid {
			t.Errorf("%s: expected valid=%v, got error: %v", tc.name, tc.isValid, err)
		}
	}
}

func TestNewEmailConfig(t *testing.T) {
	cases := []struct {
		name    string
		cfg     email.Config
		isValid bool
	}{
		{
			name: "valid config",
			cfg: email.Config{
				Email:    "test@example.com",
				Password: "password",
			},
			isValid: true,
		},
		{
			name: "empty email",
			cfg: email.Config{
				Password: "password",
			},
		},
		{
			name: "empty password",
			cfg: email.Config{
				Email: "test@example.com",
			},
		},
		{
			name: "invalid email",
			cfg: email.Config{
				Email:    "bad-email",
				Password: "password",
			},
		},
	}

	for _, tc := range cases {
		_, err := email.NewEmailConfig(tc.cfg.Email, tc.cfg.Password)
		if (err == nil) != tc.isValid {
			t.Errorf("%s: expected valid=%v, got error: %v", tc.name, tc.isValid, err)
		}
	}
}
