package middleware

import (
	"context"
	types "festwrap/internal"
	"fmt"
	"net/http"

	"festwrap/internal/logging"
	"festwrap/internal/user"
)

type UserIdExtractor struct {
	userIdKey      types.ContextKey
	userRepository user.UserRepository
	logger         logging.Logger
}

// Adds the current user identifier into the context by using the provided user repository
func NewUserIdExtractor(userRepository user.UserRepository, logger logging.Logger) UserIdExtractor {
	return UserIdExtractor{userIdKey: types.ContextKey("user_id"), userRepository: userRepository, logger: logger}
}

func (m UserIdExtractor) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		currentUserId, err := m.userRepository.GetCurrentUserId(r.Context())
		if err != nil {
			m.logger.Error(fmt.Sprintf("could not retrieve user id: %v", err))
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
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
