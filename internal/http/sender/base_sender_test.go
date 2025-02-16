package httpsender

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	httpclient "festwrap/internal/http/client"
	"festwrap/internal/testtools"

	"github.com/stretchr/testify/assert"
)

func defaultRequestBody() []byte {
	return []byte("{\"request\": \"some_request\"}")
}

func defaultOptions() HTTPRequestOptions {
	options := NewHTTPRequestOptions("https://some_url", POST, 200)
	options.SetBody(defaultRequestBody())
	return options
}

func defaultResponseBody() []byte {
	return []byte("{\"response\": \"some_response\"}")
}

func defaultResponse() *http.Response {
	return &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBuffer(defaultResponseBody())),
	}
}

func errorStatusResponse() *http.Response {
	return &http.Response{Status: "500 Unexpected error", StatusCode: 500}
}

func errorBodyResponse() *http.Response {
	return &http.Response{Status: "200 OK", StatusCode: 200, Body: testtools.NewErrorReader()}
}

func testSetup() (*httpclient.FakeHTTPClient, HTTPRequestSender, HTTPRequestOptions) {
	client := httpclient.NewFakeHTTPClient()
	client.SetResponse(defaultResponse())
	options := defaultOptions()
	sender := NewBaseHTTPRequestSender(&client)
	return &client, &sender, options
}

func TestSendRequestHasProvidedMethod(t *testing.T) {
	client, sender, options := testSetup()

	_, err := sender.Send(options)

	expected := client.GetRequestArg()
	actual := string(options.GetMethod())
	assert.Equal(t, actual, expected.Method)
	assert.Nil(t, err)
}

func TestSendRequestHasProvidedUrl(t *testing.T) {
	client, sender, options := testSetup()

	_, err := sender.Send(options)

	expected := client.GetRequestArg()
	actual := options.GetUrl()
	assert.Equal(t, actual, expected.URL.String())
	assert.Nil(t, err)
}

func TestSendRequestHasProvidedBody(t *testing.T) {
	client, sender, options := testSetup()

	_, err := sender.Send(options)

	expected := client.GetRequestArg()
	actual := options.GetBody()
	assert.Equal(t, actual, readBodyFromRequest(t, expected))
	assert.Nil(t, err)
}

func TestSendRequestDoesNotIncludeBodyIfNotProvided(t *testing.T) {
	client, sender, options := testSetup()
	options.SetBody(nil)

	_, err := sender.Send(options)

	expected := client.GetRequestArg()

	assert.Nil(t, expected.Body)
	assert.Nil(t, err)
}

func TestSendRequestUsesHeaders(t *testing.T) {
	client, sender, options := testSetup()
	headers := map[string]string{
		"Something":      "some_value",
		"Something_else": "some_other_value",
	}
	options.SetHeaders(headers)

	_, err := sender.Send(options)

	expected := client.GetRequestArg()
	assertHeadersMatch(t, headers, expected.Header)
	assert.Nil(t, err)
}

func TestSendRequestUsesNoHeadersIfNotProvided(t *testing.T) {
	client, sender, options := testSetup()

	_, err := sender.Send(options)

	expected := client.GetRequestArg()
	if len(expected.Header) > 0 {
		t.Errorf("Headers should be empty, found %v", expected.Header)
	}
	assert.Nil(t, err)
}

func TestSendRequestReturnsErrorOnErrorRequestOptions(t *testing.T) {
	_, sender, options := testSetup()
	options.SetUrl("https://bad url")

	_, err := sender.Send(options)

	assert.NotNil(t, err)
}

func TestSendRequestReturnsErrorOnClientError(t *testing.T) {
	client, sender, options := testSetup()
	client.SetError(errors.New("Test client error"))

	_, err := sender.Send(options)

	assert.NotNil(t, err)
}

func TestSendRequestReturnsErrorWhenStatusNotExpected(t *testing.T) {
	client, sender, options := testSetup()
	client.SetResponse(errorStatusResponse())

	_, err := sender.Send(options)

	assert.NotNil(t, err)
}

func TestSendRequestReturnsErrorOnResponseBodyError(t *testing.T) {
	client, sender, options := testSetup()
	client.SetResponse(errorBodyResponse())

	_, err := sender.Send(options)

	assert.NotNil(t, err)
}

func TestSendRequestReturnsResponseBody(t *testing.T) {
	_, sender, options := testSetup()

	body, err := sender.Send(options)

	assert.Equal(t, string(*body), string(defaultResponseBody()))
	assert.Nil(t, err)
}

func readBodyFromRequest(t *testing.T, request *http.Request) []byte {
	requestBody, err := io.ReadAll(request.Body)
	if err != nil {
		t.Fatal("Could not read body from request")
	}
	defer request.Body.Close()
	return requestBody
}

func assertHeadersMatch(t *testing.T, expected map[string]string, actual http.Header) {
	for name, values := range actual {
		if len(values) != 1 {
			t.Errorf("Expected a single value for header %v, found %d", name, len(values))
		}

		if values[0] != expected[name] {
			t.Errorf("Expected value %s for key %s, found %v", expected[name], name, values[0])
		}
	}
}
