package playlist

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"festwrap/internal/logging"
	mocks "festwrap/internal/playlist/mocks"
	"festwrap/internal/serialization"
	"festwrap/internal/testtools"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func defaultPlaylistId() string {
	return "myId"
}

func defaultUpdateBody() []byte {
	return []byte(`{artists:[{"name":"Silverstein",{"name":"Chinese Football"}]}`)
}

func defaultDeserializedBody() playlistUpdate {
	return playlistUpdate{Artists: []updateArtists{{Name: "Silverstein"}, {"Chinese Football"}}}
}

func alwaysErrorPlaylistService() *mocks.PlaylistServiceMock {
	playlistService := &mocks.PlaylistServiceMock{}
	playlistService.On("AddSetlist", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("test error"))
	return playlistService
}

func partialErrorPlaylistService() *mocks.PlaylistServiceMock {
	playlistService := &mocks.PlaylistServiceMock{}
	playlistService.On("AddSetlist", mock.Anything, defaultPlaylistId(), "Silverstein").Return(errors.New("test error"))
	playlistService.On("AddSetlist", mock.Anything, defaultPlaylistId(), mock.Anything).Return(nil)
	return playlistService
}

func buildRequest(t *testing.T, playlistId string, body []byte) *http.Request {
	t.Helper()
	requestUrl, err := url.Parse("https://example.com/playlist/{playlistId}/update")
	if err != nil {
		t.Errorf("Could not create request: %v", err.Error())
	}

	request := httptest.NewRequest("GET", requestUrl.String(), bytes.NewBuffer(body))
	if playlistId != "" {
		request.SetPathValue("playlistId", playlistId)
	}
	return request
}

func updatePlaylistHandler(t *testing.T) UpdatePlaylistHandler {
	t.Helper()
	deserializer := serialization.FakeDeserializer[playlistUpdate]{}
	deserializer.SetResponse(defaultDeserializedBody())

	playlistService := mocks.NewPlaylistServiceMock()
	playlistService.On("AddSetlist", mock.Anything, defaultPlaylistId(), mock.Anything).Return(nil)

	handler := NewUpdatePlaylistHandler(&playlistService, logging.NoopLogger{})
	handler.SetDeserializer(&deserializer)
	return handler
}

func TestBadRequestIfPlaylistIdNotProvided(t *testing.T) {
	request := buildRequest(t, "", defaultUpdateBody())
	writer := httptest.NewRecorder()
	handler := updatePlaylistHandler(t)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, writer.Code, http.StatusBadRequest)
}

func TestBadRequestOnInvalidBody(t *testing.T) {
	invalidBody := []byte(`{"someInvalidBody":"value"}`)
	request := buildRequest(t, "", invalidBody)
	writer := httptest.NewRecorder()
	handler := updatePlaylistHandler(t)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, writer.Code, http.StatusBadRequest)
}

func TestBadRequestOnMoreArtistsThanAllowed(t *testing.T) {
	request := buildRequest(t, "", defaultUpdateBody())
	writer := httptest.NewRecorder()
	handler := updatePlaylistHandler(t)
	handler.SetMaxArtists(1)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, writer.Code, http.StatusBadRequest)
}

func TestServerStatusOnNoErrors(t *testing.T) {
	request := buildRequest(t, defaultPlaylistId(), defaultUpdateBody())
	writer := httptest.NewRecorder()
	handler := updatePlaylistHandler(t)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusCreated, writer.Code)
}

func TestServerErrorReturnedWhenAllArtistsFailed(t *testing.T) {
	request := buildRequest(t, defaultPlaylistId(), defaultUpdateBody())
	writer := httptest.NewRecorder()
	handler := updatePlaylistHandler(t)
	handler.SetPlaylistService(alwaysErrorPlaylistService())

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestServerStatusOnPartialErrors(t *testing.T) {
	request := buildRequest(t, defaultPlaylistId(), defaultUpdateBody())
	writer := httptest.NewRecorder()
	handler := updatePlaylistHandler(t)
	handler.SetPlaylistService(partialErrorPlaylistService())

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusMultiStatus, writer.Code)
}

func TestHandlerUsesDeserializedBodyIntegration(t *testing.T) {
	testtools.SkipOnShortRun(t)

	request := buildRequest(t, defaultPlaylistId(), defaultUpdateBody())
	writer := httptest.NewRecorder()
	handler := updatePlaylistHandler(t)
	handler.SetDeserializer(serialization.NewJsonDeserializer[playlistUpdate]())

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusCreated, writer.Code)
}
