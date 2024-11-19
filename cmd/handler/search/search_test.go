package search

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"festwrap/internal/logging"
	"festwrap/internal/serialization"
	"festwrap/internal/testtools"
)

type Result struct {
	Id    string `json:"id"`
	Value int    `json:"value"`
}

func defaultResults() []Result {
	return []Result{
		{Id: "1", Value: 1},
		{Id: "2", Value: 2},
		{Id: "3", Value: 3},
		{Id: "4", Value: 4},
	}
}

func defaultQueryParams() map[string]string {
	return map[string]string{
		"name":  "someName",
		"limit": "5",
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

func defaultSearcher() Searcher[Result] {
	searcher := NewFakeSearcher[Result]()
	searcher.SetSearchResult(defaultResults())
	return searcher
}

func unmarshalSearchResponse(t *testing.T, bytes []byte) []Result {
	t.Helper()
	var response []Result
	err := json.Unmarshal(bytes, &response)
	if err != nil {
		t.Errorf("Error unmarshalling body: %v", err)
	}
	return response
}

func setup(t *testing.T, params map[string]string) (*httptest.ResponseRecorder, *http.Request, SearchHandler[Result]) {
	t.Helper()
	handler := NewSearchHandler(defaultSearcher(), "someType", logging.NoopLogger{})
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

func TestSearcherCalledWithParams(t *testing.T) {
	searcher := FakeSearcher[Result]{}
	searcher.SetSearchResult(defaultResults())
	handler := NewSearchHandler(&searcher, "someType", logging.NoopLogger{})
	request := buildRequestWithParams(t, defaultQueryParams())

	handler.ServeHTTP(httptest.NewRecorder(), request)

	actual := searcher.GetSearchArgs()
	testtools.AssertEqual(t, actual.Context, request.Context())
	testtools.AssertEqual(t, fmt.Sprint(actual.Limit), defaultQueryParams()["limit"])
	testtools.AssertEqual(t, actual.Name, defaultQueryParams()["name"])
}

func TestSearchReturnsInternalErrorOnEncoderError(t *testing.T) {
	encoder := serialization.FakeEncoder[[]Result]{}
	encoder.SetError(errors.New("test error"))
	writer, request, handler := setup(t, defaultQueryParams())
	handler.SetEncoder(encoder)

	handler.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusInternalServerError)
}

func TestSearchReturnsExpectedResult(t *testing.T) {
	testtools.SkipOnShortRun(t)

	writer, request, handler := setup(t, defaultQueryParams())

	handler.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusOK)
	testtools.AssertEqual(t, unmarshalSearchResponse(t, writer.Body.Bytes()), defaultResults())
}
