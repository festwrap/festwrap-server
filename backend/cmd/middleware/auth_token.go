package middleware

import (
	"context"
	"net/http"
	"strings"

	types "festwrap/internal"
)

// Extracts the Bearer Auth token in the header and stores in the context variable with the given key
type AuthTokenMiddleware struct {
	tokenKey types.ContextKey
	handler  http.Handler
}

func NewAuthTokenMiddleware(handler http.Handler) AuthTokenMiddleware {
	return AuthTokenMiddleware{tokenKey: types.ContextKey("token"), handler: handler}
}

func (m AuthTokenMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	m.handler.ServeHTTP(w, requestWithToken)
}

func (m *AuthTokenMiddleware) SetTokenKey(key types.ContextKey) {
	m.tokenKey = key
}
