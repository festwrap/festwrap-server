package spotify

import (
	"context"
	"errors"
	"fmt"
	"testing"

	types "festwrap/internal"
	httpsender "festwrap/internal/http/sender"
	"festwrap/internal/playlist"
	"festwrap/internal/serialization"
	"festwrap/internal/song"

	"github.com/stretchr/testify/assert"
)

const (
	addSongsPlaylistId  = "testId"
	createPlaylistId    = "someId"
	searchPlaylistName  = "searchPlaylist"
	searchPlaylistLimit = 5
	token               = "abcdefg12345" // gitleaks:allow
	tokenKey            = "token"
	userId              = "qrRwLBFxQL9fknW8NzBn4JprRNgS"
	userIdKey           = "userId"
)

func emptyResponseSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	emptyResponse := []byte("")
	sender.SetResponse(&emptyResponse)
	return &sender
}

func errorSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	sender.SetError(errors.New("test send error"))
	return &sender
}

func nonJsonResponseSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	nonJsonResponse := []byte("{abc}")
	sender.SetResponse(&nonJsonResponse)
	return &sender
}

func createPlaylistSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	response := fmt.Appendf(nil, `{"id":"%s"}`, createPlaylistId)
	sender.SetResponse(&response)
	return &sender
}

func searchPlaylistSender() *httpsender.FakeHTTPSender {
	sender := httpsender.FakeHTTPSender{}
	response := searchedPlaylistsResponseBody()
	sender.SetResponse(&response)
	return &sender
}

func searchedPlaylistsResponseBody() []byte {
	return []byte(`
		{
			"playlists": {
				"items": [
					{
						"id":"id1",
						"name":"first playlist",
						"description":"First description",
						"public":true,
						"owner":{"id":"qrRwLBFxQL9fknW8NzBn4JprRNgS"}
					},
					{
						"id":"id2",
						"name":"second playlist",
						"description":"Second description",
						"public":false,
						"owner":{"id":"another_owner_id"}
					}
				]
			}
		}
	`)
}

func songsToAdd() []song.Song {
	return []song.Song{song.NewSong("uri1"), song.NewSong("uri2")}
}

func playlistToCreate() playlist.Playlist {
	return playlist.Playlist{Id: createPlaylistId, Name: "my-playlist", Description: "some playlist", IsPublic: false}
}

func searchedPlaylists() []playlist.Playlist {
	return []playlist.Playlist{
		{
			Id:          "id1",
			Name:        "first playlist",
			Description: "First description",
			IsPublic:    true,
		},
	}
}

func testContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ContextKey(tokenKey), token)
	ctx = context.WithValue(ctx, types.ContextKey(userIdKey), userId)
	return ctx
}

func addSongsHttpOptions() httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", addSongsPlaylistId)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	options.SetHeaders(authHeaders())
	options.SetBody([]byte(`{"uris":["uri1","uri2"]}`))
	return options
}

func createPlaylistHttpOptions() httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userId)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	options.SetHeaders(authHeaders())
	createPlaylistBody := []byte(`{"name":"my-playlist","description":"some playlist","public":false}`)
	options.SetBody(createPlaylistBody)
	return options
}

func searchPlaylistHttpOptions() httpsender.HTTPRequestOptions {
	url := fmt.Sprintf(
		"https://api.spotify.com/v1/search?limit=%d&q=%s&type=playlist",
		searchPlaylistLimit,
		searchPlaylistName,
	)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.GET, 200)
	options.SetHeaders(authHeaders())
	return options
}

func authHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
}

func spotifyPlaylistRepository(sender httpsender.HTTPRequestSender) SpotifyPlaylistRepository {
	repository := NewSpotifyPlaylistRepository(sender)
	repository.SetTokenKey(tokenKey)
	repository.SetUserIdKey(userIdKey)
	return repository
}

func TestAddSongsReturnsErrorWhenNoSongsProvided(t *testing.T) {
	repository := spotifyPlaylistRepository(emptyResponseSender())

	err := repository.AddSongs(testContext(), addSongsPlaylistId, []song.Song{})

	assert.NotNil(t, err)
}

func TestAddSongsSendsRequestUsingProperOptions(t *testing.T) {
	sender := emptyResponseSender()
	repository := spotifyPlaylistRepository(sender)

	repository.AddSongs(testContext(), addSongsPlaylistId, songsToAdd())

	actual := sender.GetSendArgs()
	assert.Equal(t, addSongsHttpOptions(), actual)
}

