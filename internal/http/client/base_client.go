package httpclient

import (
	"net/http"
)

type BaseHTTPClient struct {
	client *http.Client
}

func NewBaseHTTPClient(client *http.Client) BaseHTTPClient {
	return BaseHTTPClient{client: client}
}

func (c *BaseHTTPClient) Send(request *http.Request) (*http.Response, error) {
	return c.client.Do(request)
}
