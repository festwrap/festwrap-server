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
	"festwrap/internal/testtools"

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

func fakeSender() *httpsender.FakeHTTPSender {
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

func songsToAdd() []song.Song {
	return []song.Song{song.NewSong("uri1"), song.NewSong("uri2")}
}

func songsToAddResponseBody() []byte {
	return []byte(`{"uris":["uri1","uri2"]}`)
}

func playlistToCreate() playlist.Playlist {
	return playlist.Playlist{Name: "my-playlist", Description: "some playlist", IsPublic: false}
}

func createPlaylistBody() []byte {
	return []byte(`{"name":"my-playlist","description":"some playlist","is_public":false}`)
}

func searchedPlaylistsResponseBody() []byte {
	return []byte(`
		{
			"playlists": {
				"items": [
					{
						"name":"first playlist",
						"description":"First description",
						"public":true,
						"owner":{"id":"qrRwLBFxQL9fknW8NzBn4JprRNgS"}
					},
					{
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

func searchedPlaylists() SpotifySearchPlaylistResponse {
	return SpotifySearchPlaylistResponse{
		Playlists: SpotifySearchPlaylists{
			Items: []SpotifySearchPlaylist{
				{
					Name:          "first playlist",
					Description:   "First description",
					Public:        true,
					OwnerMetadata: SpotifyPlaylistOwnerMetadata{Id: userId},
				},
				{
					Name:          "second playlist",
					Description:   "Second description",
					Public:        false,
					OwnerMetadata: SpotifyPlaylistOwnerMetadata{Id: "another_owner_id"},
				},
			},
		},
	}
}

func createPlaylistResponseBody() []byte {
	return []byte(fmt.Sprintf(`{"id":"%s"}`, createPlaylistId))
}

func createdPlaylist() SpotifyCreatePlaylistResponse {
	return SpotifyCreatePlaylistResponse{Id: createPlaylistId}
}

func searchedFilteredPlaylists() []playlist.Playlist {
	return []playlist.Playlist{
		{
			Name:        "first playlist",
			Description: "First description",
			IsPublic:    true,
		},
	}
}

func songSerializer() *serialization.FakeSerializer[SpotifySongs] {
	serializer := serialization.FakeSerializer[SpotifySongs]{}
	serializer.SetResponse(songsToAddResponseBody())
	return &serializer
}

func playlistCreateSerializer() *serialization.FakeSerializer[SpotifyPlaylist] {
	serializer := serialization.FakeSerializer[SpotifyPlaylist]{}
	serializer.SetResponse(createPlaylistBody())
	return &serializer
}

func playlistCreateDeserializer() *serialization.FakeDeserializer[SpotifyCreatePlaylistResponse] {
	deserializer := serialization.FakeDeserializer[SpotifyCreatePlaylistResponse]{}
	deserializer.SetResponse(createdPlaylist())
	return &deserializer
}

func playlistSearchDeserializer() *serialization.FakeDeserializer[SpotifySearchPlaylistResponse] {
	deserializer := serialization.FakeDeserializer[SpotifySearchPlaylistResponse]{}
	deserializer.SetResponse(searchedPlaylists())
	return &deserializer
}

func defaultContext() context.Context {
	ctx := context.Background()
	ctx = context.WithValue(ctx, types.ContextKey(tokenKey), token)
	ctx = context.WithValue(ctx, types.ContextKey(userIdKey), userId)
	return ctx
}

func expectedAddSongsHttpOptions() httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks", addSongsPlaylistId)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	options.SetHeaders(authHeaders())
	options.SetBody(songsToAddResponseBody())
	return options
}

func expectedCreatePlaylistHttpOptions() httpsender.HTTPRequestOptions {
	url := fmt.Sprintf("https://api.spotify.com/v1/users/%s/playlists", userId)
	options := httpsender.NewHTTPRequestOptions(url, httpsender.POST, 201)
	options.SetHeaders(authHeaders())
	options.SetBody(createPlaylistBody())
	return options
}

func expectedSearchPlaylistHttpOptions() httpsender.HTTPRequestOptions {
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

func spotifyPlaylistRepository() SpotifyPlaylistRepository {
	repository := NewSpotifyPlaylistRepository(fakeSender())

	repository.SetTokenKey(tokenKey)
	repository.SetUserIdKey(userIdKey)

	repository.SetSongSerializer(songSerializer())
	repository.SetPlaylistCreateSerializer(playlistCreateSerializer())
	repository.SetPlaylistSearchDeserializer(playlistSearchDeserializer())
	repository.SetPlaylistCreateDeserializer((playlistCreateDeserializer()))

	return repository
}

func TestAddSongsReturnsErrorWhenNoSongsProvided(t *testing.T) {
	repository := spotifyPlaylistRepository()

	err := repository.AddSongs(defaultContext(), addSongsPlaylistId, []song.Song{})

	assert.NotNil(t, err)
}

func TestAddSongsSerializesInputSongs(t *testing.T) {
	repository := spotifyPlaylistRepository()

	repository.AddSongs(defaultContext(), addSongsPlaylistId, songsToAdd())

	expected := SpotifySongs{Uris: []string{"uri1", "uri2"}}
	actual := repository.GetSongSerializer().(*serialization.FakeSerializer[SpotifySongs]).GetArgs()
	assert.Equal(t, actual, expected)
}

func TestAddSongsReturnsErrorOnNonSerializationError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	serializer := songSerializer()
	serializer.SetError(errors.New("test songs error"))
	repository.SetSongSerializer(serializer)

	err := repository.AddSongs(defaultContext(), addSongsPlaylistId, songsToAdd())

	assert.NotNil(t, err)
}

func TestAddSongsSendsRequestUsingProperOptions(t *testing.T) {
	repository := spotifyPlaylistRepository()

	repository.AddSongs(defaultContext(), addSongsPlaylistId, songsToAdd())

	actual := repository.GetHTTPSender().(*httpsender.FakeHTTPSender).GetSendArgs()
	assert.Equal(t, actual, expectedAddSongsHttpOptions())
}

func TestAddSongsReturnsErrorOnSendError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(errorSender())

	err := repository.AddSongs(defaultContext(), addSongsPlaylistId, songsToAdd())

	assert.NotNil(t, err)
}

func TestAddSongsSerializesInputPlaylist(t *testing.T) {
	repository := spotifyPlaylistRepository()

	repository.CreatePlaylist(defaultContext(), playlistToCreate())

	expected := SpotifyPlaylist{Name: "my-playlist", Description: "some playlist", IsPublic: false}
	actual := repository.GetPlaylistCreateSerializer().(*serialization.FakeSerializer[SpotifyPlaylist]).GetArgs()
	assert.Equal(t, actual, expected)
}

func TestCreatePlaylistReturnsErrorOnPlaylistSerializationError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	serializer := playlistCreateSerializer()
	serializer.SetError(errors.New("test playlist error"))
	repository.SetPlaylistCreateSerializer(serializer)

	_, err := repository.CreatePlaylist(defaultContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistSendsCreateRequestWithOptions(t *testing.T) {
	repository := spotifyPlaylistRepository()

	repository.CreatePlaylist(defaultContext(), playlistToCreate())

	actual := repository.GetHTTPSender().(*httpsender.FakeHTTPSender).GetSendArgs()
	assert.Equal(t, actual, expectedCreatePlaylistHttpOptions())
}

func TestCreatePlaylistReturnsErrorOnDeserializationError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	deserializer := playlistCreateDeserializer()
	deserializer.SetError(errors.New("test create playlist error"))
	repository.SetPlaylistCreateDeserializer(deserializer)

	_, err := repository.CreatePlaylist(defaultContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestCreatePlaylistReturnsIdFromDeserializedResponse(t *testing.T) {
	repository := spotifyPlaylistRepository()
	deserialized := createdPlaylist()
	deserializer := playlistCreateDeserializer()
	deserializer.SetResponse(deserialized)
	repository.SetPlaylistCreateDeserializer(deserializer)

	actual, _ := repository.CreatePlaylist(defaultContext(), playlistToCreate())

	assert.Equal(t, actual, deserialized.Id)
}

func TestCreatePlaylistReturnsErrorOnSendError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(errorSender())

	_, err := repository.CreatePlaylist(defaultContext(), playlistToCreate())

	assert.NotNil(t, err)
}

func TestSearchPlaylistSendsCreateRequestWithOptions(t *testing.T) {
	repository := spotifyPlaylistRepository()

	repository.SearchPlaylist(defaultContext(), searchPlaylistName, searchPlaylistLimit)

	actual := repository.GetHTTPSender().(*httpsender.FakeHTTPSender).GetSendArgs()
	assert.Equal(t, actual, expectedSearchPlaylistHttpOptions())
}

func TestSearchPlaylistReturnsErrorOnSendError(t *testing.T) {
	repository := spotifyPlaylistRepository()
	repository.SetHTTPSender(errorSender())

	_, err := repository.SearchPlaylist(defaultContext(), searchPlaylistName, searchPlaylistLimit)

	assert.NotNil(t, err)
}

func TestSearchPlaylistReturnsErrorOnDeserializationError(t *testing.T) {
	deserializer := playlistSearchDeserializer()
	deserializer.SetError(errors.New("test deserialization error"))
	repository := spotifyPlaylistRepository()
	repository.SetPlaylistSearchDeserializer(deserializer)

	_, err := repository.SearchPlaylist(defaultContext(), searchPlaylistName, searchPlaylistLimit)

	assert.NotNil(t, err)
}

func TestSearchPlaylistReturnsDeserializedContent(t *testing.T) {
	repository := spotifyPlaylistRepository()

	actual, err := repository.SearchPlaylist(defaultContext(), searchPlaylistName, searchPlaylistLimit)

	assert.Nil(t, err)
	assert.Equal(t, actual, searchedFilteredPlaylists())
}

func TestAddSongsPlaylistSendsOptionsUsingSerializerIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	serializer := serialization.NewJsonSerializer[SpotifySongs]()
	repository := spotifyPlaylistRepository()
	repository.SetSongSerializer(&serializer)

	repository.AddSongs(defaultContext(), addSongsPlaylistId, songsToAdd())

	actual := repository.GetHTTPSender().(*httpsender.FakeHTTPSender).GetSendArgs()
	assert.Equal(t, actual, expectedAddSongsHttpOptions())
}

func TestCreatePlaylistSendsOptionsUsingSerializerIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	serializer := serialization.NewJsonSerializer[SpotifyPlaylist]()
	repository := spotifyPlaylistRepository()
	repository.SetPlaylistCreateSerializer(&serializer)

	repository.CreatePlaylist(defaultContext(), playlistToCreate())

	actual := repository.GetHTTPSender().(*httpsender.FakeHTTPSender).GetSendArgs()
	assert.Equal(t, actual, expectedCreatePlaylistHttpOptions())
}

func TestCreatePlaylistReturnsResultsUsingDeserializerIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	repository := spotifyPlaylistRepository()
	sender := fakeSender()
	response := createPlaylistResponseBody()
	sender.SetResponse(&response)
	repository.SetHTTPSender(sender)
	deserializer := serialization.NewJsonDeserializer[SpotifyCreatePlaylistResponse]()
	repository.SetPlaylistCreateDeserializer(&deserializer)

	id, err := repository.CreatePlaylist(defaultContext(), playlistToCreate())

	assert.Nil(t, err)
	assert.Equal(t, id, createPlaylistId)
}

func TestSearchPlaylistReturnsResultsUsingDeserializerIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	sender := fakeSender()
	response := searchedPlaylistsResponseBody()
	sender.SetResponse(&response)
	deserializer := serialization.NewJsonDeserializer[SpotifySearchPlaylistResponse]()
	repository := spotifyPlaylistRepository()
	repository.SetPlaylistSearchDeserializer(deserializer)
	repository.SetHTTPSender(sender)

	actual, err := repository.SearchPlaylist(defaultContext(), searchPlaylistName, searchPlaylistLimit)

	assert.Equal(t, actual, searchedFilteredPlaylists())
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
			repository := spotifyPlaylistRepository()
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
		userIdVal           interface{}
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
			repository := spotifyPlaylistRepository()
			repository.SetUserIdKey(test.repositoryUserIdKey)

			_, err := repository.SearchPlaylist(ctx, searchPlaylistName, searchPlaylistLimit)
			assert.NotNil(t, err)

			_, err = repository.CreatePlaylist(ctx, playlistToCreate())
			assert.NotNil(t, err)
		})
	}
}
