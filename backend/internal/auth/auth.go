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
// "Authorization: Bearer <token>" header. Requests without a match get 401
// and never reach next.
func RequireToken(token string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got := r.URL.Query().Get("token")
		if got == "" {
			if authz := r.Header.Get("Authorization"); strings.HasPrefix(authz, "Bearer ") {
				got = strings.TrimPrefix(authz, "Bearer ")
			}
		}

		if subtle.ConstantTimeCompare([]byte(got), []byte(token)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
