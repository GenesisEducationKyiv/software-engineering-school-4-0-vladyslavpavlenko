package main

import (
	"testing"

	"github.com/go-chi/chi"
)

func TestRoutes(t *testing.T) {
	mux := routes()

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing
	default:
		t.Errorf("type is not chi.Mux, but is %T", v)
	}
}
