package jsonutils

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Error   bool   `json:"error"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	out, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/jsonutils")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

// ErrorJSON writes the error JSON. If status is not provided, the default http.StatusBadRequest in sent.
func ErrorJSON(w http.ResponseWriter, err error, status ...int) error {
	statusCode := http.StatusBadRequest

	if len(status) > 0 {
		statusCode = status[0]
	}

	payload := Response{
		Error:   true,
		Message: err.Error(),
	}

	return WriteJSON(w, statusCode, payload)
}
