package auth

import (
	"context"
	"fmt"
	"net/http"

	types "festwrap/internal"
	"festwrap/internal/logging"
)

// Obtains an auth access token and stores in the context variable with the given key
type AuthTokenExtractor struct {
	tokenKey   types.ContextKey
	authClient AuthClient
	logger     logging.Logger
}

func NewAuthTokenExtractor(authClient AuthClient, logger logging.Logger) AuthTokenExtractor {
	return AuthTokenExtractor{
		tokenKey:   types.ContextKey("token"),
		authClient: authClient,
		logger:     logger,
	}
}

func (m AuthTokenExtractor) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken, err := m.authClient.GetAccessToken()
		if err != nil {
			m.logger.Error(fmt.Sprintf("could not obtain access token: %v", err))
			http.Error(w, "Unexpected error", http.StatusInternalServerError)
			return
		}
		ctxWithToken := context.WithValue(r.Context(), m.tokenKey, accessToken)
		requestWithToken := r.WithContext(ctxWithToken)
		next.ServeHTTP(w, requestWithToken)
	})
}

func (m *AuthTokenExtractor) SetTokenKey(key types.ContextKey) {
	m.tokenKey = key
}

func (m *AuthTokenExtractor) SetAuthClient(client AuthClient) {
	m.authClient = client
}
