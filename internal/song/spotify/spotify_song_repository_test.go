package spotify

import (
	"context"
	"errors"
	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/song"
	"festwrap/internal/testtools"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	token     = "some_token"
	tokenKey  = types.ContextKey("token")
	artist    = "Movements"
	songTitle = "Daylily"
)

func getSongHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://api.spotify.com/v1/search?q=artist%3AMovements+track%3ADaylily&type=track"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{"Authorization": "Bearer some_token"},
	)
	return options
}

func searchSongResponseBody(t *testing.T) []byte {
	return testtools.LoadTestDataOrError(
		t,
		filepath.Join(testtools.GetParentDir(t), "testdata", "search_song_response.json"),
	)
}

func noSongsSearchSongResponseBody(t *testing.T) []byte {
	return testtools.LoadTestDataOrError(
		t,
		filepath.Join(testtools.GetParentDir(t), "testdata", "no_songs_search_song_response.json"),
	)
}

func songsSender(t *testing.T) *httpsender.FakeHTTPSender {
	sender := &httpsender.FakeHTTPSender{}
	response := searchSongResponseBody(t)
	sender.SetResponse(&response)
	return sender
}

func testContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, tokenKey, token)
	return ctx
}

func TestGetSongSendsRequestWithProperOptions(t *testing.T) {
	sender := songsSender(t)
	repository := NewSpotifySongRepository(sender)

	_, err := repository.GetSong(testContext(), artist, songTitle)

	assert.Nil(t, err)
	assert.Equal(t, getSongHttpOptions(), sender.GetSendArgs())
}

func TestGetSongReturnsErrorOnSendError(t *testing.T) {
	errorSender := songsSender(t)
	errorSender.SetError(errors.New("test error"))
	repository := NewSpotifySongRepository(errorSender)

	_, err := repository.GetSong(testContext(), artist, songTitle)

	assert.NotNil(t, err)
}

func TestGetSongReturnsErrorOnNonJsonSearchResponseBody(t *testing.T) {
	sender := songsSender(t)
	nonJsonBody := []byte("{some_non_json")
	sender.SetResponse(&nonJsonBody)
	repository := NewSpotifySongRepository(sender)

	_, err := repository.GetSong(testContext(), artist, songTitle)

	assert.NotNil(t, err)
}

func TestGetSongReturnsErrorIfNoSongsFound(t *testing.T) {
	sender := songsSender(t)
	noSongsBody := noSongsSearchSongResponseBody(t)
	sender.SetResponse(&noSongsBody)
	repository := NewSpotifySongRepository(sender)

	_, err := repository.GetSong(testContext(), artist, songTitle)

	assert.NotNil(t, err)
}

func TestGetSongReturnsFirstSongFound(t *testing.T) {
	repository := NewSpotifySongRepository(songsSender(t))

	actual, err := repository.GetSong(testContext(), artist, songTitle)

	expected := song.NewSong("spotify:track:4rH1kFLYW0b28UNRyn7dK3")
	assert.Equal(t, expected, *actual)
	assert.Nil(t, err)
}

func TestGetSongReturnsErrorWhenInvalidToken(t *testing.T) {
	tests := map[string]struct {
		repositoryTokenKey types.ContextKey
		tokenKey           types.ContextKey
		tokenVal           any
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
			repository := NewSpotifySongRepository(songsSender(t))
			repository.SetTokenKey(test.repositoryTokenKey)

			_, err := repository.GetSong(ctx, artist, songTitle)

			assert.NotNil(t, err)
		})
	}
}
