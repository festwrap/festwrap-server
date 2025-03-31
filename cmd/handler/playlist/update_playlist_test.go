package playlist

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"festwrap/internal/logging"
	"festwrap/internal/playlist"
	playlistmocks "festwrap/internal/playlist/mocks"
	buildermocks "festwrap/internal/playlist/update_builders/mocks"

	"github.com/stretchr/testify/assert"
)

const (
	playlistId     = "someId"
	playlistIdPath = "playlistId"
)

func updateArtists() []playlist.PlaylistArtist {
	return []playlist.PlaylistArtist{{Name: "Comeback Kid"}, {Name: "Municipal Waste"}}
}

func alwaysSuccessPlaylistService(request *http.Request) *playlistmocks.PlaylistServiceMock {
	playlistService := &playlistmocks.PlaylistServiceMock{}
	playlistService.On("AddSetlist", request.Context(), playlistId, "Municipal Waste").Return(nil)
	playlistService.On("AddSetlist", request.Context(), playlistId, "Comeback Kid").Return(nil)
	return playlistService
}

func alwaysErrorPlaylistService(request *http.Request) *playlistmocks.PlaylistServiceMock {
	playlistService := &playlistmocks.PlaylistServiceMock{}
	playlistService.On("AddSetlist", request.Context(), playlistId, "Municipal Waste").Return(errors.New("error 1"))
	playlistService.On("AddSetlist", request.Context(), playlistId, "Comeback Kid").Return(errors.New("error 2"))
	return playlistService
}

func partialErrorPlaylistService(request *http.Request) *playlistmocks.PlaylistServiceMock {
	playlistService := &playlistmocks.PlaylistServiceMock{}
	playlistService.On("AddSetlist", request.Context(), playlistId, "Municipal Waste").Return(nil)
	playlistService.On("AddSetlist", request.Context(), playlistId, "Comeback Kid").Return(errors.New("error 1"))
	return playlistService
}

func buildRequest(t *testing.T) *http.Request {
	t.Helper()
	requestUrl, err := url.Parse("https://example.com/playlist/")
	if err != nil {
		t.Errorf("Could not create request: %v", err.Error())
	}
	body := []byte(`{"artists":[{"name":"Comeback Kid"}, {"name":"Municipal Waste"}]}`)
	return httptest.NewRequest("GET", requestUrl.String(), bytes.NewBuffer(body))
}

func setup(t *testing.T) (UpdatePlaylistHandler, *http.Request, *httptest.ResponseRecorder) {
	t.Helper()

	request := buildRequest(t)
	writer := httptest.NewRecorder()

	builder := buildermocks.PlaylistUpdateBuilderMock{}
	builder.On("Build", request).Return(playlist.PlaylistUpdate{PlaylistId: playlistId, Artists: updateArtists()}, nil)

	playlistService := alwaysSuccessPlaylistService(request)

	handler := NewUpdatePlaylistHandler(playlistService, &builder, logging.NoopLogger{})
	return handler, request, writer
}

func TestUpdatePlaylistHandlerReturnsErrorOnUpdateBuilderError(t *testing.T) {
	handler, request, writer := setup(t)
	builder := buildermocks.PlaylistUpdateBuilderMock{}
	builder.On("Build", request).Return(playlist.PlaylistUpdate{}, errors.New("test error"))
	handler.SetPlaylistUpdateBuilder(&builder)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestUpdatePlaylistHandlerReturnsErrorIfArtistsOutOfBounds(t *testing.T) {
	tests := map[string]struct {
		artists    []playlist.PlaylistArtist
		maxArtists int
	}{
		"no artists": {
			artists:    []playlist.PlaylistArtist{},
			maxArtists: 5,
		},
		"more artists than limit": {
			artists:    updateArtists(),
			maxArtists: 1,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler, request, writer := setup(t)
			builder := buildermocks.PlaylistUpdateBuilderMock{}
			builder.On("Build", request).Return(
				playlist.PlaylistUpdate{PlaylistId: playlistId, Artists: test.artists},
				nil,
			)
			handler.SetPlaylistUpdateBuilder(&builder)
			handler.SetMaxArtists(test.maxArtists)

			handler.ServeHTTP(writer, request)

			assert.Equal(t, http.StatusBadRequest, writer.Code)
		})
	}
}

func TestUpdatePlaylistHandlerAddsSetlisttWithArtistsFromBuilder(t *testing.T) {
	handler, request, writer := setup(t)

	handler.ServeHTTP(writer, request)

	playlistService := handler.GetPlaylistService().(*playlistmocks.PlaylistServiceMock)
	playlistService.AssertExpectations(t)
}

func TestUpdatePlaylistHandlerStatusOnNoErrors(t *testing.T) {
	handler, request, writer := setup(t)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusCreated, writer.Code)
}

func TestUpdatePlaylistHandlerStatusOnAllFailures(t *testing.T) {
	handler, request, writer := setup(t)
	playlistService := alwaysErrorPlaylistService(request)
	handler.SetPlaylistService(playlistService)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestUpdatePlaylistHandlerStatusOnPartialErrors(t *testing.T) {
	handler, request, writer := setup(t)
	playlistService := partialErrorPlaylistService(request)
	handler.SetPlaylistService(playlistService)

	handler.ServeHTTP(writer, request)

	assert.Equal(t, http.StatusMultiStatus, writer.Code)
}

func TestUpdatePlaylistHandlerShouldReturnResponseIfEnabled(t *testing.T) {
	tests := map[string]struct {
		returnResponse bool
		expected       string
	}{
		"enabled": {
			returnResponse: true,
			expected:       fmt.Sprintf("{\"playlist\":{\"id\":\"%s\"}}\n", playlistId),
		},
		"disabled": {
			returnResponse: false,
			expected:       "",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			handler, request, writer := setup(t)
			handler.ReturnResponse(test.returnResponse)

			handler.ServeHTTP(writer, request)

			assert.Equal(t, http.StatusCreated, writer.Code)
			assert.Equal(t, test.expected, writer.Body.String())
		})
	}
}