func TestAddSongsReturnsErrorOnSendError(t *testing.T) {
	repository := spotifyPlaylistRepository(errorSender())

	err := repository.AddSongs(testContext(), addSongsPlaylistId, songsToAdd())

	assert.NotNil(t, err)
}

func TestAddSongsReturnsNoError(t *testing.T) {
	repository := spotifyPlaylistRepository(emptyResponseSender())

	err := repository.AddSongs(testContext(), addSongsPlaylistId, songsToAdd())

	assert.Nil(t, err)
}

func TestCreatePlaylistReturnsErrorOnPlaylistSerializationError(t *testing.T) {
	repository := spotifyPlaylistRepository(createPlaylistSender())
	serializer := serialization.FakeSerializer[SpotifyPlaylist]{}
	serializer.SetError(errors.New("test error"))
	repository.SetPlaylistCreateSerializer(&serializer)

	_, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistSendsRequestWithOptions(t *testing.T) {
	sender := createPlaylistSender()
	repository := spotifyPlaylistRepository(sender)

	repository.CreatePlaylist(testContext(), playlistToCreate())

	actual := sender.GetSendArgs()
	assert.Equal(t, createPlaylistHttpOptions(), actual)
}

func TestCreatePlaylistReturnsErrorIfSenderResponseIsNotJson(t *testing.T) {
	repository := spotifyPlaylistRepository(nonJsonResponseSender())

	_, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistReturnsErrorOnSenderError(t *testing.T) {
	repository := spotifyPlaylistRepository(errorSender())

	_, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistReturnsCreatedPlaylistId(t *testing.T) {
	repository := spotifyPlaylistRepository(createPlaylistSender())

	actual, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.Equal(t, createPlaylistId, actual)
	assert.Nil(t, err)
}

func TestSearchPlaylistSendsRequestWithOptions(t *testing.T) {
	sender := searchPlaylistSender()
	repository := spotifyPlaylistRepository(sender)

	repository.SearchPlaylist(testContext(), searchPlaylistName, searchPlaylistLimit)

	actual := sender.GetSendArgs()
	assert.Equal(t, searchPlaylistHttpOptions(), actual)
}

func TestSearchPlaylistReturnsErrorOnSendError(t *testing.T) {
	repository := spotifyPlaylistRepository(errorSender())

	_, err := repository.SearchPlaylist(testContext(), searchPlaylistName, searchPlaylistLimit)

	assert.NotNil(t, err)
}

func TestSearchPlaylistReturnsErrorIfSenderResponseIsNotJson(t *testing.T) {
	repository := spotifyPlaylistRepository(nonJsonResponseSender())

	_, err := repository.SearchPlaylist(testContext(), searchPlaylistName, searchPlaylistLimit)

	assert.NotNil(t, err)
}

func TestSearchPlaylistReturnsSearchedPlaylists(t *testing.T) {
	repository := spotifyPlaylistRepository(searchPlaylistSender())

	actual, err := repository.SearchPlaylist(testContext(), searchPlaylistName, searchPlaylistLimit)

	assert.Equal(t, searchedPlaylists(), actual)
	assert.Nil(t, err)
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
			repository := spotifyPlaylistRepository(emptyResponseSender())
			repository.SetTokenKey(test.repositoryTokenKey)

			err := repository.AddSongs(ctx, addSongsPlaylistId, songsToAdd())
			assert.NotNil(t, err)

			_, err = repository.CreatePlaylist(ctx, playlistToCreate())
			assert.NotNil(t, err)

			_, err = repository.SearchPlaylist(ctx, searchPlaylistName, searchPlaylistLimit)
			assert.NotNil(t, err)
		})
	}
}

func TestRepositoryMethodsReturnErrorWhenInvalidUserId(t *testing.T) {
	tests := map[string]struct {
		repositoryUserIdKey types.ContextKey
		userIdKey           types.ContextKey
		userIdVal           any
	}{
		"returns error when userId is wrong type": {
			repositoryUserIdKey: "matchingKey",
			userIdKey:           "matchingKey",
			userIdVal:           1234,
		},
		"returns error when userId is missing": {
			repositoryUserIdKey: "someKey",
			userIdKey:           "otherKey",
			userIdVal:           "myUser",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			ctx = context.WithValue(ctx, test.userIdKey, test.userIdVal)
			repository := spotifyPlaylistRepository(emptyResponseSender())
			repository.SetUserIdKey(test.repositoryUserIdKey)

			_, err := repository.SearchPlaylist(ctx, searchPlaylistName, searchPlaylistLimit)
			assert.NotNil(t, err)

			_, err = repository.CreatePlaylist(ctx, playlistToCreate())
			assert.NotNil(t, err)
		})
	}
}
