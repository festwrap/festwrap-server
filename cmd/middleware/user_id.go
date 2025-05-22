package middleware

import (
	"context"
	types "festwrap/internal"
	"net/http"

	"festwrap/internal/user"
)

type UserIdExtractor struct {
	userIdKey      types.ContextKey
	userRepository user.UserRepository
}

// Adds the current user identifier into the context by using the provided user repository
func NewUserIdExtractor(userRepository user.UserRepository) UserIdExtractor {
	return UserIdExtractor{userIdKey: types.ContextKey("user_id"), userRepository: userRepository}
}

func (m UserIdExtractor) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUserId, err := m.userRepository.GetCurrentUserId(r.Context())
		if err != nil {
			http.Error(w, "Unexpected error: could not retrieve user id", http.StatusInternalServerError)
		}

		ctxWithUserId := context.WithValue(r.Context(), m.userIdKey, currentUserId)
		requestWithUserId := r.WithContext(ctxWithUserId)
		next.ServeHTTP(w, requestWithUserId)
	})
}

func (m *UserIdExtractor) SetUserIdKey(key types.ContextKey) {
	m.userIdKey = key
}

func (m UserIdExtractor) GetUserRepository() user.UserRepository {
	return m.userRepository
}

func (m *UserIdExtractor) SetUserRepository(repository user.UserRepository) {
	m.userRepository = repository
}
