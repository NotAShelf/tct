package tct

import (
	"net/http"
	"testing"
)

type MockHttpClient struct {
	GetFunc func(url string) (*http.Response, error)
}

func (m *MockHttpClient) Get(url string) (*http.Response, error) {
	return m.GetFunc(url)
}

func TestClientMakeRequest(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{
			name: "valid URL",
			url:  "http://example.com",
		},
		{
			name: "another valid URL",
			url:  "https://google.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockHttpClient{
				GetFunc: func(url string) (*http.Response, error) {
					return &http.Response{
						StatusCode: 200,
						Body:       http.NoBody,
					}, nil
				},
			}

			client := &Client{httpClient: mockClient}
			client.MakeRequest(tt.url)
		})
	}
}
