package handlers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/handlers"
)

func decodeJSONResponse(t *testing.T, r *http.Response, target interface{}) {
	t.Helper()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(target); err != nil {
		t.Fatal(err)
	}
}

func TestWriteJSON(t *testing.T) {
	data := struct {
		Name string
	}{
		Name: "Test",
	}

	rr := httptest.NewRecorder()
	err := handlers.WriteJSON(rr, http.StatusOK, data)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if status := rr.Result().StatusCode; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}

	var response struct {
		Name string
	}
	decodeJSONResponse(t, rr.Result(), &response)
	if response.Name != "Test" {
		t.Errorf("Expected Name to be 'Test', got '%s'", response.Name)
	}
}

func TestErrorJSON(t *testing.T) {
	rr := httptest.NewRecorder()
	err := handlers.ErrorJSON(rr, errors.New("test error"), http.StatusInternalServerError)
	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if status := rr.Result().StatusCode; status != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, status)
	}

	var response struct {
		Error   bool
		Message string
	}
	decodeJSONResponse(t, rr.Result(), &response)
	if !response.Error || response.Message != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", response.Message)
	}
}

type badJSON struct{}

func (b badJSON) MarshalJSON() ([]byte, error) {
	return nil, errors.New("intentional marshal error")
}

func TestWriteJSON_MarshalError(t *testing.T) {
	data := badJSON{}

	rr := httptest.NewRecorder()
	err := handlers.WriteJSON(rr, http.StatusOK, data)
	if err == nil {
		t.Errorf("Expected error, but got none")
	}
}
