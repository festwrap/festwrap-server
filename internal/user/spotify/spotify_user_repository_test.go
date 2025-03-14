package spotify

import (
	"context"
	"errors"
	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	tokenKey = types.ContextKey("myKey")
	token    = "some_token"
)

func testContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, tokenKey, token)
	return ctx
}

func userIdSender() *httpsender.FakeHTTPSender {
	sender := &httpsender.FakeHTTPSender{}
	response := []byte(`{"id":"my_id"}`)
	sender.SetResponse(&response)
	return sender
}

func spotifyUserRepository(sender httpsender.HTTPRequestSender) SpotifyUserRepository {
	repository := NewSpotifyUserRepository(sender)
	repository.SetTokenKey(tokenKey)
	return repository
}

func getUserHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://api.spotify.com/v1/me"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", token)},
	)
	return options
}

func TestRepositoryMethodsReturnErrorWhenInvalidToken(t *testing.T) {
	tests := map[string]struct {
		repositoryTokenKey types.ContextKey
		tokenKey           types.ContextKey
		tokenVal           interface{}
	}{
		"returns error when token is wrong type": {
			repositoryTokenKey: "matchingKey",
			tokenKey:           "matchingKey",
			tokenVal:           1234,
		},
		"returns error when token is missing": {
			repositoryTokenKey: "someKey",
			tokenKey:           "otherKey",
			tokenVal:           "myToken",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctx = context.WithValue(ctx, test.tokenKey, test.tokenVal)
			repository := spotifyUserRepository(userIdSender())
			repository.SetTokenKey(test.repositoryTokenKey)

			_, err := repository.GetCurrentUserId(ctx)
			assert.NotNil(t, err)
		})
	}
}

func TestGetCurrentUserIdSendsRequestWithProperOptions(t *testing.T) {
	sender := userIdSender()
	repository := spotifyUserRepository(sender)

	_, err := repository.GetCurrentUserId(testContext())

	assert.Nil(t, err)
	assert.Equal(t, getUserHttpOptions(), sender.GetSendArgs())
}

func TestGetCurrentUserReturnsErrorOnSendError(t *testing.T) {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test error"))
	repository := spotifyUserRepository(sender)

	_, err := repository.GetCurrentUserId(testContext())

	assert.NotNil(t, err)
}

func TestGetCurrentUserReturnsErrorOnNonJsonUserIdBody(t *testing.T) {
	sender := userIdSender()
	nonJsonResponse := []byte("{non_json")
	sender.SetResponse(&nonJsonResponse)
	repository := spotifyUserRepository(sender)

	_, err := repository.GetCurrentUserId(testContext())

	assert.NotNil(t, err)
}

func TestGetCurrentUserReturnsUserId(t *testing.T) {
	repository := spotifyUserRepository(userIdSender())

	actual, err := repository.GetCurrentUserId(testContext())

	expected := "my_id"
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
