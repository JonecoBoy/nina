package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBasicAuthMiddleware(t *testing.T) {
	username := "admin"
	password := "password"
	middleware := BasicAuthMiddleware(username, password, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}))

	tests := []struct {
		name       string
		username   string
		password   string
		wantStatus int
	}{
		{"Valid credentials", "admin", "password", http.StatusOK},
		{"Invalid username", "user", "password", http.StatusUnauthorized},
		{"Invalid password", "admin", "pass", http.StatusUnauthorized},
		{"No credentials", "", "", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello", nil)
			if tt.username != "" || tt.password != "" {
				req.SetBasicAuth(tt.username, tt.password)
			}
			rr := httptest.NewRecorder()

			middleware.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
