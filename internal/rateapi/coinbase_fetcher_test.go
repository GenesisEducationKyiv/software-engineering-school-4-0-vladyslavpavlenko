package rateapi_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/vladyslavpavlenko/genesis-api-project/internal/rateapi"
)

func TestFetchRate_NetworkError(t *testing.T) {
	mockClient := &rateapi.MockHTTPClient{
		Err: fmt.Errorf("network error"),
	}
	fetcher := rateapi.Fetcher{Client: mockClient}
	_, err := fetcher.Fetch()
	if err == nil {
		t.Errorf("Expected network error, got none")
	}
}

func TestFetchRate_NonOKStatusCode(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusInternalServerError,
		Body:       io.NopCloser(bytes.NewBufferString("Internal Server Error")),
	}
	mockClient := &rateapi.MockHTTPClient{
		Resp: mockResp,
	}
	fetcher := rateapi.Fetcher{Client: mockClient}
	_, err := fetcher.Fetch()
	if err == nil {
		t.Errorf("Expected error for non-OK status code, got none")
	}
}

func TestFetchRate_ReadBodyError(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(errReader{}),
	}
	mockClient := &rateapi.MockHTTPClient{
		Resp: mockResp,
	}
	fetcher := rateapi.Fetcher{Client: mockClient}
	_, err := fetcher.Fetch()
	if err == nil {
		t.Errorf("Expected error reading the body, got none")
	}
}

type errReader struct{}

func (e errReader) Read(_ []byte) (int, error) {
	return 0, fmt.Errorf("simulated read error")
}

func TestFetchRate_UnmarshalError(t *testing.T) {
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewBufferString(`{"invalid_json"}`)),
	}
	mockClient := &rateapi.MockHTTPClient{
		Resp: mockResp,
	}
	fetcher := rateapi.Fetcher{Client: mockClient}
	_, err := fetcher.Fetch()
	if err == nil {
		t.Errorf("Expected JSON unmarshal error, got none")
	}
}
