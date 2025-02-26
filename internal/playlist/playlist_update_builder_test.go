package playlist

import (
	"bytes"
	"errors"
	"festwrap/internal/serialization"
	"festwrap/internal/testtools"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	playlistId     = "myId"
	playlistIdPath = "playlistId"
)

func updateBody() []byte {
	return []byte(`{"artists":[{"name":"Silverstein"},{"name":"Chinese Football"}]}`)
}

func updateArtists() []PlaylistArtist {
	return []PlaylistArtist{{Name: "Silverstein"}, {Name: "Chinese Football"}}
}

func buildRequest(t *testing.T, playlistId string, body []byte) *http.Request {
	t.Helper()
	requestUrl, err := url.Parse(fmt.Sprintf("https://example.com/playlist/{%s}", playlistIdPath))
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
	request := buildRequest(t, "", updateBody())
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestExistingUpdateBuilderReturnsErrorOnIncorrectBody(t *testing.T) {
	invalidBody := []byte("`some_incorrect_body}")
	request := buildRequest(t, playlistId, invalidBody)
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestExistingUpdateBuilderReturnsErrorOnDeserializationError(t *testing.T) {
	request := buildRequest(t, playlistId, updateBody())
	deserializer := serialization.FakeDeserializer[PlaylistArtists]{}
	deserializer.SetError(errors.New("some error"))
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)
	builder.SetDeserializer(&deserializer)

	_, err := builder.Build(request)

	assert.NotNil(t, err)
}

func TestExistingUpdateBuilderReturnsDeserializedContent(t *testing.T) {
	request := buildRequest(t, playlistId, updateBody())
	deserializer := serialization.FakeDeserializer[PlaylistArtists]{}
	deserializer.SetResponse(PlaylistArtists{Artists: updateArtists()})
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)
	builder.SetDeserializer(&deserializer)

	actual, err := builder.Build(request)

	expected := PlaylistUpdate{PlaylistId: playlistId, Artists: updateArtists()}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
	assert.Equal(t, deserializer.GetArgs(), updateBody())
}

func TestExistingUpdateBuilderReturnsExpectedResultIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	request := buildRequest(t, playlistId, updateBody())
	builder := NewExistingPlaylistUpdateBuilder(playlistIdPath)

	actual, err := builder.Build(request)

	expected := PlaylistUpdate{PlaylistId: playlistId, Artists: updateArtists()}
	assert.Equal(t, expected, actual)
	assert.Nil(t, err)
}
