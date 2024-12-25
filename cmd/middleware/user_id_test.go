package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	types "festwrap/internal"
	"festwrap/internal/testtools"
	"festwrap/internal/user"
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

func defaultContext() context.Context {
	ctx := context.Background()
	context.WithValue(ctx, "some key", "some value")
	return ctx
}

func userIdMiddlewareTestSetup() (UserIdMiddleware, *http.Request, *httptest.ResponseRecorder) {
	userRepository := user.FakeUserRepository{}
	userRepository.SetGetCurrentIdValue(user.GetCurrentIdValue{UserId: "some_id", Err: nil})
	middleware := NewUserIdMiddleware(GetUserIdHandler{}, &userRepository)
	middleware.SetUserIdKey(defaultUserIdKey())
	request := httptest.NewRequest("GET", "http://example.com", nil)
	writer := httptest.NewRecorder()
	return middleware, request, writer
}

func TestGetUserIdCallsRepositoryWithRequestContext(t *testing.T) {
	middleware, request, writer := userIdMiddlewareTestSetup()

	middleware.ServeHTTP(writer, request)

	fakeRepository := middleware.GetUserRepository().(*user.FakeUserRepository)
	testtools.AssertEqual(t, fakeRepository.GetGetCurrentIdArgs().Context, request.Context())
}

func TestGetUserReturnsInternalErrorOnRepositoryError(t *testing.T) {
	middleware, request, writer := userIdMiddlewareTestSetup()
	userRepository := user.FakeUserRepository{}
	userRepository.SetGetCurrentIdValue(user.GetCurrentIdValue{UserId: "", Err: errors.New("test error")})
	middleware.SetUserRepository(&userRepository)

	middleware.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Result().StatusCode, http.StatusInternalServerError)
}

func TestUserIsPlacedInExpectedContextKey(t *testing.T) {
	middleware, request, writer := userIdMiddlewareTestSetup()

	middleware.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Body.String(), "some_id")
}

func TestUserIdMiddlewareReturnsStatusCodeofTheHandler(t *testing.T) {
	middleware, request, writer := userIdMiddlewareTestSetup()

	middleware.ServeHTTP(writer, request)

	testtools.AssertEqual(t, writer.Code, http.StatusContinue)
}
