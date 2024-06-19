package email

import "errors"

// Config holds the email configuration.
type Config struct {
	Email    string
	Password string
}

// NewEmailConfig creates an instance of Config.
func NewEmailConfig(email, password string) (Config, error) {
	cfg := Config{
		Email:    email,
		Password: password,
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Validate validates an instance of Config to check whether the provided parameters are correct.
func (cfg *Config) Validate() error {
	if cfg.Email == "" {
		return errors.New("config email is empty")
	}

	if cfg.Password == "" {
		return errors.New("config password is empty")
	}

	if !Email(cfg.Email).Validate() {
		return errors.New("email is invalid")
	}

	return nil
}
