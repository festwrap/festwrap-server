package spotify

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	types "festwrap/internal"
	"festwrap/internal/artist"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/testtools"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	searchName = "Movements"
	limit      = 2
	tokenKey   = types.ContextKey("myKey")
	authToken  = "some_token"
)

func testContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, tokenKey, authToken)
	return ctx
}

func searchArtistHttpOptions() httpsender.HTTPRequestOptions {
	url := fmt.Sprintf(
		"https://api.spotify.com/v1/search?limit=%d&q=artist%%3A%s&type=artist",
		limit,
		searchName,
	)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(
		map[string]string{"Authorization": fmt.Sprintf("Bearer %s", authToken)},
	)
	return options
}

func artistSearchResponse(t *testing.T) *[]byte {
	t.Helper()

	path := filepath.Join(testtools.GetParentDir(t), "testdata", "spotify_artist_search_response.json")
	result := testtools.LoadTestDataOrError(t, path)
	return &result
}

func artistSearchEmptyResponse(t *testing.T) *[]byte {
	t.Helper()

	path := filepath.Join(testtools.GetParentDir(t), "testdata", "spotify_artist_search_empty_response.json")
	result := testtools.LoadTestDataOrError(t, path)
	return &result
}

func searchedArtists() []artist.Artist {
	return []artist.Artist{
		artist.NewArtistWithImageUri(
			"The Beatles",
			"https://i.scdn.co/image/ab6761610000f178e9348cc01ff5d55971b22433",
		),
		artist.NewArtistWithImageUri(
			"The Beatles Tribute Band",
			"https://i.scdn.co/image/ab67616d00004851a53d58fac4e46d5264adc122",
		),
		artist.NewArtist("The Beatles Recovered Band"),
		artist.NewArtistWithImageUri(
			"The Beatles Greatest Hits Performed By The Frank Berman Band",
			"https://i.scdn.co/image/ab67616d00004851f903d75acdce7727b3c4aa2c",
		),
		artist.NewArtist("The Beatles Revival Band"),
	}
}

func sender(t *testing.T) *httpsender.FakeHTTPSender {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetResponse(artistSearchResponse(t))
	return sender
}

func spotifySongRepository(sender httpsender.HTTPRequestSender) SpotifyArtistRepository {
	repository := NewSpotifyArtistRepository(sender)
	repository.SetTokenKey(tokenKey)
	return repository
}

func TestSearchArtistSendsRequestWithProperOptions(t *testing.T) {
	testSender := sender(t)
	repository := spotifySongRepository(testSender)

	_, err := repository.SearchArtist(testContext(), searchName, limit)

	assert.Nil(t, err)
	assert.Equal(t, searchArtistHttpOptions(), testSender.GetSendArgs())
}

func TestSearchArtistReturnsErrorOnWrongKeyType(t *testing.T) {
	ctx := testContext()
	ctx = context.WithValue(ctx, tokenKey, 42)
	repository := spotifySongRepository(sender(t))

	_, err := repository.SearchArtist(ctx, searchName, limit)

	assert.NotNil(t, err)
}

func TestSearchArtistReturnsErrorOnSendError(t *testing.T) {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test error"))
	repository := spotifySongRepository(sender)

	_, err := repository.SearchArtist(testContext(), searchName, limit)

	assert.NotNil(t, err)
}

func TestSearchArtistsReturnsErrorOnInvalidBody(t *testing.T) {
	sender := &httpsender.FakeHTTPSender{}
	invalidBody := []byte("{some_invalid_json}")
	sender.SetResponse(&invalidBody)
	repository := spotifySongRepository(sender)

	_, err := repository.SearchArtist(testContext(), searchName, limit)

	assert.NotNil(t, err)
}

func TestSearchArtistReturnsEmptyIfNoneFound(t *testing.T) {
	sender := &httpsender.FakeHTTPSender{}
	sender.SetResponse(artistSearchEmptyResponse(t))
	repository := spotifySongRepository(sender)

	artists, _ := repository.SearchArtist(testContext(), searchName, limit)

	assert.Equal(t, []artist.Artist{}, artists)
}

func TestSearchArtistReturnsArtists(t *testing.T) {
	repository := spotifySongRepository(sender(t))

	artists, _ := repository.SearchArtist(testContext(), searchName, limit)

	assert.Equal(t, artists, searchedArtists())
}
