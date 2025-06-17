package httpsender

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	httpclient "festwrap/internal/http/client"
)

type BaseHTTPRequestSender struct {
	client httpclient.HTTPClient
}

func NewBaseHTTPRequestSender(client httpclient.HTTPClient) BaseHTTPRequestSender {
	return BaseHTTPRequestSender{client: client}
}

func (c *BaseHTTPRequestSender) Send(options HTTPRequestOptions) (*[]byte, error) {
	var body io.Reader = nil
	if options.body != nil {
		body = bytes.NewBuffer(options.body)
	}

	request, err := http.NewRequest(string(options.GetMethod()), options.GetUrl(), body)
	addHeadersToRequest(options.GetHeaders(), request)
	if err != nil {
		return nil, fmt.Errorf("could not create HTTP request for options %v: %s", options, err.Error())
	}

	response, err := c.client.Send(request)
	if err != nil {
		return nil, fmt.Errorf("error sending HTTP request for options %v: %s", options, err.Error())
	}

	if response.StatusCode != options.GetExpectedStatusCode() {
		errorMsg := fmt.Sprintf(
			"request to %s failed. Expected status code %d, found %d",
			options.url,
			options.GetExpectedStatusCode(),
			response.StatusCode,
		)
		return nil, errors.New(errorMsg)
	}

	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body for options %v: %s", options, err.Error())
	}

	return &responseBody, nil
}

func addHeadersToRequest(headers map[string]string, request *http.Request) {
	for key, name := range headers {
		request.Header.Add(key, name)
	}
}
