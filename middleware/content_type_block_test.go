package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ninaRouter "github.com/JonecoBoy/nina/router"
)

func TestContentTypeBlockMiddleware(t *testing.T) {
	blockedTypes := []string{"application/json", "text/xml"}
	middleware := ContentTypeBlockMiddleware(blockedTypes)

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
		{"Blocked content type - application/json", "application/json", http.StatusUnsupportedMediaType},
		{"Blocked content type - text/xml", "text/xml", http.StatusUnsupportedMediaType},
		{"Allowed content type - text/plain", "text/plain", http.StatusOK},
		{"No content type", "", http.StatusOK},
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
