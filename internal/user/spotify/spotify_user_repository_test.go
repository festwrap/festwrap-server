package spotify

import (
	"context"
	"errors"
	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/testtools"
	"testing"
)

func defaultTokenKey() types.ContextKey {
	return "myKey"
}

func defaultContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, defaultTokenKey(), "some_token")
	return ctx
}

func defaultResponse() []byte {
	return []byte(`{"id":"my_id"}`)
}

func defaultDeserializedResponse() spotifyUserResponse {
	return spotifyUserResponse{UserId: "userId"}
}

func defaultSender() *httpsender.FakeHTTPSender {
	sender := &httpsender.FakeHTTPSender{}
	response := defaultResponse()
	sender.SetResponse(&response)
	return sender
}

func defaultDeserializer() *serialization.FakeDeserializer[spotifyUserResponse] {
	deserializer := &serialization.FakeDeserializer[spotifyUserResponse]{}
	deserializer.SetResponse(defaultDeserializedResponse())
	return deserializer
}

func spotifyUserRepository(sender httpsender.HTTPRequestSender) SpotifyUserRepository {
	repository := NewSpotifyUserRepository(sender)
	repository.SetDeserializer(defaultDeserializer())
	repository.SetTokenKey(defaultTokenKey())
	return repository
}

func expectedHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://api.spotify.com/v1/me"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{"Authorization": "Bearer some_token"},
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
			repository := spotifyUserRepository(defaultSender())
			repository.SetTokenKey(test.repositoryTokenKey)

			_, err := repository.GetCurrentUserId(ctx)
			testtools.AssertErrorIsNotNil(t, err)
		})
	}
}

func TestGetCurrentUserIdSendsRequestWithProperOptions(t *testing.T) {
	sender := defaultSender()
	repository := spotifyUserRepository(sender)

	_, err := repository.GetCurrentUserId(defaultContext())

	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, sender.GetSendArgs(), expectedHttpOptions())
}

func TestGetCurrentUserReturnsErrorOnSendError(t *testing.T) {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test error"))
	repository := spotifyUserRepository(sender)

	_, err := repository.GetCurrentUserId(defaultContext())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetCurrentUserCallsDeserializeWithSendResponseBody(t *testing.T) {
	repository := spotifyUserRepository(defaultSender())

	_, err := repository.GetCurrentUserId(defaultContext())

	testtools.AssertErrorIsNil(t, err)
	deserializer := repository.GetDeserializer().(*serialization.FakeDeserializer[spotifyUserResponse])
	testtools.AssertEqual(t, deserializer.GetArgs(), defaultResponse())
}

func TestGetCurrentUserReturnsErrorOnResponseBodyDeserializationError(t *testing.T) {
	deserializer := defaultDeserializer()
	deserializer.SetError(errors.New("test error"))
	repository := spotifyUserRepository(defaultSender())
	repository.SetDeserializer(deserializer)

	_, err := repository.GetCurrentUserId(defaultContext())

	testtools.AssertErrorIsNotNil(t, err)
}

func TestGetCurrentUserReturnsFirstSongFoundIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	repository := spotifyUserRepository(defaultSender())
	repository.SetDeserializer(serialization.NewJsonDeserializer[spotifyUserResponse]())

	actual, err := repository.GetCurrentUserId(defaultContext())

	expected := "my_id"
	testtools.AssertErrorIsNil(t, err)
	testtools.AssertEqual(t, actual, expected)
}
