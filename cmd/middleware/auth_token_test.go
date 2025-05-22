package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "festwrap/internal"

	"github.com/stretchr/testify/assert"
)

type GetTokenHandler struct{}

func (h GetTokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, _ := r.Context().Value(defaultTokenKey()).(string)
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprint(w, token)
}

func defaultTokenKey() types.ContextKey {
	var tokenKey types.ContextKey = "token"
	return tokenKey
}

func tokenAuthExtractorTestSetup() (AuthTokenExtractor, *http.Request, *httptest.ResponseRecorder) {
	middleware := NewAuthTokenExtractor()
	request := httptest.NewRequest("GET", "http://example.com", nil)
	writer := httptest.NewRecorder()
	return middleware, request, writer
}

func TestBadRequestErrorOnMissingAuthHeader(t *testing.T) {
	extractor, request, writer := tokenAuthExtractorTestSetup()

	extractor.Middleware(GetTokenHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, http.StatusBadRequest, writer.Code)
}

func TestUnprocessableEntityErrorOnWronglyFormattedAuthHeader(t *testing.T) {
	extractor, request, writer := tokenAuthExtractorTestSetup()
	request.Header.Set("Authorization", "something")

	extractor.Middleware(GetTokenHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, http.StatusUnprocessableEntity, writer.Code)
}

func TestTokenIsPlacedInExpectedContextKey(t *testing.T) {
	extractor, request, writer := tokenAuthExtractorTestSetup()
	request.Header.Set("Authorization", "Bearer 1234")

	extractor.Middleware(GetTokenHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, "1234", writer.Body.String())
}

func TestMiddlewareReturnsStatusCodeofTheHandler(t *testing.T) {
	extractor, request, writer := tokenAuthExtractorTestSetup()
	request.Header.Set("Authorization", "Bearer 1234")

	extractor.Middleware(GetTokenHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, http.StatusAccepted, writer.Code)
}
