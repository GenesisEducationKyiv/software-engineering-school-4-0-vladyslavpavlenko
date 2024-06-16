package rateapi_test

import (
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
)

func TestCode_Validate(t *testing.T) {
	tests := []struct {
		name string
		code rateapi.Code
		want bool
	}{
		{"valid code", "USD", true},
		{"lowercase code", "usd", true},
		{"invalid code", "US", false},
		{"numeric code", "123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.code.Validate(); got != tt.want {
				t.Errorf("Code.Validate() = %v, want %v", got, tt.want)
			}
		})
	}
}
