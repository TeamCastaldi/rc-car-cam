// Package auth provides shared-secret token middleware for protecting
// backend routes.
package auth

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// RequireToken wraps next with a check for a matching shared-secret token,
// supplied as either a "token" query parameter or an
// "Authorization: Bearer <token>" header — either one matching is enough.
// Requests without a match get 401 and never reach next. A blank token
// (misconfiguration) always fails closed, rather than treating a blank
// request credential as a match.
func RequireToken(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token == "" || (!tokenMatches(r.URL.Query().Get("token"), token) && !tokenMatches(bearerToken(r), token)) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func bearerToken(r *http.Request) string {
	authz := r.Header.Get("Authorization")
	if !strings.HasPrefix(authz, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authz, "Bearer ")
}

func tokenMatches(candidate, token string) bool {
	return subtle.ConstantTimeCompare([]byte(candidate), []byte(token)) == 1
}
