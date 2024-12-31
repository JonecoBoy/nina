package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	ninaRouter "github.com/JonecoBoy/nina/router"
)

func TestCSRFMiddleware(t *testing.T) {
	// Create a new router
	nr := ninaRouter.NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Register the route with the handler and middleware for unsafe methods only
	nr.GET("/hello", helloHandler, nil)
	nr.HEAD("/hello", helloHandler, nil)
	nr.OPTIONS("/hello", helloHandler, nil)
	nr.TRACE("/hello", helloHandler, nil)
	nr.POST("/hello", helloHandler, []ninaRouter.Middleware{CSRFGenerateMiddleware, CSRFValidateMiddleware})

	tests := []struct {
		name       string
		method     string
		csrfToken  string
		wantStatus int
	}{
		{"Safe method - GET", http.MethodGet, "", http.StatusOK},
		{"Safe method - HEAD", http.MethodHead, "", http.StatusOK},
		{"Safe method - OPTIONS", http.MethodOptions, "", http.StatusOK},
		{"Safe method - TRACE", http.MethodTrace, "", http.StatusOK},
		{"Unsafe method - POST with valid token", http.MethodPost, "valid-token", http.StatusOK},
		{"Unsafe method - POST with invalid token", http.MethodPost, "invalid-token", http.StatusForbidden},
		{"Unsafe method - POST without token", http.MethodPost, "", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/hello", nil)
			if tt.csrfToken != "" {
				req.Header.Set(csrfTokenHeader, tt.csrfToken)
				ctx := context.WithValue(req.Context(), csrfTokenHeader, tt.csrfToken)
				req = req.WithContext(ctx)
			}
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
