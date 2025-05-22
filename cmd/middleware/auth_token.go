package middleware

import (
	"context"
	"net/http"
	"strings"

	types "festwrap/internal"
)

// Extracts the Bearer Auth token in the header and stores in the context variable with the given key
type AuthTokenExtractor struct {
	tokenKey types.ContextKey
}

func NewAuthTokenExtractor() AuthTokenExtractor {
	return AuthTokenExtractor{tokenKey: types.ContextKey("token")}
}

func (m AuthTokenExtractor) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusBadRequest)
			return
		} else if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Unexpected authorization header format", http.StatusUnprocessableEntity)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		ctxWithToken := context.WithValue(r.Context(), m.tokenKey, token)
		requestWithToken := r.WithContext(ctxWithToken)
		next.ServeHTTP(w, requestWithToken)
	})
}

func (m *AuthTokenExtractor) SetTokenKey(key types.ContextKey) {
	m.tokenKey = key
}
