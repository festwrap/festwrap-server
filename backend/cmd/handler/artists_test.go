package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"festwrap/internal/artist"
	"festwrap/internal/serialization"
	"festwrap/internal/testtools"
)

func defaultArtists() []artist.Artist {
	return []artist.Artist{
		artist.NewArtist("Brutus"),
		artist.NewArtistWithImageUri("Boysetsfire", "http://some/image.jpg"),
	}
}

func defaultQueryParams() map[string]string {
	return map[string]string{
		"name":  "Brutus",
		"limit": "8",
	}
}

func buildRequestWithParams(t *testing.T, params map[string]string) *http.Request {
	t.Helper()
	requestUrl, err := url.Parse("https://example.com/api")
	if err != nil {
		t.Errorf("Could not create request: %v", err.Error())
	}

	queryParams := requestUrl.Query()
	for name, value := range params {
		queryParams.Add(name, value)
	}

	requestUrl.RawQuery = queryParams.Encode()
	return httptest.NewRequest("GET", requestUrl.String(), nil)
}

func createSearchArtistHandler() SearchArtistHandler {
	repository := artist.FakeArtistRepository{}
	repository.SetSearchReturnValue(defaultArtists())
	return NewSearchArtistHandler(&repository)
}

func unmarshalSearchArtistResponse(t *testing.T, bytes []byte) []artist.Artist {
	t.Helper()
	var response []artist.Artist
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		t.Errorf("Error unmarshalling body: %v", err)
	}
	return response
}

func setup(t *testing.T, params map[string]string) (*httptest.ResponseRecorder, *http.Request, SearchArtistHandler) {
	t.Helper()
	handler := createSearchArtistHandler()
	writer := httptest.NewRecorder()
	request := buildRequestWithParams(t, params)
	return writer, request, handler
}

func TestBadRequestIfNameNotProvided(t *testing.T) {
	params := defaultQueryParams()
	delete(params, "name")
	writer, request, handler := setup(t, params)

	handler.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusBadRequest)
}

func TestLimitStatusCodeDependingOnValue(t *testing.T) {
	tests := map[string]struct {
		limit    string
		maxLimit int
		status   int
	}{
		"below one": {
			limit:    "0",
			maxLimit: 5,
			status:   http.StatusUnprocessableEntity,
		},
		"above max": {
			limit:    "6",
			maxLimit: 5,
			status:   http.StatusUnprocessableEntity,
		},
		"not an integer": {
			limit:    "something",
			maxLimit: 5,
			status:   http.StatusUnprocessableEntity,
		},
		"within limits": {
			limit:    "3",
			maxLimit: 5,
			status:   http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			params := defaultQueryParams()
			params["limit"] = test.limit
			writer, request, handler := setup(t, params)
			handler.SetMaxLimit(test.maxLimit)

			handler.ServeHTTP(writer, request)

			testtools.AssertEqual(t, writer.Code, test.status)
		})
	}
}

func TestSearchArtistRepositoryCalledWithParams(t *testing.T) {
	repository := artist.FakeArtistRepository{}
	repository.SetSearchReturnValue(defaultArtists())
	handler := NewSearchArtistHandler(&repository)
	request := buildRequestWithParams(t, defaultQueryParams())

	handler.ServeHTTP(httptest.NewRecorder(), request)

	actual := repository.GetSearchArtistArgs()
	testtools.AssertEqual(t, actual.Context, request.Context())
	testtools.AssertEqual(t, fmt.Sprint(actual.Limit), defaultQueryParams()["limit"])
	testtools.AssertEqual(t, actual.Name, defaultQueryParams()["name"])
}

func TestSearchArtistReturnsArtists(t *testing.T) {
	writer, request, handler := setup(t, defaultQueryParams())

	handler.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusOK)
	testtools.AssertEqual(t, unmarshalSearchArtistResponse(t, writer.Body.Bytes()), defaultArtists())
}

func TestSearchArtistReturnsInternalErrorOnEncoderError(t *testing.T) {
	encoder := serialization.FakeEncoder[[]artist.Artist]{}
	encoder.SetError(errors.New("test error"))
	writer, request, handler := setup(t, defaultQueryParams())
	handler.SetEncoder(encoder)

	handler.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusInternalServerError)
}
