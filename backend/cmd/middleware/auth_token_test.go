package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "festwrap/internal"
	"festwrap/internal/testtools"
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

func testSetup() (AuthTokenMiddleware, *http.Request, *httptest.ResponseRecorder) {
	middleware := NewAuthTokenMiddleware(defaultTokenKey(), GetTokenHandler{})
	request := httptest.NewRequest("GET", "http://example.com", nil)
	writer := httptest.NewRecorder()
	return middleware, request, writer
}

func TestBadRequestErrorOnMissingAuthHeader(t *testing.T) {
	middleware, request, writer := testSetup()

	middleware.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusBadRequest)
}

func TestUnprocessableEntityErrorOnWronglyFormattedAuthHeader(t *testing.T) {
	middleware, request, writer := testSetup()
	request.Header.Set("Authorization", "something")

	middleware.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusUnprocessableEntity)
}

func TestTokenIsPlacedInExpectedContextKey(t *testing.T) {
	middleware, request, writer := testSetup()
	request.Header.Set("Authorization", "Bearer 1234")

	middleware.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Body.String(), "1234")
}

func TestMiddlewareReturnsStatusCodeofTheHandler(t *testing.T) {
	middleware, request, writer := testSetup()
	request.Header.Set("Authorization", "Bearer 1234")

	middleware.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusAccepted)
}
