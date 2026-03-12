package tidal

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
	addSongsPlaylistId = "testId"
	createPlaylistId   = "someId"
	token              = "abcdefg12345" // gitleaks:allow
	tokenKey           = "token"
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
	response := fmt.Appendf(nil, `{"data": {"id":"%s"}}`, createPlaylistId)
	sender.SetResponse(&response)
	return &sender
}

func songsToAdd() []song.Song {
	return []song.Song{song.NewSong("uri1"), song.NewSong("uri2")}
}

func playlistToCreate() playlist.PlaylistDetails {
	return playlist.PlaylistDetails{Name: "my-playlist", Description: "some playlist", IsPublic: false}
}

func testContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ContextKey(tokenKey), token)
	return ctx
}

func addSongsHttpOptions() httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://openapi.tidal.com/v2/playlists/%s/relationships/items", addSongsPlaylistId)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	options.SetHeaders(authHeaders())
	options.SetBody([]byte(`{"data":[{"id":"uri1","type":"tracks"},{"id":"uri2","type":"tracks"}]}`))
	return options
}

func createPlaylistHttpOptions() httpsender.HTTPRequestOptions {
	url := "https://openapi.tidal.com/v2/playlists"
	options := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	options.SetHeaders(authHeaders())
	createPlaylistBody := []byte(`{"data":{"accessType":"PUBLIC","description":"some playlist","name":"my-playlist"},"type":"playlist"}`)
	options.SetBody(createPlaylistBody)
	return options
}

func authHeaders() map[string]string {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", token),
		"Content-Type":  "application/json",
	}
}

func tidalPlaylistRepository(sender httpsender.HTTPRequestSender) TidalPlaylistRepository {
	repository := NewTidalPlaylistRepository(sender)
	repository.SetTokenKey(tokenKey)
	return repository
}

func TestAddSongsReturnsErrorWhenNoSongsProvided(t *testing.T) {
	repository := tidalPlaylistRepository(emptyResponseSender())

	err := repository.AddSongs(testContext(), addSongsPlaylistId, []song.Song{})

	assert.NotNil(t, err)
}

func TestAddSongsSendsRequestUsingProperOptions(t *testing.T) {
	sender := emptyResponseSender()
	repository := tidalPlaylistRepository(sender)

	repository.AddSongs(testContext(), addSongsPlaylistId, songsToAdd())

	actual := sender.GetSendArgs()
	assert.Equal(t, addSongsHttpOptions(), actual)
}

func TestAddSongsReturnsErrorOnSendError(t *testing.T) {
	repository := tidalPlaylistRepository(errorSender())

	err := repository.AddSongs(testContext(), addSongsPlaylistId, songsToAdd())

	assert.NotNil(t, err)
}

func TestAddSongsReturnsNoError(t *testing.T) {
	repository := tidalPlaylistRepository(emptyResponseSender())

	err := repository.AddSongs(testContext(), addSongsPlaylistId, songsToAdd())

	assert.Nil(t, err)
}

func TestCreatePlaylistReturnsErrorOnPlaylistSerializationError(t *testing.T) {
	repository := tidalPlaylistRepository(createPlaylistSender())
	serializer := serialization.FakeSerializer[tidalPlaylist]{}
	serializer.SetError(errors.New("test error"))
	repository.SetPlaylistCreateSerializer(&serializer)

	_, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistSendsRequestWithOptions(t *testing.T) {
	sender := createPlaylistSender()
	repository := tidalPlaylistRepository(sender)

	repository.CreatePlaylist(testContext(), playlistToCreate())

	actual := sender.GetSendArgs()
	assert.Equal(t, createPlaylistHttpOptions(), actual)
}

func TestCreatePlaylistReturnsErrorIfSenderResponseIsNotJson(t *testing.T) {
	repository := tidalPlaylistRepository(nonJsonResponseSender())

	_, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistReturnsErrorOnSenderError(t *testing.T) {
	repository := tidalPlaylistRepository(errorSender())

	_, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistReturnsCreatedPlaylistId(t *testing.T) {
	repository := tidalPlaylistRepository(createPlaylistSender())

	actual, err := repository.CreatePlaylist(testContext(), playlistToCreate())

	assert.Equal(t, createPlaylistId, actual)
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
			repository := tidalPlaylistRepository(emptyResponseSender())
			repository.SetTokenKey(test.repositoryTokenKey)

			err := repository.AddSongs(ctx, addSongsPlaylistId, songsToAdd())
			assert.NotNil(t, err)

			_, err = repository.CreatePlaylist(ctx, playlistToCreate())
			assert.NotNil(t, err)
		})
	}
}

// expected : httpsender.HTTPRequestOptions{body:[]uint8{0x7b, 0x22, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x7b, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x6d, 0x79, 0x2d, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x2c, 0x22, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x2c, 0x22, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x22, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x22, 0x7d, 0x2c, 0x22, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x22, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x7d}, url:"https://openapi.tidal.com/v2/playlists", method:"POST", headers:map[string]string{"Authorization":"Bearer abcdefg12345", "Content-Type":"application/json"}, expectedStatusCode:201}
// actual   : httpsender.HTTPRequestOptions{body:[]uint8{0x7b, 0x22, 0x64, 0x61, 0x74, 0x61, 0x22, 0x3a, 0x7b, 0x22, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x22, 0x50, 0x55, 0x42, 0x4c, 0x49, 0x43, 0x22, 0x2c, 0x22, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x3a, 0x22, 0x73, 0x6f, 0x6d, 0x65, 0x20, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x2c, 0x22, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x3a, 0x22, 0x6d, 0x79, 0x2d, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x7d, 0x2c, 0x22, 0x74, 0x79, 0x70, 0x65, 0x22, 0x3a, 0x22, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x22, 0x7d}, url:"https://openapi.tidal.com/v2/playlists", method:"POST", headers:map[string]string{"Authorization":"Bearer abcdefg12345", "Content-Type":"application/json"}, expectedStatusCode:201}
