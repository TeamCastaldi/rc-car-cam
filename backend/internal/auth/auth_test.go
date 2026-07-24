package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func stubHandler(called *bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		w.WriteHeader(http.StatusOK)
	})
}

func TestRequireToken_ValidTokenInQueryParam(t *testing.T) {
	var called bool
	h := RequireToken("secret", stubHandler(&called))

	req := httptest.NewRequest("GET", "/stream?token=secret", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if !called {
		t.Error("expected next handler to be called")
	}
}

func TestRequireToken_ValidTokenInAuthorizationHeader(t *testing.T) {
	var called bool
	h := RequireToken("secret", stubHandler(&called))

	req := httptest.NewRequest("GET", "/stream", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	if !called {
		t.Error("expected next handler to be called")
	}
}

func TestRequireToken_MissingTokenReturns401(t *testing.T) {
	var called bool
	h := RequireToken("secret", stubHandler(&called))

	req := httptest.NewRequest("GET", "/stream", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
	if called {
		t.Error("expected next handler NOT to be called")
	}
}

func TestRequireToken_WrongTokenReturns401(t *testing.T) {
	var called bool
	h := RequireToken("secret", stubHandler(&called))

	req := httptest.NewRequest("GET", "/stream?token=wrong", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", rec.Code)
	}
	if called {
		t.Error("expected next handler NOT to be called")
	}
}
