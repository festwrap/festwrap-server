package spotify

import (
	"fmt"
	"net/http"
	"testing"

	httpsender "festwrap/internal/http/sender"
	httpsendermocks "festwrap/internal/http/sender/mocks"

	"github.com/stretchr/testify/assert"
)

const (
	refreshToken        = "cached_token"
	clientId            = "some_client_id"
	clientSecret        = "some_client_secret"
	encodedIdAndSecret  = "c29tZV9jbGllbnRfaWQ6c29tZV9jbGllbnRfc2VjcmV0" // gitleaks:allow
	authResponse        = `{"access_token": "new_token", "expires_in": 25}`
	authExpiredResponse = `{"access_token": "another_token", "expires_in": 0}`
)

func expectedSenderArgs() httpsender.HTTPRequestOptions {
	headers := map[string]string{
		"Content-Type":  "application/x-www-form-urlencoded",
		"Authorization": fmt.Sprintf("Basic %s", encodedIdAndSecret),
	}
	url := "https://accounts.spotify.com/api/token?client_id=some_client_id&grant_type=refresh_token&refresh_token=cached_token"
	opts := httpsender.NewHTTPRequestOptions(url, httpsender.POST, http.StatusOK)
	opts.SetHeaders(headers)
	return opts
}

func createSender(response string) *httpsendermocks.HTTPSenderMock {
	sender := httpsendermocks.HTTPSenderMock{}
	responseBytes := []byte(response)
	sender.On("Send", expectedSenderArgs()).Return(&responseBytes, nil)
	return &sender
}

func TestAccessTokenReturnedFromSender(t *testing.T) {
	sender := createSender(authResponse)
	client := NewSpotifyAuthClient(sender, refreshToken, clientId, clientSecret)

	token, err := client.GetAccessToken()

	assert.Nil(t, err)
	assert.Equal(t, "new_token", token)
	sender.AssertExpectations(t)
}

func TestAccessTokenReturnedFromCacheWhenPreviousNotExpired(t *testing.T) {
	sender := createSender(authResponse)
	client := NewSpotifyAuthClient(sender, refreshToken, clientId, clientSecret)
	// Obtain token that updates expiration time
	previousToken, _ := client.GetAccessToken()

	nextToken, err := client.GetAccessToken()

	assert.Nil(t, err)
	assert.Equal(t, previousToken, nextToken)
	sender.AssertNumberOfCalls(t, "Send", 1)
}

func TestAccessTokenReturnedFromSenderWhenPreviousExpired(t *testing.T) {
	sender := createSender(authExpiredResponse)
	client := NewSpotifyAuthClient(sender, refreshToken, clientId, clientSecret)
	// Obtain token that expires immediately
	client.GetAccessToken()

	_, err := client.GetAccessToken()

	assert.Nil(t, err)
	sender.AssertNumberOfCalls(t, "Send", 2)
}
