package tct

import (
	"fmt"
	"io"
	"net/http"
)

type HttpClient interface {
	Get(url string) (*http.Response, error)
}

type Client struct {
	httpClient HttpClient
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
	}
}

func (c *Client) MakeRequest(url string) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	_ = body
}
