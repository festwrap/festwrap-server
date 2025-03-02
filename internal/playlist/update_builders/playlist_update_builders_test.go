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
	"festwrap/internal/serialization"
	"festwrap/internal/testtools"

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

func updateArtists() []PlaylistArtist {
	return []PlaylistArtist{{Name: "Silverstein"}, {Name: "Chinese Football"}}
}

func newPlaylistUpdate() NewPlaylistUpdate {
	return NewPlaylistUpdate{
		ExistingPlaylistUpdate: ExistingPlaylistUpdate{Artists: updateArtists()},
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

func playlistServiceMock() *mocks.PlaylistServiceMock {
	service := mocks.NewPlaylistServiceMock()
	playlist := newPlaylistUpdate().Playlist.toPlaylist()
	service.On("CreatePlaylist", context.Background(), playlist).Return(playlistId, nil)
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
		request.SetPathValue(playlistIdPath, playlistId)
	}
	return request
}

func TestExistingUpdateBuilderReturnsErrorIfPlaylistIdNotProvided(t *testing.T) {
	request := buildRequest(t, "", existingPlaylistUpdateBody())
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestBuildersReturnErrorOnIncorrectBody(t *testing.T) {
	playlistService := playlistServiceMock()
	existingPlaylistBuilder := NewExistingPlaylistUpdateBuilder(playlistIdPath)
	newPlaylistBuilder := NewNewPlaylistUpdateBuilder(playlistService)
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

func TestExistingUpdateBuilderReturnsErrorOnDeserializationError(t *testing.T) {
	request := buildRequest(t, playlistId, existingPlaylistUpdateBody())
	deserializer := serialization.FakeDeserializer[ExistingPlaylistUpdate]{}
	deserializer.SetError(errors.New("some error"))
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)
	builder.SetDeserializer(&deserializer)

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestNewUpdateBuilderReturnsErrorOnDeserializationError(t *testing.T) {
	request := buildRequest(t, playlistId, newPlaylistUpdateBody())
	deserializer := serialization.FakeDeserializer[NewPlaylistUpdate]{}
	deserializer.SetError(errors.New("some error"))
	playlistService := playlistServiceMock()
	builder := NewNewPlaylistUpdateBuilder(playlistService)
	builder.SetDeserializer(&deserializer)

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestExistingUpdateBuilderReturnsDeserializedContent(t *testing.T) {
	request := buildRequest(t, playlistId, existingPlaylistUpdateBody())
	deserializer := serialization.FakeDeserializer[ExistingPlaylistUpdate]{}
	deserializer.SetResponse(ExistingPlaylistUpdate{Artists: updateArtists()})
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)
	builder.SetDeserializer(&deserializer)

	actual, err := builder.Build(request)

	expected := playlistUpdate()
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
	assert.Equal(t, deserializer.GetArgs(), existingPlaylistUpdateBody())
}

func TestNewUpdateBuilderReturnsExpectedUpdate(t *testing.T) {
	request := buildRequest(t, playlistId, newPlaylistUpdateBody())
	deserializer := serialization.FakeDeserializer[NewPlaylistUpdate]{}
	deserializer.SetResponse(newPlaylistUpdate())
	playlistService := playlistServiceMock()
	builder := NewNewPlaylistUpdateBuilder(playlistService)
	builder.SetDeserializer(&deserializer)

	actual, err := builder.Build(request)

	expected := playlistUpdate()
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}

func TestBuildReturnsExpectedResultIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	playlistService := playlistServiceMock()
	existingPlaylistBuilder := NewExistingPlaylistUpdateBuilder(playlistIdPath)
	newPlaylistBuilder := NewNewPlaylistUpdateBuilder(playlistService)
	tests := map[string]struct {
		builder PlaylistUpdateBuilder
		request *http.Request
	}{
		"existing playlist builder": {
			builder: &existingPlaylistBuilder,
			request: buildRequest(t, playlistId, existingPlaylistUpdateBody()),
		},
		"new playlist builder": {
			builder: &newPlaylistBuilder,
			request: buildRequest(t, playlistId, newPlaylistUpdateBody()),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			actual, err := test.builder.Build(test.request)

			expected := playlistUpdate()
			assert.Equal(t, expected, actual)
			assert.Nil(t, err)
		})
	}
}
