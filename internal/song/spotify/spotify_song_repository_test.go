package spotify

import (
	"context"
	"errors"
	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/serialization"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func defaultToken() string {
	return "some_token"
}

func defaultTokenKey() types.ContextKey {
	return "token"
}

func defaultArtist() string {
	return "Movements"
}

func defaultTitle() string {
	return "Daylily"
}

func expectedHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://api.spotify.com/v1/search?q=%2Bartist%3AMovements%2Btrack%3ADaylily&type=track"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{"Authorization": "Bearer some_token"},
	)
	return options
}

func defaultResponse() []byte {
	return []byte(`{"tracks":{"items":[{"uri":"some uri"},{}"uri":"another uri"]}}`)
}

func integrationResponse(t *testing.T) []byte {
	return testtools.LoadTestDataOrError(
		t,
		filepath.Join(testtools.GetParentDir(t), "testdata", "search_song_response.json"),
	)
}

func defaultDeserializedResponse() spotifyResponse {
	return spotifyResponse{
		Tracks: spotifyTracks{
			Songs: []spotifySong{
				{"some uri"},
				{"another uri"},
			},
		},
	}
}

func defaultSender() *httpsender.FakeHTTPSender {
	sender := &httpsender.FakeHTTPSender{}
	response := defaultResponse()
	sender.SetResponse(&response)
	return sender
}

func defaultDeserializer() *serialization.FakeDeserializer[spotifyResponse] {
	deserializer := &serialization.FakeDeserializer[spotifyResponse]{}
	deserializer.SetResponse(defaultDeserializedResponse())
	return deserializer
}

func defaultContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, defaultTokenKey(), defaultToken())
	return ctx
}

func spotifySongRepository(sender httpsender.HTTPRequestSender) SpotifySongRepository {
	repository := NewSpotifySongRepository(sender)
	repository.SetDeserializer(defaultDeserializer())
	return *repository
}

func TestGetSongSendsRequestWithProperOptions(t *testing.T) {
	sender := defaultSender()
	repository := spotifySongRepository(sender)

	_, err := repository.GetSong(defaultContext(), defaultArtist(), defaultTitle())

	assert.Nil(t, err)
	assert.Equal(t, sender.GetSendArgs(), expectedHttpOptions())
}

func TestGetSongReturnsErrorOnSendError(t *testing.T) {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test error"))
	repository := spotifySongRepository(sender)

	_, err := repository.GetSong(defaultContext(), defaultArtist(), defaultTitle())

	assert.NotNil(t, err)
}

func TestGetSongCallsDeserializeWithSendResponseBody(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	deserializer := defaultDeserializer()
	repository.SetDeserializer(deserializer)

	_, err := repository.GetSong(defaultContext(), defaultArtist(), defaultTitle())

	assert.Nil(t, err)
	assert.Equal(t, deserializer.GetArgs(), defaultResponse())
}

func TestGetSongReturnsErrorOnResponseBodyDeserializationError(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	deserializer := defaultDeserializer()
	deserializer.SetError(errors.New("test error"))
	repository.SetDeserializer(deserializer)

	_, err := repository.GetSong(defaultContext(), defaultArtist(), defaultTitle())

	assert.NotNil(t, err)
}

func TestGetSongReturnsErrorIfNoSongsFound(t *testing.T) {
	repository := spotifySongRepository(defaultSender())
	deserializer := defaultDeserializer()
	emptyResponse := spotifyResponse{Tracks: spotifyTracks{Songs: []spotifySong{}}}
	deserializer.SetResponse(emptyResponse)
	repository.SetDeserializer(deserializer)

	_, err := repository.GetSong(defaultContext(), defaultArtist(), defaultTitle())

	assert.NotNil(t, err)
}

func TestGetSongReturnsFirstSongFound(t *testing.T) {
	repository := spotifySongRepository(defaultSender())

	actual, err := repository.GetSong(defaultContext(), defaultArtist(), defaultTitle())

	expected := song.NewSong(defaultDeserializedResponse().Tracks.Songs[0].Uri)
	assert.Nil(t, err)
	assert.Equal(t, *actual, expected)
}

func TestGetSongReturnErrorWhenInvalidToken(t *testing.T) {
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
			repository := spotifySongRepository(defaultSender())
			repository.SetTokenKey(test.repositoryTokenKey)

			_, err := repository.GetSong(ctx, defaultArtist(), defaultTitle())

			assert.NotNil(t, err)
		})
	}
}

func TestGetSongReturnsFirstSongFoundIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	sender := defaultSender()
	response := integrationResponse(t)
	sender.SetResponse(&response)
	repository := NewSpotifySongRepository(sender)

	actual, err := repository.GetSong(defaultContext(), defaultArtist(), defaultTitle())

	expected := song.NewSong("spotify:track:4rH1kFLYW0b28UNRyn7dK3")
	assert.Nil(t, err)
	assert.Equal(t, *actual, expected)
}
