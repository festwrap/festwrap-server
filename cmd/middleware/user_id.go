package middleware

import (
	"context"
	types "festwrap/internal"
	"net/http"

	"festwrap/internal/user"
)

type UserIdMiddleware struct {
	userIdKey      types.ContextKey
	userRepository user.UserRepository
	handler        http.Handler
}

// Adds the current user identifier into the context by using the provided user repository
func NewUserIdMiddleware(handler http.Handler, userRepository user.UserRepository) UserIdMiddleware {
	return UserIdMiddleware{userIdKey: types.ContextKey("user_id"), userRepository: userRepository, handler: handler}
}

func (m UserIdMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	currentUserId, err := m.userRepository.GetCurrentUserId(r.Context())
	if err != nil {
		http.Error(w, "Unexpected error: could not retrieve user id", http.StatusInternalServerError)
	}

	ctxWithUserId := context.WithValue(r.Context(), m.userIdKey, currentUserId)
	requestWithUserId := r.WithContext(ctxWithUserId)
	m.handler.ServeHTTP(w, requestWithUserId)
}

func (m *UserIdMiddleware) SetUserIdKey(key types.ContextKey) {
	m.userIdKey = key
}

func (m UserIdMiddleware) GetUserRepository() user.UserRepository {
	return m.userRepository
}

func (m *UserIdMiddleware) SetUserRepository(repository user.UserRepository) {
	m.userRepository = repository
}
