package playlist

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	services "festwrap/cmd/services"
	playlistmocks "festwrap/cmd/services/mocks"
	"festwrap/internal/logging"
	"festwrap/internal/playlist"

	"github.com/stretchr/testify/assert"
)

const (
	playlistId                = "someId"
	playlistName              = "my playlist"
	requestBodyString         = `{"playlist": {"name": "my playlist"}, "artists":[{"name":"Comeback Kid"}, {"name":"Municipal Waste"}]}`
	emptyArtistsBodyString    = `{"playlist": {"name": "my playlist"}, "artists":[]}`
	emptyArtistNameBodyString = `{"playlist": {"name": "my playlist"}, "artists":[{"name":""}, {"name":"Municipal Waste"}]}`
)

func playlistArtists() []string {
	return []string{"Comeback Kid", "Municipal Waste"}
}

func buildRequest(t *testing.T, requestBody []byte) *http.Request {
	t.Helper()
	requestUrl, err := url.Parse("https://example.com/playlist/")
	if err != nil {
		t.Errorf("Could not create request: %v", err.Error())
	}
	return httptest.NewRequest("GET", requestUrl.String(), bytes.NewBuffer(requestBody))
}

func buildPlaylistServiceMock(
	ctx context.Context,
	result services.PlaylistCreation,
	err error,
) *playlistmocks.PlaylistServiceMock {
	playlistService := &playlistmocks.PlaylistServiceMock{}
	playlistService.On(
		"CreatePlaylistWithArtists",
		ctx,
		playlist.PlaylistDetails{Name: playlistName, Description: "", IsPublic: true},
		playlistArtists(),
	).Return(
		result,
		err,
	)
	return playlistService
}

func setup(t *testing.T) (CreatePlaylistHandler, *http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	request := buildRequest(t, []byte(requestBodyString))
	writer := httptest.NewRecorder()

	playlistService := buildPlaylistServiceMock(
		request.Context(),
		services.PlaylistCreation{PlaylistId: playlistId, Status: services.Success},
		nil,
	)

	handler := NewCreatePlaylistHandler(playlistService, logging.NoopLogger{})
	return handler, request, writer
}

func TestCreatePlaylistHandlerReturnsErrorOnIncorrectRequestBody(t *testing.T) {
	handler, _, writer := setup(t)
	request := buildRequest(t, []byte("`some_incorrect_body}"))

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)

}

func TestCreatePlaylistHandlerReturnsErrorOnInvalidRequest(t *testing.T) {
	tests := map[string]struct {
		requestBody         string
		maxArtists          int
		maxArtistNameLength int
	}{
		"no artists": {
			requestBody:         emptyArtistsBodyString,
			maxArtists:          5,
			maxArtistNameLength: 50,
		},
		"more artists than limit": {
			requestBody:         requestBodyString,
			maxArtists:          1,
			maxArtistNameLength: 50,
		},
		"artist length exceeds limit": {
			requestBody:         requestBodyString,
			maxArtists:          1,
			maxArtistNameLength: 5,
		},
		"empty artist name": {
			requestBody:         emptyArtistNameBodyString,
			maxArtists:          5,
			maxArtistNameLength: 50,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler, _, writer := setup(t)
			handler.SetMaxArtistNameLength(test.maxArtistNameLength)
			request := buildRequest(t, []byte(test.requestBody))
			handler.SetMaxArtists(test.maxArtists)

			handler.ServeHTTP(writer, request)

			assert.Equal(t, http.StatusBadRequest, writer.Code)
		})
	}
}

func TestCreatePlaylistHandlerCallsServiceWithExpectedArgs(t *testing.T) {
	handler, request, writer := setup(t)

	handler.ServeHTTP(writer, request)

	playlistService := handler.GetPlaylistService().(*playlistmocks.PlaylistServiceMock)
	playlistService.AssertExpectations(t)
}

func TestCreatePlaylistHandlerReturnsCreatedStatusOnSuccess(t *testing.T) {
	handler, request, writer := setup(t)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusCreated, writer.Code)
}

func TestCreatePlaylistHandlerReturnsMultiStatusOnPartialFailure(t *testing.T) {
	handler, request, writer := setup(t)
	playlistservice := buildPlaylistServiceMock(
		request.Context(),
		services.PlaylistCreation{PlaylistId: playlistId, Status: services.PartialFailure},
		nil,
	)
	handler.SetPlaylistService(playlistservice)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusMultiStatus, writer.Code)
}

func TestCreatePlaylistHandlerReturnsInternalErrorOnServiceError(t *testing.T) {
	handler, request, writer := setup(t)
	playlistservice := buildPlaylistServiceMock(
		request.Context(),
		services.PlaylistCreation{},
		errors.New("test service error"),
	)
	handler.SetPlaylistService(playlistservice)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestCreatePlaylistHandlerReturnsCreatedPlaylistInfo(t *testing.T) {
	handler, request, writer := setup(t)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusCreated, writer.Code)
	expectedBody := fmt.Sprintf("{\"playlist\":{\"id\":\"%s\"}}\n", playlistId)
	assert.Equal(t, expectedBody, writer.Body.String())
}
