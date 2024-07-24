package httpclient

import (
	"net/http"
)

type FakeHTTPClient struct {
	requestArg *http.Request
	response   *http.Response
	err        error
}

func (c *FakeHTTPClient) Send(request *http.Request) (*http.Response, error) {
	c.requestArg = request

	if c.err != nil {
		return nil, c.err
	}

	return c.response, nil
}

func (c *FakeHTTPClient) GetRequestArg() *http.Request {
	return c.requestArg
}

func (c *FakeHTTPClient) SetResponse(response *http.Response) {
	c.response = response
}

func (c *FakeHTTPClient) SetError(err error) {
	c.err = err
}

func NewFakeHTTPClient() FakeHTTPClient {
	return FakeHTTPClient{}
}
