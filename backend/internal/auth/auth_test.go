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

func TestRequireToken_EmptyConfiguredTokenAlwaysFails(t *testing.T) {
	var called bool
	h := RequireToken("", stubHandler(&called))

	// Even a request that also supplies an empty/blank credential must not
	// be treated as a match against a misconfigured (blank) token.
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

func TestRequireToken_ValidHeaderTokenSucceedsDespiteWrongQueryToken(t *testing.T) {
	var called bool
	h := RequireToken("secret", stubHandler(&called))

	req := httptest.NewRequest("GET", "/stream?token=stale-or-wrong", nil)
	req.Header.Set("Authorization", "Bearer secret")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 (a valid header token should authorize even with a wrong/stale query token present)", rec.Code)
	}
	if !called {
		t.Error("expected next handler to be called")
	}
}

func TestRequireToken_ValidQueryTokenSucceedsDespiteWrongHeaderToken(t *testing.T) {
	var called bool
	h := RequireToken("secret", stubHandler(&called))

	req := httptest.NewRequest("GET", "/stream?token=secret", nil)
	req.Header.Set("Authorization", "Bearer wrong")
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200 (a valid query token should authorize even with a wrong header token present)", rec.Code)
	}
	if !called {
		t.Error("expected next handler to be called")
	}
}
