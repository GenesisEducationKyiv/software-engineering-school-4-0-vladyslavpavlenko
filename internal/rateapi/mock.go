package rateapi

import "net/http"

// MockHTTPClient defines the interface for a mock HTTP client.
type MockHTTPClient struct {
	Resp *http.Response
	Err  error
}

func (m *MockHTTPClient) Do(_ *http.Request) (*http.Response, error) {
	return m.Resp, m.Err
}
