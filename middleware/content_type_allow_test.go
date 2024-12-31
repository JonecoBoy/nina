package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ninaRouter "github.com/JonecoBoy/nina/router"
)

func TestAllowedContentTypeMiddleware(t *testing.T) {
	allowedTypes := []string{"application/json", "text/xml"}
	middleware := AllowedContentTypeMiddleware(allowedTypes)

	// Create a new router
	nr := ninaRouter.NewRouter()

	// Define a simple handler
	helloHandler := func(w http.ResponseWriter, r *ninaRouter.NinaRequest) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello, World!"))
	}

	// Register the route with the handler and middleware
	nr.GET("/hello", helloHandler, []ninaRouter.Middleware{middleware})

	tests := []struct {
		name        string
		contentType string
		wantStatus  int
	}{
		{"Allowed content type - application/json", "application/json", http.StatusOK},
		{"Allowed content type - text/xml", "text/xml", http.StatusOK},
		{"Blocked content type - text/plain", "text/plain", http.StatusUnsupportedMediaType},
		{"No content type", "", http.StatusUnsupportedMediaType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello", nil)
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
