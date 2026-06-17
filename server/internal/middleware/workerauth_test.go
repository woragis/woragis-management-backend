package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWorkerAuthRejectsMissingKey(t *testing.T) {
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	h := WorkerAuth("secret", next)

	req := httptest.NewRequest(http.MethodGet, "/v1/internal/content/leetcode/dispatch", nil)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status: got %d want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestWorkerAuthAcceptsHeaderAndBearer(t *testing.T) {
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true
		w.WriteHeader(http.StatusOK)
	})
	h := WorkerAuth("secret", next)

	for _, tc := range []struct {
		name   string
		header http.Header
	}{
		{"x-worker-key", http.Header{"X-Worker-Key": []string{"secret"}}},
		{"bearer", http.Header{"Authorization": []string{"Bearer secret"}}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			called = false
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header = tc.header
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK || !called {
				t.Fatalf("expected authorized request")
			}
		})
	}
}
