package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "festwrap/internal"
	"festwrap/internal/user"

	"github.com/stretchr/testify/assert"
)

type GetUserIdHandler struct{}

func (h GetUserIdHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, _ := r.Context().Value(defaultUserIdKey()).(string)
	w.WriteHeader(http.StatusContinue)
	fmt.Fprint(w, userId)
}

func defaultUserIdKey() types.ContextKey {
	var userIdKey types.ContextKey = "user_key"
	return userIdKey
}

func userIdExtractorTestSetup() (UserIdExtractor, *http.Request, *httptest.ResponseRecorder) {
	userRepository := user.FakeUserRepository{}
	userRepository.SetGetCurrentIdValue(user.GetCurrentIdValue{UserId: "some_id", Err: nil})
	middleware := NewUserIdExtractor(&userRepository)
	middleware.SetUserIdKey(defaultUserIdKey())
	request := httptest.NewRequest("GET", "http://example.com", nil)
	writer := httptest.NewRecorder()
	return middleware, request, writer
}

func TestGetUserIdCallsRepositoryWithRequestContext(t *testing.T) {
	extractor, request, writer := userIdExtractorTestSetup()

	extractor.Middleware(GetUserIdHandler{}).ServeHTTP(writer, request)

	fakeRepository := extractor.GetUserRepository().(*user.FakeUserRepository)
	assert.Equal(t, request.Context(), fakeRepository.GetGetCurrentIdArgs().Context)
}

func TestGetUserReturnsInternalErrorOnRepositoryError(t *testing.T) {
	extractor, request, writer := userIdExtractorTestSetup()
	userRepository := user.FakeUserRepository{}
	userRepository.SetGetCurrentIdValue(user.GetCurrentIdValue{UserId: "", Err: errors.New("test error")})
	extractor.SetUserRepository(&userRepository)

	extractor.Middleware(GetUserIdHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, http.StatusInternalServerError, writer.Result().StatusCode)
}

func TestUserIsPlacedInExpectedContextKey(t *testing.T) {
	extractor, request, writer := userIdExtractorTestSetup()

	extractor.Middleware(GetUserIdHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, "some_id", writer.Body.String())
}

func TestUserIdMiddlewareReturnsStatusCodeofTheHandler(t *testing.T) {
	extractor, request, writer := userIdExtractorTestSetup()

	extractor.Middleware(GetUserIdHandler{}).ServeHTTP(writer, request)

	assert.Equal(t, http.StatusContinue, writer.Code)
}
