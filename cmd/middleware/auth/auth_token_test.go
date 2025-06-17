package auth

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "festwrap/internal"
	"festwrap/internal/logging"

	"github.com/stretchr/testify/assert"
)

const (
	accessToken = "clientToken"
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
	authClient := AuthClientMock{}
	authClient.Mock.On("GetAccessToken").Return(accessToken, nil)
	middleware := NewAuthTokenExtractor(&authClient, logging.NoopLogger{})
	middleware.SetTokenKey(defaultTokenKey())
	request := httptest.NewRequest("GET", "http://example.com", nil)
	writer := httptest.NewRecorder()
	return middleware, request, writer
}

func errorAuthClient() AuthClient {
	authClient := AuthClientMock{}
	authClient.Mock.On("GetAccessToken").Return("", errors.New("test auth client error"))
	return &authClient
}

func TestInternalErrorOnAuthClientError(t *testing.T) {
	extractor, request, writer := tokenAuthExtractorTestSetup()
	extractor.SetAuthClient(errorAuthClient())

	extractor.Middleware(GetTokenHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, http.StatusInternalServerError, writer.Code)
}

func TestTokenIsPlacedInExpectedContextKey(t *testing.T) {
	extractor, request, writer := tokenAuthExtractorTestSetup()

	extractor.Middleware(GetTokenHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, accessToken, writer.Body.String())
}

func TestMiddlewareReturnsStatusCodeofTheHandler(t *testing.T) {
	extractor, request, writer := tokenAuthExtractorTestSetup()

	extractor.Middleware(GetTokenHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, http.StatusAccepted, writer.Code)
}
