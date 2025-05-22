package update_builders

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"festwrap/internal/playlist"
	mocks "festwrap/internal/playlist/mocks"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	playlistId     = "myId"
	playlistIdPath = "playlistId"
)

func existingPlaylistUpdateBody() []byte {
	return []byte(`{"artists":[{"name":"Silverstein"},{"name":"Chinese Football"}]}`)
}

func newPlaylistUpdateBody() []byte {
	return []byte(`{
        "playlist": {
            "name": "Emo songs",
            "description": "Classic emo songs",
            "isPublic": false
        },
        "artists": [
            {"name": "Silverstein"},
            {"name": "Chinese Football"}
        ]
    }`)
}

func newPlaylistUpdate() NewPlaylistUpdate {
	artists := []PlaylistArtist{{Name: "Silverstein"}, {Name: "Chinese Football"}}
	return NewPlaylistUpdate{
		ExistingPlaylistUpdate: ExistingPlaylistUpdate{Artists: artists},
		Playlist:               NewPlaylist{Name: "Emo songs", Description: "Classic emo songs", IsPublic: false},
	}
}

func playlistUpdate() playlist.PlaylistUpdate {
	return playlist.PlaylistUpdate{
		PlaylistId: playlistId,
		Artists: []playlist.PlaylistArtist{
			{Name: "Silverstein"},
			{Name: "Chinese Football"},
		},
	}
}

func playlistService() *mocks.PlaylistServiceMock {
	service := mocks.NewPlaylistServiceMock()
	playlist := newPlaylistUpdate().Playlist.toPlaylist()
	service.On("CreatePlaylist", context.Background(), playlist).Return(playlistId, nil)
	return &service
}

func playlistErrorService() *mocks.PlaylistServiceMock {
	service := mocks.NewPlaylistServiceMock()
	playlist := newPlaylistUpdate().Playlist.toPlaylist()
	service.On("CreatePlaylist", context.Background(), playlist).Return("", errors.New("test error"))
	return &service
}

func buildRequest(t *testing.T, playlistId string, body []byte) *http.Request {
	t.Helper()
	requestUrl, err := url.Parse("https://some_url")
	if err != nil {
		t.Errorf("Could not create request: %v", err.Error())
	}

	request := httptest.NewRequest("GET", requestUrl.String(), bytes.NewBuffer(body))
	if playlistId != "" {
		request = mux.SetURLVars(request, map[string]string{
			playlistIdPath: playlistId,
		})
	}
	return request
}

func existingPlaylistUpdateRequest(t *testing.T) *http.Request {
	t.Helper()
	return buildRequest(t, playlistId, existingPlaylistUpdateBody())
}

func newPlaylistUpdateRequest(t *testing.T) *http.Request {
	t.Helper()
	return buildRequest(t, "", newPlaylistUpdateBody())
}

func TestExistingUpdateBuilderReturnsErrorIfPlaylistIdNotProvided(t *testing.T) {
	request := buildRequest(t, "", existingPlaylistUpdateBody())
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestBuildersReturnErrorOnIncorrectBody(t *testing.T) {
	existingPlaylistBuilder := NewExistingPlaylistUpdateBuilder(playlistIdPath)
	newPlaylistBuilder := NewNewPlaylistUpdateBuilder(playlistService())
	tests := map[string]struct {
		builder PlaylistUpdateBuilder
	}{
		"existing playlist builder": {
			builder: &existingPlaylistBuilder,
		},
		"new playlist builder": {
			builder: &newPlaylistBuilder,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			invalidBody := []byte("`some_incorrect_body}")
			request := buildRequest(t, playlistId, invalidBody)

			_, err := test.builder.Build(request)

			assert.NotNil(t, err)
		})
	}
}

func TestExistingUpdateBuilderReturnsUpdate(t *testing.T) {
	request := existingPlaylistUpdateRequest(t)
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)

	actual, err := builder.Build(request)

	expected := playlistUpdate()
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestNewUpdateBuilderReturnsErrorOnPlaylistCreation(t *testing.T) {
	request := newPlaylistUpdateRequest(t)
	builder := NewNewPlaylistUpdateBuilder(playlistErrorService())

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestNewUpdateBuilderReturnsUpdate(t *testing.T) {
	request := newPlaylistUpdateRequest(t)
	builder := NewNewPlaylistUpdateBuilder(playlistService())

	actual, err := builder.Build(request)

	expected := playlistUpdate()
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
