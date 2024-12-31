package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	ninaRouter "github.com/JonecoBoy/nina/router"
)

func TestAllowIPMiddleware(t *testing.T) {
	allowedIPs := []string{"192.168.1.1", "10.0.0.*"}
	middleware := AllowIPMiddleware(allowedIPs)

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
		name       string
		ip         string
		wantStatus int
	}{
		{"Allowed IP", "192.168.1.1", http.StatusOK},
		{"Allowed IP with wildcard", "10.0.0.5", http.StatusOK},
		{"Not Allowed IP", "192.168.1.2", http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/hello", nil)
			req.RemoteAddr = tt.ip + ":12345"
			rr := httptest.NewRecorder()

			nr.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("got status %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
